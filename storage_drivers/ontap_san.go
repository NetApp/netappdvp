// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"
)

// OntapSANStorageDriverName is the constant name for this Ontap NAS storage driver
const OntapSANStorageDriverName = "ontap-san"

func init() {
	san := &OntapSANStorageDriver{}
	san.Initialized = false
	Drivers[san.Name()] = san
	log.Debugf("Registered driver '%v'", san.Name())
}

func lunName(name string) string {
	return fmt.Sprintf("/vol/%v/lun0", name)
}

// OntapSANStorageDriver is for iSCSI storage provisioning
type OntapSANStorageDriver struct {
	Initialized bool
	Config      OntapStorageDriverConfig
	API         *ontap.Driver
}

func (d *OntapSANStorageDriver) GetConfig() *OntapStorageDriverConfig {
	return &d.Config
}

func (d *OntapSANStorageDriver) GetAPI() *ontap.Driver {
	return d.API
}

// Name is for returning the name of this driver
func (d OntapSANStorageDriver) Name() string {
	return OntapSANStorageDriverName
}

// Initialize from the provided config
func (d *OntapSANStorageDriver) Initialize(configJSON string, commonConfig *CommonStorageDriverConfig) error {
	log.Debugf("OntapSANStorageDriver#Initialize(...)")

	config := &OntapStorageDriverConfig{}
	config.CommonStorageDriverConfig = commonConfig

	config.IgroupName = "netappdvp"

	// decode configJSON into OntapStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	log.WithFields(log.Fields{
		"Version":           config.Version,
		"StorageDriverName": config.StorageDriverName,
		"DisableDelete":     config.DisableDelete,
	}).Debugf("Reparsed into ontapConfig")

	d.Config = *config
	d.API, err = InitializeOntapDriver(d.Config)
	if err != nil {
		return fmt.Errorf("Problem while initializing: %v", err)
	}

	defaultsErr := PopulateConfigurationDefaults(&d.Config)
	if defaultsErr != nil {
		return fmt.Errorf("Cannot populate configuration defaults: %v", defaultsErr)
	}

	log.WithFields(log.Fields{
		"StoragePrefix": *d.Config.StoragePrefix,
	}).Debugf("Configuration defaults")

	validationErr := d.Validate()
	if validationErr != nil {
		return fmt.Errorf("Problem validating OntapSANStorageDriver: %v", validationErr)
	}

	// log an informational message on a heartbeat
	EmsInitialized(d.Name(), d.API, &d.Config)
	StartEmsHeartbeat(d.Name(), d.API, &d.Config)

	d.Initialized = true

	return nil
}

// Validate the driver configuration and execution environment
func (d *OntapSANStorageDriver) Validate() error {
	log.Debugf("OntapSANStorageDriver#Validate()")

	zr := &azgo.ZapiRunner{
		ManagementLIF: d.Config.ManagementLIF,
		SVM:           d.Config.SVM,
		Username:      d.Config.Username,
		Password:      d.Config.Password,
		Secure:        true,
	}

	r0, err0 := azgo.NewSystemGetVersionRequest().ExecuteUsing(zr)
	if err0 != nil {
		return fmt.Errorf("Could not validate credentials for %v@%v, error: %v", d.Config.Username, d.Config.SVM, err0)
	}

	// Add system version validation, if needed, this is a sanity check right now
	systemVersion := r0.Result
	if systemVersion.VersionPtr == nil {
		return fmt.Errorf("Could not determine system version for %v@%v", d.Config.Username, d.Config.SVM)
	}

	r1, err1 := d.API.NetInterfaceGet()
	if err1 != nil {
		return fmt.Errorf("Problem checking network interfaces; error: %v", err1)
	}

	// if they didn't set a lif to use in the config, we'll set it to the first iscsi lif we happen to find
	if d.Config.DataLIF == "" {
		for _, attrs := range r1.Result.AttributesList() {
			for _, protocol := range attrs.DataProtocols() {
				if protocol == "iscsi" {
					log.Debugf("Setting iSCSI protocol access to '%v'", attrs.Address())
					d.Config.DataLIF = string(attrs.Address())
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
				log.Debugf("Comparing iSCSI protocol access on: '%v' vs: '%v'", attrs.Address(), d.Config.DataLIF)
				if string(attrs.Address()) == d.Config.DataLIF {
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

	// Make sure this host is logged into the ONTAP iSCSI target
	err := utils.EnsureIscsiSession(d.Config.DataLIF)
	if err != nil {
		return fmt.Errorf("Could not establish iSCSI session. %v", err)
	}

	return nil
}

// Create a volume+LUN with the specified options
func (d *OntapSANStorageDriver) Create(name string, sizeBytes uint64, opts map[string]string) error {
	log.Debugf("OntapSANStorageDriver#Create(%v)", name)

	// If the volume already exists, bail out
	response, _ := d.API.VolumeSize(name)
	if isPassed(response.Result.ResultStatusAttr) {
		return fmt.Errorf("Volume already exists")
	}

	// get options with default fallback values
	// see also: ontap_common.go#PopulateConfigurationDefaults
	size := strconv.FormatUint(sizeBytes, 10)
	spaceReserve := utils.GetV(opts, "spaceReserve", d.Config.SpaceReserve)
	snapshotPolicy := utils.GetV(opts, "snapshotPolicy", d.Config.SnapshotPolicy)
	unixPermissions := utils.GetV(opts, "unixPermissions", d.Config.UnixPermissions)
	snapshotDir := utils.GetV(opts, "snapshotDir", d.Config.SnapshotDir)
	exportPolicy := utils.GetV(opts, "exportPolicy", d.Config.ExportPolicy)
	aggregate := utils.GetV(opts, "aggregate", d.Config.Aggregate)
	securityStyle := utils.GetV(opts, "securityStyle", d.Config.SecurityStyle)

	log.WithFields(log.Fields{
		"name":            name,
		"size":            size,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
		"snapshotDir":     snapshotDir,
		"exportPolicy":    exportPolicy,
		"aggregate":       aggregate,
		"securityStyle":   securityStyle,
	}).Debug("Creating volume with values")

	// create the volume
	response1, error1 := d.API.VolumeCreate(name, aggregate, size, spaceReserve, snapshotPolicy, unixPermissions, exportPolicy, securityStyle)
	if !isPassed(response1.Result.ResultStatusAttr) || error1 != nil {
		if response1.Result.ResultErrnoAttr != azgo.EAPIERROR {
			return fmt.Errorf("Error creating volume\n%verror: %v", response1.Result, error1)
		} else {
			if !strings.HasSuffix(strings.TrimSpace(response1.Result.ResultReasonAttr), "Job exists") {
				return fmt.Errorf("Error creating volume\n%verror: %v", response1.Result, error1)
			} else {
				log.Warnf("%v volume create job already exists, skipping volume create on this node...", name)
				return nil
			}
		}
	}

	lunPath := lunName(name)
	osType := "linux"
	spaceReserved := false

	// create the lun
	response2, err2 := d.API.LunCreate(lunPath, int(sizeBytes), osType, spaceReserved)
	if !isPassed(response2.Result.ResultStatusAttr) || err2 != nil {
		return fmt.Errorf("Error creating LUN\n%verror: %v", response2.Result, err2)
	}

	return nil
}

// Create a volume clone
func (d *OntapSANStorageDriver) CreateClone(name, source, snapshot string) error {
	return CreateOntapClone(name, source, snapshot, d.API)
}

// Destroy the requested (volume,lun) storage tuple
func (d *OntapSANStorageDriver) Destroy(name string) error {
	log.Debugf("OntapSANStorageDriver#Destroy(%v)", name)

	lunPath := lunName(name)

	// validate LUN+volume exists before trying to destroy
	response0, _ := d.API.VolumeSize(name)
	if !isPassed(response0.Result.ResultStatusAttr) {
		log.Debugf("%v already deleted, skipping destroy", name)
		return nil
	}

	// lun offline
	response, err := d.API.LunOffline(lunPath)
	if !isPassed(response.Result.ResultStatusAttr) || err != nil {
		log.Warnf("Error attempting to offline lun: %v\n%verror: %v", lunPath, response.Result, err)
	}

	// lun destroy
	response2, err2 := d.API.LunDestroy(lunPath)
	if !isPassed(response2.Result.ResultStatusAttr) || err2 != nil {
		log.Warnf("Error destroying lun: %v\n%verror: %v", lunPath, response2.Result, err2)
	}

	// perform rediscovery to remove the deleted LUN
	utils.MultipathFlush() // flush unused paths
	utils.IscsiRescan()

	response3, error3 := d.API.VolumeDestroy(name, true)
	if !isPassed(response3.Result.ResultStatusAttr) || error3 != nil {
		if response3.Result.ResultErrnoAttr != azgo.EVOLUMEDOESNOTEXIST {
			return fmt.Errorf("Error destroying volume: %v\n%verror: %v", name, response3.Result, error3)
		} else {
			log.Warnf("Volume already deleted while destroying volume: %v\n%verror: %v", name, response3.Result, error3)
		}
	}

	return nil
}

// Attach the lun
func (d *OntapSANStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("OntapSANStorageDriver#Attach(%v, %v, %v)", name, mountpoint, opts)

	// error if no 'iscsi session' exsits for the specified iscsi portal
	sessionExists, sessionExistsErr := utils.IscsiSessionExists(d.Config.DataLIF)
	if sessionExistsErr != nil {
		return fmt.Errorf("Unexpected iSCSI session error: %v", sessionExistsErr)
	}
	if !sessionExists {
		// TODO automatically login for the user if no session detected?
		return fmt.Errorf("Expected iSCSI session %v NOT found, please login to the iscsi portal", d.Config.DataLIF)
	}

	igroupName := d.Config.IgroupName
	lunPath := lunName(name)

	// igroup create
	response, err := d.API.IgroupCreate(igroupName, "iscsi", "linux")
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
		response2, err2 := d.API.IgroupAdd(igroupName, iqn)
		if !isPassed(response2.Result.ResultStatusAttr) {
			if response2.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_HAS_NODE {
				return fmt.Errorf("Problem adding iqn: %v to igroup: %v\n%verror: %v", iqn, igroupName, response2.Result, err2)
			}
		}
	}

	// check if already mapped, so we don't map again
	lunID := 0
	alreadyMapped := false
	response5, _ := d.API.LunMapListInfo(lunPath)
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
			response4, err4 := d.API.LunMap(igroupName, lunPath, i)
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
		if e.PortalIP == d.Config.DataLIF {
			sessionInfoToUse = sessionInfo[i]
		}
	}

	for i, e := range info {
		log.WithFields(log.Fields{
			"i":                i,
			"scsiHost":         e.Host,
			"scsiChannel":      e.Channel,
			"scsiTarget":       e.Target,
			"scsiLun":          e.LUN,
			"multipathDevFile": e.MultipathDevice,
			"devFile":          e.Device,
			"fsType":           e.Filesystem,
			"iqn":              e.IQN,
		}).Debug("Found")
	}

	// look for the expected mapped lun
	for i, e := range info {

		log.WithFields(log.Fields{
			"i":                i,
			"scsiHost":         e.Host,
			"scsiChannel":      e.Channel,
			"scsiTarget":       e.Target,
			"scsiLun":          e.LUN,
			"multipathDevFile": e.MultipathDevice,
			"devFile":          e.Device,
			"fsType":           e.Filesystem,
			"iqn":              e.IQN,
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

// Return the list of snapshots associated with the named volume
func (d *OntapSANStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {
	return GetSnapshotList(name, d.API)
}

// Return the list of volumes associated with this tenant
func (d *OntapSANStorageDriver) List() ([]string, error) {
	return GetVolumeList(d.API, &d.Config)
}

// Test for the existence of a volume
func (d *OntapSANStorageDriver) Get(name string) error {
	return GetVolume(name, d.API)
}
