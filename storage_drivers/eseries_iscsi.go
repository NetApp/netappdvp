// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/eseries"
	"github.com/netapp/netappdvp/utils"
)

func init() {
	san := &ESeriesStorageDriver{}
	san.Initialized = false
	Drivers[san.Name()] = san
	log.Debugf("Registered driver '%v'", san.Name())
}

// EseriesIscsiStorageDriverName is the name for this storage driver that is specified in the config file, etc.
const EseriesIscsiStorageDriverName = "eseries-iscsi"
const DefaultAccessGroupName = "netappdvp"
const DefaultHostType = "linux_dm_mp"

// ESeriesStorageDriver is for storage provisioning via the Web Services Proxy RESTful interface that communicates
// with E-Series controllers via the SYMbol API.
type ESeriesStorageDriver struct {
	Initialized bool
	Config      ESeriesStorageDriverConfig
	API         *eseries.ESeriesAPIDriver
}

func (d *ESeriesStorageDriver) Name() string {
	return EseriesIscsiStorageDriverName
}

func (d *ESeriesStorageDriver) Protocol() string {
	return "iscsi"
}

// Initialize from the provided config
func (d *ESeriesStorageDriver) Initialize(configJSON string, commonConfig *CommonStorageDriverConfig) error {

	// Trace logging hasn't been set up yet, so always do it here
	fields := log.Fields{
		"Method": "Initialize",
		"Type":   "ESeriesStorageDriver",
	}
	log.WithFields(fields).Debug(">>>> Initialize")
	defer log.WithFields(fields).Debug("<<<< Initialize")

	config := &ESeriesStorageDriverConfig{}
	config.CommonStorageDriverConfig = commonConfig

	// Decode configJSON into ESeriesStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return fmt.Errorf("Could not decode JSON configuration. %v", err)
	}

	// Apply config defaults
	if config.StoragePrefix == nil {
		prefix := DefaultStoragePrefix
		config.StoragePrefix = &prefix
	}
	if config.AccessGroup == "" {
		config.AccessGroup = DefaultAccessGroupName
	}
	if config.HostType == "" {
		config.HostType = DefaultHostType
	}
	if config.PoolNameSearchPattern == "" {
		config.PoolNameSearchPattern = ".+"
	}

	// Fix poorly-chosen config key
	if config.HostData_IP != "" && config.HostDataIP == "" {
		config.HostDataIP = config.HostData_IP
	}

	log.WithFields(log.Fields{
		"Version":           config.Version,
		"StorageDriverName": config.StorageDriverName,
		"DebugTraceFlags":   config.DebugTraceFlags,
		"DisableDelete":     config.DisableDelete,
		"StoragePrefix":     *config.StoragePrefix,
	}).Debug("Reparsed into ESeriesStorageDriverConfig")

	d.Config = *config

	// Ensure the config is valid
	err = d.Validate()
	if err != nil {
		return fmt.Errorf("Could not validate ESeriesStorageDriver config. %v", err)
	}

	d.API = eseries.NewDriver(eseries.DriverConfig{
		WebProxyHostname:      config.WebProxyHostname,
		WebProxyPort:          config.WebProxyPort,
		WebProxyUseHTTP:       config.WebProxyUseHTTP,
		WebProxyVerifyTLS:     config.WebProxyVerifyTLS,
		Username:              config.Username,
		Password:              config.Password,
		ControllerA:           config.ControllerA,
		ControllerB:           config.ControllerB,
		PasswordArray:         config.PasswordArray,
		PoolNameSearchPattern: config.PoolNameSearchPattern,
		HostDataIP:            config.HostDataIP,
		Protocol:              d.Protocol(),
		AccessGroup:           config.AccessGroup,
		HostType:              config.HostType,
		DriverName:            config.CommonStorageDriverConfig.StorageDriverName,
		DriverVersion:         DriverVersion,
		ConfigVersion:         config.CommonStorageDriverConfig.Version,
		DebugTraceFlags:       config.CommonStorageDriverConfig.DebugTraceFlags,
	})

	// Make sure this host is logged into the E-series iSCSI target
	err = utils.EnsureIscsiSession(d.Config.HostDataIP)
	if err != nil {
		return fmt.Errorf("Could not establish iSCSI session. %v", err)
	}

	// Connect to web services proxy
	_, err = d.API.Connect()
	if err != nil {
		return fmt.Errorf("Could not connect to Web Services Proxy. %v", err)
	}

	d.Initialized = true

	return nil
}

// Validate the driver configuration
func (d *ESeriesStorageDriver) Validate() error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Validate",
			"Type":   "ESeriesStorageDriver",
		}
		log.WithFields(fields).Debug(">>>> Validate")
		defer log.WithFields(fields).Debug("<<<< Validate")
	}

	// Make sure the essential information was specified in the config
	if d.Config.WebProxyHostname == "" {
		return errors.New("WebProxyHostname is empty! You must specify the host/IP for the Web Services Proxy.")
	}
	if d.Config.ControllerA == "" || d.Config.ControllerB == "" {
		return errors.New("ControllerA or ControllerB are empty! You must specify the host/IP for the E-Series storage array. " +
			"If it is a simplex array just specify the same host/IP twice.")
	}
	if d.Config.HostDataIP == "" {
		return errors.New("HostDataIP is empty! You need to specify at least one of the iSCSI interface IP addresses that " +
			"is connected to the E-Series array.")
	}

	return nil
}

// Create is called by Docker to create a container volume. Besides the volume name, a few optional parameters such as size
// and disk media type may be provided in the opts map. If more than one pool on the storage controller can satisfy the request, the
// one with the most free space is selected.
func (d *ESeriesStorageDriver) Create(name string, sizeBytes uint64, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Create",
			"Type":   "ESeriesStorageDriver",
			"name":   name,
			"opts":   opts,
		}
		log.WithFields(fields).Debug(">>>> Create")
		defer log.WithFields(fields).Debug("<<<< Create")
	}

	// Get media type, or default to "hdd" if not specified
	mediaType := utils.GetV(opts, "mediaType", "hdd")

	// Get pool name, or default to all pools if not specified
	poolName := utils.GetV(opts, "pool", "")

	pools, err := d.API.GetVolumePools(mediaType, sizeBytes, poolName)
	if err != nil {
		return fmt.Errorf("Create failed. %v", err)
	} else if len(pools) == 0 {
		return errors.New("Create failed. No storage pools matched specified parameters.")
	}

	log.Debugf("Got pools for create: %v", pools)

	// Pick the pool with the largest free space
	sort.Sort(sort.Reverse(eseries.ByFreeSpace(pools)))
	pool := pools[0]

	// Create the volume
	vol, err := d.API.CreateVolume(name, pool.VolumeGroupRef, sizeBytes, mediaType)
	if err != nil {
		return fmt.Errorf("Could not create volume %s. %v", name, err)
	}

	log.WithFields(log.Fields{
		"Name":      name,
		"Size":      sizeBytes,
		"MediaType": mediaType,
		"VolumeRef": vol.VolumeRef,
		"Pool":      pool.Label,
	}).Debug("Create succeeded.")

	return nil
}

// Create is called by Docker to delete a container volume.
func (d *ESeriesStorageDriver) Destroy(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Destroy",
			"Type":   "ESeriesStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> Destroy")
		defer log.WithFields(fields).Debug("<<<< Destroy")
	}

	vol, err := d.API.GetVolume(name)
	if err != nil {
		return fmt.Errorf("Could not find volume %s. %v", name, err)
	}

	if d.API.IsRefValid(vol.VolumeRef) {

		// Destroy volume on storage array
		err = d.API.DeleteVolume(vol)
		if err != nil {
			return fmt.Errorf("Could not destroy volume %s. %v", name, err)
		}

	} else {

		// If volume was deleted on this storage for any reason, don't fail it here.
		log.WithField("Name", name).Warn("Could not find volume on array. Allowing deletion to proceed.")
	}

	// perform rediscovery to remove the deleted LUN
	utils.MultipathFlush() // flush unused paths
	utils.IscsiRescan()

	return nil
}

// Attach is called by Docker when attaching a container volume to a container. This method is expected to map the volume
// to the local host, discover it on the SCSI bus, format it with a filesystem, and mount it at the specified mount point.
// This method has an opts parameter, but no options are presently handled by this method.
func (d *ESeriesStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Attach",
			"Type":       "ESeriesStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
			"opts":       opts,
		}
		log.WithFields(fields).Debug(">>>> Attach")
		defer log.WithFields(fields).Debug("<<<< Attach")
	}

	// Get the volume
	vol, err := d.API.GetVolume(name)
	if err != nil {
		return fmt.Errorf("Could not find volume %s. %v", name, err)
	}
	if !d.API.IsRefValid(vol.VolumeRef) {
		return fmt.Errorf("Could not find volume %s.", name)
	}

	// Map the volume to the local host
	mapping, err := d.MapVolumeToLocalHost(vol)
	if err != nil {
		return fmt.Errorf("Could not map volume %s. %v", name, err)
	}

	// Rescan the SCSI bus to ensure the host sees the LUN
	err = utils.IscsiRescan()
	if err != nil {
		return fmt.Errorf("Could not rescan the SCSI bus. %v", err)
	}

	// Get the SCSI device information
	deviceInfo, err := utils.GetDeviceInfoForLuns()
	if err != nil {
		return fmt.Errorf("Could not get SCSI device information. %v", err)
	}

	// Get the iSCSI session information
	sessionInfo, err := utils.GetIscsiSessionInfo()
	if err != nil {
		return fmt.Errorf("Could not get iSCSI session information. %v", err)
	}

	sessionInfoToUse := utils.IscsiSessionInfo{}
	for i, e := range sessionInfo {
		if e.PortalIP == d.Config.HostDataIP {
			sessionInfoToUse = sessionInfo[i]
			break
		}
	}
	if sessionInfoToUse.TargetName == "" {
		return errors.New("Could not get iSCSI session information.")
	}

	deviceToUse := d.findDevice(mapping.LunNumber, sessionInfoToUse, deviceInfo)
	if deviceToUse.Device == "" {
		return fmt.Errorf("Could not determine device to use for volume %s.", vol.Label)
	}

	deviceRef := deviceToUse.Device
	if deviceToUse.MultipathDevice != "" {
		deviceRef = deviceToUse.MultipathDevice
	}

	// Put a filesystem on the volume if there isn't one already there
	if deviceToUse.Filesystem == "" {
		err := utils.FormatVolume(deviceRef, "ext4") // TODO externalize fsType
		if err != nil {
			return fmt.Errorf("Could not format volume %s, device %v. %v", name, deviceToUse, err)
		}
	}

	// Mount the volume
	err = utils.Mount(deviceRef, mountpoint)
	if err != nil {
		return fmt.Errorf("Could not mount volume %s, device %v at mount point %s. %v", name, deviceToUse, mountpoint, err)
	}

	return nil
}

// MapVolumeToLocalHost gets the iSCSI identity of the local host, ensures a corresponding Host definition exists on the array
// (defining a Host & HostGroup if not), maps the specified volume to the host/group (if it isn't already), and returns the mapping info.
func (d *ESeriesStorageDriver) MapVolumeToLocalHost(volume eseries.VolumeEx) (eseries.LUNMapping, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "MapVolumeToLocalHost",
			"Type":   "ESeriesStorageDriver",
			"volume": volume.Label,
		}
		log.WithFields(fields).Debug(">>>> MapVolumeToLocalHost")
		defer log.WithFields(fields).Debug("<<<< MapVolumeToLocalHost")
	}

	// Get the IQN for this host
	iqns, err := utils.GetInitiatorIqns()
	if err != nil {
		return eseries.LUNMapping{}, fmt.Errorf("Could not determine host initiator IQNs. %v", err)
	}
	if len(iqns) == 0 {
		return eseries.LUNMapping{}, errors.New("Could not determine host initiator IQNs.")
	}
	iqn := iqns[0]

	// Ensure we have an E-series host to which to map the volume
	host, err := d.API.EnsureHostForIQN(iqn)
	if err != nil {
		return eseries.LUNMapping{}, fmt.Errorf("Could not define array host for IQN %s. %v", iqn, err)
	}

	// Map the volume
	mapping, err := d.API.MapVolume(volume, host)
	if err != nil {
		return eseries.LUNMapping{}, fmt.Errorf("Could not map volume %s to host %s. %v", volume.Label, host.Label, err)
	}

	return mapping, nil
}

// findDevice combs through a list of SCSI devices to find the one matching the specified LUN number and iSCSI session info. If no
// match is found, this method returns an empty structure, so the caller should check for empty values in the result.
func (d *ESeriesStorageDriver) findDevice(
	volumeLunNumber int, sessionInfo utils.IscsiSessionInfo, devices []utils.ScsiDeviceInfo) utils.ScsiDeviceInfo {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":          "findDevice",
			"Type":            "ESeriesStorageDriver",
			"volumeLunNumber": volumeLunNumber,
			"sessionInfo":     sessionInfo,
		}
		log.WithFields(fields).Debug(">>>> findDevice")
		defer log.WithFields(fields).Debug("<<<< findDevice")
	}

	// Look for the expected mapped LUN
	for i, device := range devices {

		log.WithFields(log.Fields{"i": i, "device": device}).Debug("Checking device.")

		// LUN number must match
		if device.LUN != strconv.Itoa(volumeLunNumber) {
			log.WithFields(log.Fields{"lunID": device.LUN}).Debug("Skipping device, LUN ID does not match.")
			continue
		}

		// Target IQN must match
		if !strings.HasPrefix(device.IQN, sessionInfo.TargetName) {
			log.WithFields(log.Fields{"IQN": device.IQN}).Debug("Skipping device, IQN does not match.")
			continue
		}

		log.WithFields(log.Fields{"Device": device}).Debug("Using device.")

		return device
	}
	return utils.ScsiDeviceInfo{}
}

// Attach is called by Docker when detaching a container volume from a container. This method merely unmounts the volume; it does
// not rescan the bus, unmap the volume, or undo any of the other actions taken by the Attach method.
func (d *ESeriesStorageDriver) Detach(name, mountpoint string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Detach",
			"Type":       "ESeriesStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
		}
		log.WithFields(fields).Debug(">>>> Detach")
		defer log.WithFields(fields).Debug("<<<< Detach")
	}

	cmd := fmt.Sprintf("umount %s", mountpoint)

	log.WithFields(log.Fields{"Command": cmd}).Debug("Unmounting volume")

	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.WithFields(log.Fields{"result": string(out)}).Debug("Unmount failed.")
		return fmt.Errorf("Could not unmount docker volume: %v mountpoint: %v error: %v", name, mountpoint, err)
	}

	return nil
}

// SnapshotList returns the list of snapshots associated with the named volume. The E-series volume plugin does not support snapshots,
// so this method always returns an empty array.
func (d *ESeriesStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "SnapshotList",
			"Type":   "ESeriesStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> SnapshotList")
		defer log.WithFields(fields).Debug("<<<< SnapshotList")
	}

	return make([]CommonSnapshot, 0), nil
}

// CreateClone creates a new volume from the named volume, either by direct clone or from the named snapshot. The E-series volume plugin
// does not support cloning or snapshots, so this method always returns an error.
func (d *ESeriesStorageDriver) CreateClone(name, source, snapshot string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":   "CreateClone",
			"Type":     "ESeriesStorageDriver",
			"name":     name,
			"source":   source,
			"snapshot": snapshot,
		}
		log.WithFields(fields).Debug(">>>> CreateClone")
		defer log.WithFields(fields).Debug("<<<< CreateClone")
	}

	return errors.New("Cloning with E-Series is not supported.")
}

// Return the list of volumes associated with this tenant
func (d *ESeriesStorageDriver) List() ([]string, error) {
	prefix := *d.Config.StoragePrefix

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "List",
			"Type":   "ESeriesStorageDriver",
			"prefix": prefix,
		}
		log.WithFields(fields).Debug(">>>> List")
		defer log.WithFields(fields).Debug("<<<< List")
	}

	volumeNames, err := d.API.ListVolumes()
	if err != nil {
		return nil, fmt.Errorf("Could not get the list of volumes: %v", err)
	}

	// Filter out internal volumes
	filteredVolumeNames := make([]string, 0, len(volumeNames))
	repos_regex, _ := regexp.Compile("^repos_\\d{4}$")
	for _, name := range volumeNames {
		if !repos_regex.MatchString(name) {
			filteredVolumeNames = append(filteredVolumeNames, name)
		}
	}

	if len(prefix) == 0 {

		// No prefix, so just return the whole list
		log.WithField("Count", len(filteredVolumeNames)).Debug("Returning list of all volume names.")
		return filteredVolumeNames, nil

	} else {

		// Return only the volume names with the specified prefix
		prefixedVolumeNames := make([]string, 0, len(filteredVolumeNames))
		for _, name := range filteredVolumeNames {

			if !strings.HasPrefix(name, prefix) {
				continue
			}

			// The prefix shouldn't be visible to the user
			prefixedVolumeNames = append(prefixedVolumeNames, strings.TrimPrefix(name, prefix))
		}

		log.WithFields(log.Fields{
			"Count":  len(prefixedVolumeNames),
			"Prefix": prefix,
		}).Debug("Returning list of prefixed volume names.")
		return prefixedVolumeNames, nil
	}
}

// Test for the existence of a volume
func (d *ESeriesStorageDriver) Get(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Get",
			"Type":   "ESeriesStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> Get")
		defer log.WithFields(fields).Debug("<<<< Get")
	}

	vol, err := d.API.GetVolume(name)
	if err != nil {
		return fmt.Errorf("Could not find volume %s. %v", name, err)
	} else if !d.API.IsRefValid(vol.VolumeRef) {
		return fmt.Errorf("Could not find volume %s.", name)
	}

	return nil
}
