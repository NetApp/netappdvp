// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"
)

const LS_MIRROR_IDLE_TIMEOUT_SECS = 30

type OntapStorageDriver interface {
	GetConfig() *OntapStorageDriverConfig
	GetAPI() *ontap.Driver
	Name() string
}

// InitializeOntapConfig parses the ONTAP config, mixing in the specified common config.
func InitializeOntapConfig(
	configJSON string, commonConfig *CommonStorageDriverConfig,
) (*OntapStorageDriverConfig, error) {

	if commonConfig.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "InitializeOntapConfig", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> InitializeOntapConfig")
		defer log.WithFields(fields).Debug("<<<< InitializeOntapConfig")
	}

	config := &OntapStorageDriverConfig{}
	config.CommonStorageDriverConfig = commonConfig

	// decode configJSON into OntapStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return nil, fmt.Errorf("Could not decode JSON configuration. %v", err)
	}

	return config, nil
}

// InitializeOntapDriver sets up the API client and performs all other initialization tasks
// that are common to all the ONTAP drivers.
func InitializeOntapDriver(config *OntapStorageDriverConfig) (*ontap.Driver, error) {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "InitializeOntapDriver", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> InitializeOntapDriver")
		defer log.WithFields(fields).Debug("<<<< InitializeOntapDriver")
	}

	// Get the API client
	api, err := InitializeOntapAPI(config)
	if err != nil {
		return nil, fmt.Errorf("Could not create Data ONTAP API client. %v", err)
	}

	// Make sure we're using a valid ONTAP version
	ontapi, err := api.SystemGetOntapiVersion()
	if err != nil {
		return nil, fmt.Errorf("Could not determine Data ONTAP API version. %v", err)
	}
	if !api.SupportsApiFeature(ontap.MINIMUM_ONTAPI_VERSION) {
		return nil, errors.New("Data ONTAP 8.3 or later is required.")
	}
	log.WithField("Ontapi", ontapi).Debug("Data ONTAP API version.")

	// Log cluster node serial numbers if we can get them
	config.SerialNumbers, err = api.ListNodeSerialNumbers()
	if err != nil {
		log.Warnf("Could not determine controller serial numbers. %v", err)
	} else {
		log.WithField("serialNumbers", config.SerialNumbers).Info("Controller serial numbers.")
	}

	// Load default config parameters
	err = PopulateConfigurationDefaults(config)
	if err != nil {
		return nil, fmt.Errorf("Could not populate configuration defaults. %v", err)
	}

	return api, nil
}

// InitializeOntapAPI returns an ontap.Driver ZAPI client.  If the SVM isn't specified in the config
// file, this method attempts to derive the one to use.
func InitializeOntapAPI(config *OntapStorageDriverConfig) (*ontap.Driver, error) {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "InitializeOntapAPI", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> InitializeOntapAPI")
		defer log.WithFields(fields).Debug("<<<< InitializeOntapAPI")
	}

	api := ontap.NewDriver(ontap.DriverConfig{
		ManagementLIF:   config.ManagementLIF,
		SVM:             config.SVM,
		Username:        config.Username,
		Password:        config.Password,
		DebugTraceFlags: config.DebugTraceFlags,
	})

	if config.SVM != "" {
		log.WithField("SVM", config.SVM).Debug("Using specified SVM.")
		return api, nil
	}

	// Use VserverGetIterRequest to populate config.SVM if it wasn't specified and we can derive it
	vserverResponse, err := api.VserverGetIterRequest()
	if err = ontap.GetError(vserverResponse, err); err != nil {
		return nil, fmt.Errorf("Error enumerating SVMs. %v", err)
	}

	if vserverResponse.Result.NumRecords() != 1 {
		return nil, errors.New("Cannot derive SVM to use, please specify SVM in config file.")
	}

	// Update everything to use our derived SVM
	config.SVM = vserverResponse.Result.AttributesList()[0].VserverName()
	api = ontap.NewDriver(ontap.DriverConfig{
		ManagementLIF:   config.ManagementLIF,
		SVM:             config.SVM,
		Username:        config.Username,
		Password:        config.Password,
		DebugTraceFlags: config.DebugTraceFlags,
	})
	log.WithField("SVM", config.SVM).Debug("Using derived SVM.")

	return api, nil
}

// ValidateAggregate returns an error if the configured aggregate is not available to the Vserver.
func ValidateAggregate(api *ontap.Driver, config *OntapStorageDriverConfig) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "ValidateAggregate", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> ValidateAggregate")
		defer log.WithFields(fields).Debug("<<<< ValidateAggregate")
	}

	// Get the aggregates assigned to the SVM.  There must be at least one!
	vserverAggrs, err := api.GetVserverAggregateNames()
	if err != nil {
		return err
	}
	if len(vserverAggrs) == 0 {
		return fmt.Errorf("SVM %s has no assigned aggregates.", config.SVM)
	}

	for _, aggrName := range vserverAggrs {
		if aggrName == config.Aggregate {
			log.WithFields(log.Fields{
				"SVM":       config.SVM,
				"Aggregate": config.Aggregate,
			}).Debug("Found aggregate for SVM.")
			return nil
		}
	}

	return fmt.Errorf("Aggregate %s does not exist or is not assigned to SVM %s.", config.Aggregate, config.SVM)
}

// ValidateNASDriver contains the validation logic shared between ontap-nas and ontap-nas-economy.
func ValidateNASDriver(context DriverContext, api *ontap.Driver, config *OntapStorageDriverConfig) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "ValidateNASDriver", "Type": "ontap_common", "context": context}
		log.WithFields(fields).Debug(">>>> ValidateNASDriver")
		defer log.WithFields(fields).Debug("<<<< ValidateNASDriver")
	}

	lifResponse, err := api.NetInterfaceGet()
	if err = ontap.GetError(lifResponse, err); err != nil {
		return fmt.Errorf("Error checking network interfaces. %v", err)
	}

	// If they didn't set a LIF to use in the config, we'll set it to the first nfs LIF we happen to find
	if config.DataLIF == "" {
	loop:
		for _, attrs := range lifResponse.Result.AttributesList() {
			for _, protocol := range attrs.DataProtocols() {
				if protocol == "nfs" {
					log.WithField("address", attrs.Address()).Debug("Choosing LIF for NFS.")
					config.DataLIF = string(attrs.Address())
					break loop
				}
			}
		}
	}

	foundNfs := false
loop2:
	for _, attrs := range lifResponse.Result.AttributesList() {
		for _, protocol := range attrs.DataProtocols() {
			if protocol == "nfs" {
				log.Debugf("Comparing NFS protocol access on %v vs. %v", attrs.Address(), config.DataLIF)
				if string(attrs.Address()) == config.DataLIF {
					foundNfs = true
					break loop2
				}
			}
		}
	}

	if !foundNfs {
		return fmt.Errorf("Could not find NFS Data LIF.")
	}

	if context == ContextNDVP {
		// Make sure the configured aggregate is available
		err = ValidateAggregate(api, config)
		if err != nil {
			return err
		}
	}

	return nil
}

const DefaultSpaceReserve = "none"
const DefaultSnapshotPolicy = "none"
const DefaultUnixPermissions = "---rwxrwxrwx"
const DefaultSnapshotDir = "false"
const DefaultExportPolicy = "default"
const DefaultSecurityStyle = "unix"
const DefaultNfsMountOptions = "-o nfsvers=3"
const DefaultSplitOnClone = "false"
const DefaultFileSystemType = "ext4"
const DefaultEncryption = "false"

// PopulateConfigurationDefaults fills in default values for configuration settings if not supplied in the config file
func PopulateConfigurationDefaults(config *OntapStorageDriverConfig) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "PopulateConfigurationDefaults", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> PopulateConfigurationDefaults")
		defer log.WithFields(fields).Debug("<<<< PopulateConfigurationDefaults")
	}

	if config.StoragePrefix == nil {
		prefix := DefaultStoragePrefix
		config.StoragePrefix = &prefix
	}

	if config.SpaceReserve == "" {
		config.SpaceReserve = DefaultSpaceReserve
	}

	if config.SnapshotPolicy == "" {
		config.SnapshotPolicy = DefaultSnapshotPolicy
	}

	if config.UnixPermissions == "" {
		config.UnixPermissions = DefaultUnixPermissions
	}

	if config.SnapshotDir == "" {
		config.SnapshotDir = DefaultSnapshotDir
	}

	if config.ExportPolicy == "" {
		config.ExportPolicy = DefaultExportPolicy
	}

	if config.SecurityStyle == "" {
		config.SecurityStyle = DefaultSecurityStyle
	}

	if config.NfsMountOptions == "" {
		config.NfsMountOptions = DefaultNfsMountOptions
	}

	if config.SplitOnClone == "" {
		config.SplitOnClone = DefaultSplitOnClone
	} else {
		_, err := strconv.ParseBool(config.SplitOnClone)
		if err != nil {
			return fmt.Errorf("Invalid boolean value for splitOnClone. %v", err)
		}
	}

	if config.FileSystemType == "" {
		config.FileSystemType = DefaultFileSystemType
	}

	if config.Encryption == "" {
		config.Encryption = DefaultEncryption
	}

	log.WithFields(log.Fields{
		"StoragePrefix":   config.StoragePrefix,
		"SpaceReserve":    config.SpaceReserve,
		"SnapshotPolicy":  config.SnapshotPolicy,
		"UnixPermissions": config.UnixPermissions,
		"SnapshotDir":     config.SnapshotDir,
		"ExportPolicy":    config.ExportPolicy,
		"SecurityStyle":   config.SecurityStyle,
		"NfsMountOptions": config.NfsMountOptions,
		"SplitOnClone":    config.SplitOnClone,
		"FileSystemType":  config.FileSystemType,
		"Encryption":      config.Encryption,
	}).Debugf("Configuration defaults")

	return nil
}

// ValidateEncryptionAttribute returns true/false if encryption is being requested of a backend that
// supports NetApp Volume Encryption, and nil otherwise so that the ZAPIs may be sent without
// any reference to encryption.
func ValidateEncryptionAttribute(encryption string, api *ontap.Driver) (*bool, error) {

	enableEncryption, err := strconv.ParseBool(encryption)
	if err != nil {
		return nil, fmt.Errorf("Invalid boolean value for encryption. %v", err)
	}

	if api.SupportsApiFeature(ontap.NETAPP_VOLUME_ENCRYPTION) {
		return &enableEncryption, nil
	} else {
		if enableEncryption {
			return nil, errors.New("Encrypted volumes are not supported on this storage backend.")
		} else {
			return nil, nil
		}
	}
}

// EmsInitialized logs an ASUP message that this docker volume plugin has been initialized
// view them via filer::> event log show -severity NOTICE
func EmsInitialized(driverName string, api *ontap.Driver, config *OntapStorageDriverConfig) {

	// log an informational message when this plugin starts
	myHostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Could not determine hostname. %v", err)
		myHostname = "unknown"
	}

	message := driverName + " docker volume plugin initialized, version " + FullDriverVersion + " [" + ExtendedDriverVersion + "] build " + BuildVersion
	emsResponse, err := api.EmsAutosupportLog(strconv.Itoa(ConfigVersion), false, "initialized", myHostname,
		message, 1, "netappdvp", 5)
	if err = ontap.GetError(emsResponse, err); err != nil {
		log.Warnf("Error logging EMS message. %v", err)
	}
}

// EmsHeartbeat logs an ASUP message on a timer
// view them via filer::> event log show -severity NOTICE
func EmsHeartbeat(driverName string, api *ontap.Driver, config *OntapStorageDriverConfig) {

	// log an informational message on a timer
	myHostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Could not determine hostname. %v", err)
		myHostname = "unknown"
	}

	message := driverName + " docker volume plugin, version " + FullDriverVersion + " [" + ExtendedDriverVersion + "] build " +
		BuildVersion + " SVM[" + config.SVM + "] StoragePrefix[" + *config.StoragePrefix + "]"

	emsResponse, err := api.EmsAutosupportLog(strconv.Itoa(ConfigVersion), false, "heartbeat", myHostname,
		message, 1, "netappdvp", 5)
	if err = ontap.GetError(emsResponse, err); err != nil {
		log.Warnf("Error logging EMS message. %v", err)
	}
}

const MSEC_PER_HOUR = 1000 * 60 * 60 // millis * seconds * minutes

func StartEmsHeartbeat(driverName string, api *ontap.Driver, config *OntapStorageDriverConfig) {

	heartbeatIntervalInHours := 24.0 // default to 24 hours
	if config.UsageHeartbeat != "" {
		f, err := strconv.ParseFloat(config.UsageHeartbeat, 64)
		if err != nil {
			log.WithField("interval", config.UsageHeartbeat).Warnf("Invalid heartbeat interval. %v", err)
		} else {
			heartbeatIntervalInHours = f
		}
	}
	log.WithField("intervalHours", heartbeatIntervalInHours).Debug("Configured EMS heartbeat.")

	durationInHours := time.Millisecond * time.Duration(MSEC_PER_HOUR*heartbeatIntervalInHours)
	if durationInHours > 0 {
		EmsHeartbeat(driverName, api, config)
		ticker := time.NewTicker(durationInHours)
		go func() {
			for t := range ticker.C {
				log.WithField("tick", t).Debug("Sending EMS heartbeat.")
				EmsHeartbeat(driverName, api, config)
			}
		}()
	}
}

// Create a volume clone
func CreateOntapClone(
	name, source, snapshot string, split bool, config *OntapStorageDriverConfig, api *ontap.Driver,
) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":   "CreateOntapClone",
			"Type":     "ontap_common",
			"name":     name,
			"source":   source,
			"snapshot": snapshot,
			"split":    split,
		}
		log.WithFields(fields).Debug(">>>> CreateOntapClone")
		defer log.WithFields(fields).Debug("<<<< CreateOntapClone")
	}

	// If the specified volume already exists, return an error
	volExists, err := api.VolumeExists(name)
	if err != nil {
		return fmt.Errorf("Error checking for existing volume. %v", err)
	}
	if volExists {
		return fmt.Errorf("Volume %s already exists.", name)
	}

	// If no specific snapshot was requested, create one
	if snapshot == "" {
		// This is golang being stupid: https://golang.org/pkg/time/#Time.Format
		snapshot = time.Now().UTC().Format("20060102T150405Z")
		snapResponse, err := api.SnapshotCreate(snapshot, source)
		if err = ontap.GetError(snapResponse, err); err != nil {
			return fmt.Errorf("Error creating snapshot. %v", err)
		}
	}

	// Create the clone based on a snapshot
	cloneResponse, err := api.VolumeCloneCreate(name, source, snapshot)
	if err != nil {
		return fmt.Errorf("Error creating clone. %v", err)
	}
	if zerr := ontap.NewZapiError(cloneResponse); !zerr.IsPassed() {
		if zerr.Code() == azgo.EOBJECTNOTFOUND {
			return fmt.Errorf("Snapshot %s does not exist in volume %s.", snapshot, source)
		} else {
			return fmt.Errorf("Error creating clone. %v", zerr)
		}
	}

	// Mount the new volume
	mountResponse, err := api.VolumeMount(name, "/"+name)
	if err = ontap.GetError(mountResponse, err); err != nil {
		return fmt.Errorf("Error mounting volume to junction. %v", err)
	}

	// Split the clone if requested
	if split {
		splitResponse, err := api.VolumeCloneSplitStart(name)
		if err = ontap.GetError(splitResponse, err); err != nil {
			return fmt.Errorf("Error splitting clone. %v", err)
		}
	}

	return nil
}

// Return the list of snapshots associated with the named volume
func GetSnapshotList(name string, config *OntapStorageDriverConfig, api *ontap.Driver) ([]CommonSnapshot, error) {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetSnapshotList",
			"Type":   "ontap_common",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> GetSnapshotList")
		defer log.WithFields(fields).Debug("<<<< GetSnapshotList")
	}

	snapResponse, err := api.SnapshotGetByVolume(name)
	if err = ontap.GetError(snapResponse, err); err != nil {
		return nil, fmt.Errorf("Error enumerating snapshots. %v", err)
	}

	log.Debugf("Returned %v snapshots.", snapResponse.Result.NumRecords())
	snapshots := []CommonSnapshot{}

	// AttributesList() returns []SnapshotInfoType
	for _, snap := range snapResponse.Result.AttributesList() {

		log.WithFields(log.Fields{
			"name":       snap.Name(),
			"accessTime": snap.AccessTime(),
		}).Debug("Snapshot")

		// Time format: yyyy-mm-ddThh:mm:ssZ
		snapTime := time.Unix(int64(snap.AccessTime()), 0).UTC().Format("2006-01-02T15:04:05Z")

		snapshots = append(snapshots, CommonSnapshot{snap.Name(), snapTime})
	}

	return snapshots, nil
}

// Return the list of volumes associated with the tenant
func GetVolumeList(api *ontap.Driver, config *OntapStorageDriverConfig) ([]string, error) {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "GetVolumeList", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> GetVolumeList")
		defer log.WithFields(fields).Debug("<<<< GetVolumeList")
	}

	prefix := *config.StoragePrefix

	volResponse, err := api.VolumeList(prefix)
	if err = ontap.GetError(volResponse, err); err != nil {
		return nil, fmt.Errorf("Error enumerating volumes. %v", err)
	}

	var volumes []string

	// AttributesList() returns []VolumeAttributesType
	for _, volume := range volResponse.Result.AttributesList() {
		vol_id_attrs := volume.VolumeIdAttributes()
		volName := string(vol_id_attrs.Name())[len(prefix):]
		volumes = append(volumes, volName)
	}

	return volumes, nil
}

// GetVolume checks for the existence of a volume.  It returns nil if the volume
// exists and an error if it does not (or the API call fails).
func GetVolume(name string, api *ontap.Driver, config *OntapStorageDriverConfig) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "GetVolume", "Type": "ontap_common"}
		log.WithFields(fields).Debug(">>>> GetVolume")
		defer log.WithFields(fields).Debug("<<<< GetVolume")
	}

	volExists, err := api.VolumeExists(name)
	if err != nil {
		return fmt.Errorf("Error checking for existing volume. %v", err)
	}
	if !volExists {
		log.WithField("flexvol", name).Debug("Flexvol not found.")
		return fmt.Errorf("Volume %s does not exist.", name)
	}

	return nil
}

// MountVolume accepts the mount info for an NFS share and mounts it on the local host.
func MountVolume(exportPath, mountpoint string, config *OntapStorageDriverConfig) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "MountVolume",
			"Type":       "ontap_common",
			"exportPath": exportPath,
			"mountpoint": mountpoint,
		}
		log.WithFields(fields).Debug(">>>> MountVolume")
		defer log.WithFields(fields).Debug("<<<< MountVolume")
	}

	nfsMountOptions := config.NfsMountOptions

	// Do the mount
	var cmd string
	switch runtime.GOOS {
	case utils.Linux:
		cmd = fmt.Sprintf("mount -v %s %s %s", nfsMountOptions, exportPath, mountpoint)
	case utils.Darwin:
		cmd = fmt.Sprintf("mount -v -o rw %s -t nfs %s %s", nfsMountOptions, exportPath, mountpoint)
	default:
		return fmt.Errorf("Unsupported operating system: %v", runtime.GOOS)
	}

	log.WithField("command", cmd).Debug("Mounting volume.")

	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.WithField("output", string(out)).Debug("Mount failed.")
		return fmt.Errorf("Error mounting NFS volume %v on mountpoint %v. %v", exportPath, mountpoint, err)
	}

	return nil
}

// UnmountVolume unmounts the volume mounted on the specified mountpoint.
func UnmountVolume(mountpoint string, config *OntapStorageDriverConfig) error {

	if config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "UnmountVolume",
			"Type":       "ontap_common",
			"mountpoint": mountpoint,
		}
		log.WithFields(fields).Debug(">>>> UnmountVolume")
		defer log.WithFields(fields).Debug("<<<< UnmountVolume")
	}

	cmd := fmt.Sprintf("umount %s", mountpoint)
	log.WithField("command", cmd).Debug("Unmounting volume.")

	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.WithField("output", string(out)).Debug("Unmount failed.")
		return fmt.Errorf("Error unmounting NFS volume from mountpoint %v. %v", mountpoint, err)
	}

	return nil
}

// UpdateLoadSharingMirrors checks for the present of LS mirrors on the SVM root volume, and if
// present, starts an update and waits for them to become idle.
func UpdateLoadSharingMirrors(api *ontap.Driver) {

	// We care about LS mirrors on the SVM root volume, so get the root volume name
	rootVolumeResponse, err := api.VolumeGetRootName()
	if err = ontap.GetError(rootVolumeResponse, err); err != nil {
		log.Warnf("Error getting SVM root volume. %v", err)
		return
	}
	rootVolume := rootVolumeResponse.Result.Volume()

	// Check for LS mirrors on the SVM root volume
	mirrorGetResponse, err := api.SnapmirrorGetLoadSharingMirrors(rootVolume)
	if err = ontap.GetError(rootVolumeResponse, err); err != nil {
		log.Warnf("Error getting load-sharing mirrors for SVM root volume. %v", err)
		return
	}
	if mirrorGetResponse.Result.NumRecords() == 0 {
		// None found, so nothing more to do
		log.WithField("rootVolume", rootVolume).Debug("SVM root volume has no load-sharing mirrors.")
		return
	}

	// One or more LS mirrors found, so issue an update
	mirrorSourceLocation := mirrorGetResponse.Result.AttributesList()[0].SourceLocation()
	_, err = api.SnapmirrorUpdateLoadSharingMirrors(mirrorSourceLocation)
	if err = ontap.GetError(rootVolumeResponse, err); err != nil {
		log.Warnf("Error updating load-sharing mirrors for SVM root volume. %v", err)
		return
	}

	// Wait for LS mirrors to become idle
	timeout := time.Now().Add(LS_MIRROR_IDLE_TIMEOUT_SECS * time.Second)
	for {
		time.Sleep(1 * time.Second)
		log.Debug("Load-sharing mirrors not yet idle, polling...")

		mirrorGetResponse, err = api.SnapmirrorGetLoadSharingMirrors(rootVolume)
		if err = ontap.GetError(rootVolumeResponse, err); err != nil {
			log.Warnf("Error getting load-sharing mirrors for SVM root volume. %v", err)
			break
		}
		if mirrorGetResponse.Result.NumRecords() == 0 {
			log.WithField("rootVolume", rootVolume).Debug("SVM root volume has no load-sharing mirrors.")
			break
		}

		// Ensure all mirrors are idle
		idle := true
		for _, mirror := range mirrorGetResponse.Result.AttributesList() {
			if mirror.RelationshipStatusPtr == nil || mirror.RelationshipStatus() != "idle" {
				idle = false
			}
		}
		if idle {
			log.Debug("Load-sharing mirrors idle.")
			break
		}

		// Don't wait forever
		if time.Now().After(timeout) {
			log.Warning("Load-sharing mirrors not yet idle, giving up.")
			break
		}
	}
}
