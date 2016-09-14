// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"strconv"
	"strings"

	"github.com/netapp/netappdvp/apis/eseries"
	"github.com/netapp/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

func init() {
	san := &ESeriesStorageDriver{}
	san.initialized = false
	Drivers[san.Name()] = san
	log.Debugf("Registered driver '%v'", san.Name())
}

// ESeriesStorageDriver is for storage provisioning via Web Services Proxy RESTful interface that communicates with E-Series controller via SYMbol API
type ESeriesStorageDriver struct {
	initialized bool
	config      ESeriesStorageDriverConfig
	storage     *eseries.Driver
}

// Name is for returning the name of this driver
func (d ESeriesStorageDriver) Name() string {
	return "eseries-iscsi"
}

// Initialize from the provided config
func (d *ESeriesStorageDriver) Initialize(configJSON string) error {
	log.Debugf("ESeriesStorageDriver#Initialize(...)")

	config := &ESeriesStorageDriverConfig{}

	// decode configJSON into ESeriesStorageDriverConfig object
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
	}).Debugf("Reparsed into eseriesConfig")

	d.config = *config
	d.storage = eseries.NewDriver(eseries.DriverConfig{
		WebProxyHostname:  config.WebProxyHostname,
		WebProxyPort:      config.WebProxyPort,
		WebProxyUseHTTP:   config.WebProxyUseHTTP,
		WebProxyVerifyTLS: config.WebProxyVerifyTLS,
		Username:          config.Username,
		Password:          config.Password,
		ControllerA:       config.ControllerA,
		ControllerB:       config.ControllerB,
		PasswordArray:     config.PasswordArray,
		ArrayRegistered:   config.ArrayRegistered,
		HostDataIP:        config.HostDataIP,
	})

	validationErr := d.Validate()
	if validationErr != nil {
		return fmt.Errorf("Problem validating ESeriesStorageDriver error: %v", validationErr)
	}

	//Connect to web services proxy
	response, error := d.storage.Connect()
	if error != nil {
		return fmt.Errorf("Problem connecting to Web Services Proxy - ESeriesStorageDriver error: %v", error)
	} else {
		log.Debugf("Connect to Web Services Proxy Success! response=%v", response)
	}

	d.initialized = true
	log.Infof("Successfully initialized E-Series Docker driver version %v", DriverVersion)
	return nil
}

// Validate the driver configuration and execution environment
func (d *ESeriesStorageDriver) Validate() error {
	log.Debugf("ESeriesStorageDriver#Validate()")

	//Make sure the essential information was specified in the json config
	if d.config.WebProxyHostname == "" {
		return fmt.Errorf("WebProxyHostname is empty! You must specify the host/IP for the Web Services Proxy.")
	}

	if d.config.ControllerA == "" || d.config.ControllerB == "" {
		return fmt.Errorf("ControllerA or ControllerB are empty! You must specify the host/IP for the E-Series storage array. If it is a simplex array just specify the same host/IP twice.")
	}

	if d.config.HostDataIP == "" {
		return fmt.Errorf("HostDataIP is empty! You need to specify atleast one of the iSCSI interface IP addresses that is connected to the E-Series array.")
	}

	//Make sure iSCSI is supported on system
	isIscsiSupported := utils.IscsiSupported()
	if !isIscsiSupported {
		return fmt.Errorf("iSCSI support not detected")
	}

	// error if no 'iscsi session' exsits for the specified iscsi portal
	sessionExists, sessionExistsErr := utils.IscsiSessionExists(d.config.HostDataIP)
	if sessionExistsErr != nil {
		return fmt.Errorf("Unexpected iSCSI session error: %v", sessionExistsErr)
	}
	if !sessionExists {
		// TODO automatically login for the user if no session detected?
		return fmt.Errorf("Expected iSCSI session %v NOT found, please login to the iscsi portal", d.config.HostDataIP)
	}

	return nil
}

// Create a volume+LUN with the specified options
func (d *ESeriesStorageDriver) Create(name string, opts map[string]string) error {
	log.Debugf("ESeriesStorageDriver#Create(%v)", name)

	//Example GET point for storage pools:
	//	http://10.251.228.75:8080/devmgr/v2/storage-systems/984ce9e3-46fe-402d-ac59-f4957a7c8288/storage-pools

	volumeSize := utils.GetV(opts, "size", "1g")
	mediaType := utils.GetV(opts, "mediaType", "hdd")
	//mediaSecure := utils.GetV(opts, "mediaSecure", "false")

	volumeGroupRef, error := d.storage.VerifyVolumePools(mediaType, volumeSize)
	if error != nil {
		return error
	} else {
		log.Debugf("ESeriesStorageDriver#Create(%v) - volumeGroupRef=%s", name, volumeGroupRef)
	}

	//Create the volume
	volumeRef, error1 := d.storage.CreateVolume(name, volumeGroupRef, volumeSize, mediaType)
	if error1 != nil {
		return error1
	} else {
		log.Debugf("ESeriesStorageDriver#Create(%v) - volumeRef=%s", name, volumeRef)
	}

	return nil
}

// Destroy the requested (volume,lun) storage tuple
func (d *ESeriesStorageDriver) Destroy(name string) error {
	log.Debugf("ESeriesStorageDriver#Destroy(%v)", name)

	//Grab the host IQN and verify array is aware of it
	iqns, errIqn := utils.GetInitiatorIqns()
	if errIqn != nil {
		return fmt.Errorf("Problem determining host initiator iqns error: %v", errIqn)
	}

	//Going to assume a single IQN name for our host right now
	var iqn string = iqns[0]

	log.Debugf("ESeriesStorageDriver#Destroy(%v) - iqn=%s", name, iqn)

	//We don't want to fail the operation if we can't find a host matching the IQN, but we want to log a warning
	hostRef, verifyIqnErr := d.storage.VerifyHostIQN(iqn)
	if verifyIqnErr != nil {
		log.Warnf("Host IQN (%s) not found on target E-Series array! error=%s", iqn, verifyIqnErr)
	}

	log.Debugf("ESeriesStorageDriver#Destroy(%v) - HostRef=%s", name, hostRef)

	//Now we need to verify the this instance of netappdvp is aware of this volume, and if it isn't then we need to make it aware
	error1 := d.storage.VerifyVolumeExists(name)
	if error1 != nil {
		return fmt.Errorf("Error - volume with name %s doesn't exist on array! error1=%s", name, error1)
	}

	//Verify that volume is mapped to this host already
	isMapped, _, error2 := d.storage.IsVolumeAlreadyMappedToHost(name, hostRef)
	if error2 != nil {
		return fmt.Errorf("IsVolumeAlreadyMappedToHost returned an error meaning the volume %s is already mapped to a different host! error2=%s", name, error2)
	}

	//Make sure the volume is mapped to this host so we can unmap it
	if !isMapped {
		log.Debugf("WARNING - volume %s is not already mapped to this host! Therefore we cannot unmap it! But we can destroy it...", name)
	} else {
		//TODO - make sure no running Docker container is currently using the volume we are about to destroy (perhaps Docker already does this check for us?)

		//It is mapped to host so we need to unmap it!
		errUnmap := d.storage.UnmapVolume(name)
		if errUnmap != nil {
			return fmt.Errorf("UnmapVolume returned an error for volume %s! errUnmap=%s", name, errUnmap)
		}

		// perform rediscovery to remove the deleted LUN
		utils.MultipathFlush() // flush unused paths
		utils.IscsiRescan()
	}

	//Destroy volume on storage array
	errDestroy := d.storage.DestroyVolume(name)
	if errDestroy != nil {
		return fmt.Errorf("DestroyVolume returned an error for volume %s! errDestroy=%s", name, errDestroy)
	}

	return nil
}

// Attach the lun
func (d *ESeriesStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("ESeriesStorageDriver#Attach(%v, %v, %v)", name, mountpoint, opts)

	//First lets get the IQN name for our host
	iqns, errIqn := utils.GetInitiatorIqns()
	if errIqn != nil {
		return fmt.Errorf("Problem determining host initiator iqns error: %v", errIqn)
	}

	log.Debugf("ESeriesStorageDriver#Attach(%v, %v, %v) - iqn=%s", name, mountpoint, opts, iqns[0]) //Going to assume a single IQN name for our host right now

	hostRef, error := d.storage.VerifyHostIQN(iqns[0])
	if error != nil {
		return fmt.Errorf("Host IQN (%s) not found on target E-Series array! error=%s", errIqn, error)
	}

	log.Debugf("ESeriesStorageDriver#Attach(%v, %v, %v) - HostRef=%s", name, mountpoint, opts, hostRef)

	error1 := d.storage.VerifyVolumeExists(name)
	if error1 != nil {
		return fmt.Errorf("Error - volume with name %s doesn't exist on array! error1=%s", name, error1)
	}

	//Variable that will hold the LUN number that our volume is mapped to on this host
	var volumeLunNumber int = -1

	isMapped, lunNumber, error2 := d.storage.IsVolumeAlreadyMappedToHost(name, hostRef)
	if error2 != nil {
		return fmt.Errorf("IsVolumeAlreadyMappedToHost returned an error meaning the volume %s is already mapped to a different host! error2=%s", name, error2)
	}

	//Set the volume LUN number to the already mapped value
	volumeLunNumber = lunNumber

	//Map the volume to host only if it isn't already mapped to host as well as no other hosts
	if !isMapped {
		//Now that we have verified that the host exists on the array we are ready to map the volume to the host only if the volume is not already mapped to host
		tmpLunNumber, error3 := d.storage.MapVolume(name, hostRef)
		if error3 != nil {
			return fmt.Errorf("Error while mapping volume to host! name=%s hostRef=%s iqn=%s error2=%s", name, hostRef, iqns[0], error3)
		}

		//Set the volume LUN number to newly mapped LUN
		volumeLunNumber = tmpLunNumber
	}

	if volumeLunNumber == -1 {
		panic("volumeLunNumber = -1!")
	}

	log.Debugf("ESeriesStorageDriver#Attach(%v, %v, %v) - volumeLunNumber=%v", name, mountpoint, opts, volumeLunNumber)

	//At this point we have our volume mapped to host so lets rescan the SCSI bus so host sees it
	rescanErr := utils.IscsiRescan()
	if rescanErr != nil {
		return rescanErr
	}

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
		if e.PortalIP == d.config.HostDataIP {
			sessionInfoToUse = sessionInfo[i]
		}
	}

	var deviceToUse = d.findDevice(volumeLunNumber, sessionInfoToUse, info)
	if deviceToUse == nil {
		return fmt.Errorf("Could not determine device to use for volume: %v ", name)
	}

	var deviceRef string = deviceToUse.Device
	if deviceToUse.MultipathDevice != "" {
		deviceRef = deviceToUse.MultipathDevice
	}

	// put a filesystem on it if there isn't one already there
	if deviceToUse.Filesystem == "" {
		// format it
		err := utils.FormatVolume(deviceRef, "ext4") // TODO externalize fsType
		if err != nil {
			return fmt.Errorf("Problem formatting lun: %v device: %v error: %v", name, deviceToUse, err)
		}
	}

	// make sure we use the proper device (multipath if in use)

	// mount it
	err := utils.Mount(deviceRef, mountpoint)
	if err != nil {
		return fmt.Errorf("Problem mounting lun: %v device: %v mountpoint: %v error: %v", name, deviceToUse, mountpoint, err)
	}

	return nil
}

func (d *ESeriesStorageDriver) findDevice(volumeLunNumber int, sessionInfo utils.IscsiSessionInfo, devices []utils.ScsiDeviceInfo) *utils.ScsiDeviceInfo {
	log.Debugf("ESeriesStorageDriver#findDevice(%v,...)", volumeLunNumber)
	// look for the expected mapped lun
	for i, device := range devices {

		log.WithFields(log.Fields{
			"i": i,
			"device": device,
		}).Debug("Checking")

		if device.LUN != strconv.Itoa(volumeLunNumber) {
			log.Debugf("Skipping... lun id %v != %v", device.LUN, volumeLunNumber)
			continue
		}

		if !strings.HasPrefix(device.IQN, sessionInfo.TargetName) {
			log.Debugf("Skipping... %v doesn't start with %v", device.IQN, sessionInfo.TargetName)
			continue
		}

		// if we're here then, we should be on the right info element:
		// *) we have the expected LUN ID
		// *) we have the expected iscsi session target
		log.Debugf("Using... %v", device)

		return &device
	}
	return nil
}

// Detach the volume
func (d *ESeriesStorageDriver) Detach(name, mountpoint string) error {
	log.Debugf("ESeriesStorageDriver#Detach(%v, %v)", name, mountpoint)

	cmd := fmt.Sprintf("umount %s", mountpoint)
	log.Debugf("cmd==%s", cmd)
	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.Debugf("out==%v", string(out))
		return fmt.Errorf("Problem unmounting docker volume: %v mountpoint: %v error: %v", name, mountpoint, err)
	}

	return nil
}

// DefaultStoragePrefix is the driver specific prefix for created storage, can be overridden in the config file
func (d *ESeriesStorageDriver) DefaultStoragePrefix() string {
	return "netappdvp_"
}
