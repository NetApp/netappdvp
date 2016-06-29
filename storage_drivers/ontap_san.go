// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

// OntapSANStorageDriverName is the constant name for this Ontap NAS storage driver
const OntapSANStorageDriverName = "ontap-san"

func init() {
	san := &OntapSANStorageDriver{}
	san.initialized = false
	Drivers[san.Name()] = san
	log.Debugf("Registered driver '%v'", san.Name())
}

func lunName(name string) string {
	return fmt.Sprintf("/vol/%v/lun0", name)
}

// OntapSANStorageDriver is for iSCSI storage provisioning
type OntapSANStorageDriver struct {
	initialized bool
	config      OntapStorageDriverConfig
	api         *ontap.Driver
}

// Name is for returning the name of this driver
func (d OntapSANStorageDriver) Name() string {
	return OntapSANStorageDriverName
}

// Initialize from the provided config
func (d *OntapSANStorageDriver) Initialize(configJSON string) error {
	log.Debugf("OntapSANStorageDriver#Initialize(...)")

	config := &OntapStorageDriverConfig{}
	config.IgroupName = "netappdvp"

	// decode configJSON into OntapStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	log.WithFields(log.Fields{
		"Version":           config.Version,
		"StorageDriverName": config.StorageDriverName,
		"Debug":             config.Debug,
		"DisableDelete":     config.DisableDelete,
		"StoragePrefixRaw":  string(config.StoragePrefixRaw),
	}).Debugf("Reparsed into ontapConfig")

	d.config = *config
	d.api, err = InitializeOntapDriver(d.config)
	if err != nil {
		return fmt.Errorf("Problem while initializing, error: %v", err)
	}

	validationErr := d.Validate()
	if validationErr != nil {
		return fmt.Errorf("Problem validating OntapSANStorageDriver error: %v", validationErr)
	}

	// log an informational message when this plugin starts
	EmsInitialized(d.Name(), d.api)

	d.initialized = true
	log.Infof("Successfully initialized Ontap SAN Docker driver version %v", DriverVersion)
	return nil
}

// Validate the driver configuration and execution environment
func (d *OntapSANStorageDriver) Validate() error {
	log.Debugf("OntapSANStorageDriver#Validate()")

	zr := &azgo.ZapiRunner{
		ManagementLIF: d.config.ManagementLIF,
		SVM:           d.config.SVM,
		Username:      d.config.Username,
		Password:      d.config.Password,
		Secure:        true,
	}

	r0, err0 := azgo.NewSystemGetVersionRequest().ExecuteUsing(zr)
	if err0 != nil {
		return fmt.Errorf("Could not validate credentials for %v@%v, error: %v", d.config.Username, d.config.SVM, err0)
	}

	// Add system version validation, if needed, this is a sanity check right now
	systemVersion := r0.Result
	if systemVersion.VersionPtr == nil {
		return fmt.Errorf("Could not determine system version for %v@%v", d.config.Username, d.config.SVM)
	}

	r1, err1 := d.api.NetInterfaceGet()
	if err1 != nil {
		return fmt.Errorf("Problem checking network interfaces; error: %v", err1)
	}

	// if they didn't set a lif to use in the config, we'll set it to the first iscsi lif we happen to find
	if d.config.DataLIF == "" {
		for _, attrs := range r1.Result.AttributesList() {
			for _, protocol := range attrs.DataProtocols() {
				if protocol == "iscsi" {
					log.Debugf("Setting iSCSI protocol access to '%v'", attrs.Address())
					d.config.DataLIF = string(attrs.Address())
				}
			}
		}
	}

	// now, we validate our settings
	foundIscsi := false
	iscsiLifCount := 0
	for _, attrs := range r1.Result.AttributesList() {
		for _, protocol := range attrs.DataProtocols() {
			if protocol == "iscsi" {
				log.Debugf("Comparing iSCSI protocol access on: '%v' vs: '%v'", attrs.Address(), d.config.DataLIF)
				if string(attrs.Address()) == d.config.DataLIF {
					foundIscsi = true
					iscsiLifCount++
				}
			}
		}
	}

	if iscsiLifCount > 1 {
		log.Debugf("Found multiple iSCSI lifs")
	}

	if !foundIscsi {
		return fmt.Errorf("Could not find iSCSI DataLIF")
	}

	isIscsiSupported := utils.IscsiSupported()
	if !isIscsiSupported {
		return fmt.Errorf("iSCSI support not detected")
	}

	// error if no 'iscsi session' exsits for the specified iscsi portal
	sessionExists, sessionExistsErr := utils.IscsiSessionExists(d.config.DataLIF)
	if sessionExistsErr != nil {
		return fmt.Errorf("Unexpected iSCSI session error: %v", sessionExistsErr)
	}
	if !sessionExists {
		// TODO automatically login for the user if no session detected?
		return fmt.Errorf("Expected iSCSI session %v NOT found, please login to the iscsi portal", d.config.DataLIF)
	}

	return nil
}

// Create a volume+LUN with the specified options
func (d *OntapSANStorageDriver) Create(name string, opts map[string]string) error {
	log.Debugf("OntapSANStorageDriver#Create(%v)", name)

	response, _ := d.api.VolumeSize(name)
	if isPassed(response.Result.ResultStatusAttr) {
		log.Debugf("%v already exists, skipping create...", name)
		return nil
	}

	// get options with default values if not specified in config file
	volumeSize := utils.GetV(opts, "size", "1g")
	spaceReserve := utils.GetV(opts, "spaceReserve", "none")
	snapshotPolicy := utils.GetV(opts, "snapshotPolicy", "none")
	unixPermissions := utils.GetV(opts, "unixPermissions", "---rwxr-xr-x")

	log.WithFields(log.Fields{
		"name":            name,
		"volumeSize":      volumeSize,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
	}).Debug("Creating volume with values")

	// create the volume
	response1, error1 := d.api.VolumeCreate(name, d.config.Aggregate, volumeSize, spaceReserve, snapshotPolicy, unixPermissions)
	if !isPassed(response1.Result.ResultStatusAttr) || error1 != nil {
		return fmt.Errorf("Error creating volume:\n%verror: %v", response1.Result, error1)
	}

	lunPath := lunName(name)
	osType := "linux"
	spaceReserved := false

	// lunSize takes some effort; we must convert user friendly strings to total bytes; ex "4KB" -> 4096
	convertedSize, convertErr := utils.ConvertSizeToBytes(volumeSize)
	if convertErr != nil {
		return fmt.Errorf("Cannot convert size to bytes: %v error: %v", volumeSize, convertErr)
	}
	lunSize, atoiErr := strconv.Atoi(convertedSize)
	if atoiErr != nil {
		return fmt.Errorf("Cannot convert size to bytes: %v error: %v", volumeSize, atoiErr)
	}

	// create the lun
	response2, err2 := d.api.LunCreate(lunPath, lunSize, osType, spaceReserved)
	if !isPassed(response2.Result.ResultStatusAttr) || err2 != nil {
		return fmt.Errorf("Error creating LUN\n%verror: %v", response2.Result, err2)
	}

	return nil
}

// Destroy the requested (volume,lun) storage tuple
func (d *OntapSANStorageDriver) Destroy(name string) error {
	log.Debugf("OntapSANStorageDriver#Destroy(%v)", name)

	lunPath := lunName(name)

	// validate LUN+volume exists before trying to destroy
	response0, _ := d.api.VolumeSize(name)
	if !isPassed(response0.Result.ResultStatusAttr) {
		log.Debugf("%v already deleted, skipping destroy", name)
		return nil
	}

	// lun offline
	response, err := d.api.LunOffline(lunPath)
	if !isPassed(response.Result.ResultStatusAttr) || err != nil {
		log.Warnf("Error attempting to offline lun: %v\n%verror: %v", lunPath, response.Result, err)
	}

	// lun destroy
	response2, err2 := d.api.LunDestroy(lunPath)
	if !isPassed(response2.Result.ResultStatusAttr) || err2 != nil {
		log.Warnf("Error destroying lun: %v\n%verror: %v", lunPath, response2.Result, err2)
	}

	// perform rediscovery to remove the deleted LUN
	utils.MultipathFlush() // flush unused paths
	utils.IscsiRescan()

	response3, error3 := d.api.VolumeDestroy(name, true)
	if !isPassed(response3.Result.ResultStatusAttr) || error3 != nil {
		if response3.Result.ResultErrnoAttr != azgo.EVOLUMEDOESNOTEXIST {
			return fmt.Errorf("Error destroying volume: %v error: %v", name, error3)
		} else {
			log.Warnf("Volume already deleted while destroying volume: %v\n%verror: %v", name, response3.Result, error3)
		}
	}

	return nil
}

// Attach the lun
func (d *OntapSANStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("OntapSANStorageDriver#Attach(%v, %v, %v)", name, mountpoint, opts)

	igroupName := d.config.IgroupName
	lunPath := lunName(name)

	// igroup create
	response, err := d.api.IgroupCreate(igroupName, "iscsi", "linux")
	if !isPassed(response.Result.ResultStatusAttr) {
		if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_EXISTS {
			return fmt.Errorf("Problem creating igroup: %v\n%verror: %v", igroupName, response.Result, err)
		}
	}

	// lookp host iqns
	iqns, errIqn := utils.GetInitiatorIqns()
	if errIqn != nil {
		return fmt.Errorf("Problem determining host initiator iqns error: %v", errIqn)
	}

	// igroup add each iqn we found
	for _, iqn := range iqns {
		response2, err2 := d.api.IgroupAdd(igroupName, iqn)
		if !isPassed(response2.Result.ResultStatusAttr) {
			if response2.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_HAS_NODE {
				return fmt.Errorf("Problem adding iqn: %v to igroup: %v\n%verror: %v", iqn, igroupName, response2.Result, err2)
			}
		}
	}

	// check if already mapped, so we don't map again
	lunID := 0
	alreadyMapped := false
	response5, _ := d.api.LunMapListInfo(lunPath)
	if response5.Result.ResultStatusAttr == "passed" {
		if response5.Result.InitiatorGroups() != nil {
			if len(response5.Result.InitiatorGroups()) > 0 {
				lunID = response5.Result.InitiatorGroups()[0].LunId()
				alreadyMapped = true
				log.Debugf("found already mapped lunID: %v", lunID)
			}
		}
	}

	// map IFF not already mapped
	if !alreadyMapped {
		// spin until we get a lunId that works
		// TODO find one directly instead of spinning-and-looking for one?
		for i := 0; i < 4096; i++ {
			response4, err4 := d.api.LunMap(igroupName, lunPath, i)
			if response4.Result.ResultStatusAttr == "passed" {
				lunID = i
				break
			}

			if response4.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_HAS_LUN {
				return fmt.Errorf("Problem mapping lun: %v error: %v", lunPath, err4)
			}
		}
	}
	log.Debugf("using lunID == %v ", lunID)

	// perform discovery to see the created/mapped LUN
	utils.IscsiRescan()

	// lookup all the scsi device information
	info, infoErr := utils.GetDeviceInfoForLuns()
	if infoErr != nil {
		return fmt.Errorf("Problem getting scsi device information, error: %v", infoErr)
	}

	// lookup all the iSCSI session information
	sessionInfo, sessionInfoErr := utils.GetIscsiSessionInfo()
	if sessionInfoErr != nil {
		return fmt.Errorf("Problem getting iSCSI session information, error: %v", sessionInfoErr)
	}

	sessionInfoToUse := utils.IscsiSessionInfo{}
	for i, e := range sessionInfo {
		if e.PortalIP == d.config.DataLIF {
			sessionInfoToUse = sessionInfo[i]
		}
	}

	// look for the expected mapped lun
	for i, e := range info {

		log.WithFields(log.Fields{
			"i": i,
			"e": e,
		}).Debug("Checking")

		if e.LUN != strconv.Itoa(lunID) {
			log.Debugf("Skipping... lun id %v != %v", e.LUN, lunID)
			continue
		}

		if !strings.HasPrefix(e.IQN, sessionInfoToUse.TargetName) {
			log.Debugf("Skipping... %v doesn't start with %v", e.IQN, sessionInfoToUse.TargetName)
			continue
		}

		// if we're here then, we should be on the right info element:
		// *) we have the expected LUN ID
		// *) we have the expected iscsi session target
		log.Debugf("Using... %v", e)

		// make sure we use the proper device (multipath if in use)
		deviceToUse := e.Device
		if e.MultipathDevice != "" {
			deviceToUse = e.MultipathDevice
		}

		if deviceToUse == "" {
			return fmt.Errorf("Could not determine device to use for: %v ", name)
		}

		// put a filesystem on it if there isn't one already there
		if e.Filesystem == "" {
			// format it
			err := utils.FormatVolume(deviceToUse, "ext4") // TODO externalize fsType
			if err != nil {
				return fmt.Errorf("Problem formatting lun: %v device: %v error: %v", name, deviceToUse, err)
			}
		}

		// mount it
		err := utils.Mount(deviceToUse, mountpoint)
		if err != nil {
			return fmt.Errorf("Problem mounting lun: %v device: %v mountpoint: %v error: %v", name, deviceToUse, mountpoint, err)
		}
		return nil
	}

	return nil
}

// Detach the volume
func (d *OntapSANStorageDriver) Detach(name, mountpoint string) error {
	log.Debugf("OntapSANStorageDriver#Detach(%v, %v)", name, mountpoint)

	cmd := fmt.Sprintf("umount %s", mountpoint)
	log.Debugf("cmd==%s", cmd)
	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.Debugf("out==%v", string(out))
		return fmt.Errorf("Problem unmounting docker volume: %v mountpoint: %v error: %v", name, mountpoint, err)
	}

	return nil
}

// DefaultStoragePrefix is the driver specific prefix for created storage, can be overridden in the config file
func (d *OntapSANStorageDriver) DefaultStoragePrefix() string {
	return "netappdvp_"
}
