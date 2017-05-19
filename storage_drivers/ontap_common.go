// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
)

type OntapStorageDriver interface {
	GetConfig() *OntapStorageDriverConfig
	GetAPI() *ontap.Driver
	Name() string
}

// InitializeOntapDriver will attempt to derive the SVM to use if not provided
func InitializeOntapDriver(config OntapStorageDriverConfig) (*ontap.Driver, error) {
	api := ontap.NewDriver(ontap.DriverConfig{
		ManagementLIF: config.ManagementLIF,
		SVM:           config.SVM,
		Username:      config.Username,
		Password:      config.Password,
	})

	ontapi, err := api.SystemGetOntapiVersion()
	if err != nil {
		return nil, fmt.Errorf("Could not determine Data ONTAP API version. %v", err)
	}
	if !api.SupportsApiFeature(ontap.MINIMUM_ONTAPI_VERSION) {
		return nil, errors.New("Data ONTAP 8.3 or later is required.")
	}
	log.WithField("Ontapi", ontapi).Debug("Data ONTAP API version.")

	if config.SVM != "" {
		log.Debugf("Using specified SVM: %v", config.SVM)
		return api, nil
	}

	// use VserverGetIterRequest to populate config.SVM if it wasn't specified and we can derive it
	response1, err := api.VserverGetIterRequest()
	if !isPassed(response1.Result.ResultStatusAttr) || err != nil {
		return nil, fmt.Errorf("Error enumerating SVMs:  status: %v error: %v", response1.Result.ResultStatusAttr, err)
	}
	if response1.Result.NumRecords() != 1 {
		return nil, errors.New("Cannot derive SVM to use, please specify SVM in config file.")
	}

	// update everything to use our derived svm
	config.SVM = response1.Result.AttributesList()[0].VserverName()
	api = ontap.NewDriver(ontap.DriverConfig{
		ManagementLIF: config.ManagementLIF,
		SVM:           config.SVM,
		Username:      config.Username,
		Password:      config.Password,
	})
	api.SystemGetOntapiVersion()
	log.Debugf("Using derived SVM: %v", config.SVM)
	return api, nil
}

const DefaultSpaceReserve = "none"
const DefaultSnapshotPolicy = "none"
const DefaultUnixPermissions = "---rwxrwxrwx"
const DefaultSnapshotDir = "false"
const DefaultExportPolicy = "default"
const DefaultSecurityStyle = "unix"
const DefaultNfsMountOptions = "-o nfsvers=3"

// PopulateConfigurationDefaults fills in default values for configuration settings if not supplied in the config file
func PopulateConfigurationDefaults(config *OntapStorageDriverConfig) error {
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

	return nil
}

// EmsInitialized logs an ASUP message that this docker volume plugin has been initialized
// view them via filer::> event log show -severity NOTICE
func EmsInitialized(driverName string, api *ontap.Driver, config *OntapStorageDriverConfig) {

	// log an informational message when this plugin starts
	myHostname, hostlookupErr := os.Hostname()
	if hostlookupErr != nil {
		log.Warnf("problem while looking up hostname, error: %v", hostlookupErr)
		myHostname = "unknown"
	}

	message := driverName + " docker volume plugin initialized, version " + FullDriverVersion + " [" + ExtendedDriverVersion + "] build " + BuildVersion
	_, emsErr := api.EmsAutosupportLog(strconv.Itoa(ConfigVersion), false, "initialized", myHostname,
		message,
		1, "netappdvp", 5)
	if emsErr != nil {
		log.Warnf("problem while logging ems message, error: %v", emsErr)
	}
}

// EmsHeartbeat logs an ASUP message on a timer
// view them via filer::> event log show -severity NOTICE
func EmsHeartbeat(driverName string, api *ontap.Driver, config *OntapStorageDriverConfig) {

	// log an informational message on a timer
	myHostname, hostlookupErr := os.Hostname()
	if hostlookupErr != nil {
		log.Warnf("problem while looking up hostname, error: %v", hostlookupErr)
		myHostname = "unknown"
	}

	message := driverName + " docker volume plugin, version " + FullDriverVersion + " [" + ExtendedDriverVersion + "] build " +
		BuildVersion + " SVM[" + config.SVM + "] StoragePrefix[" + *config.StoragePrefix + "]"

	_, emsErr := api.EmsAutosupportLog(strconv.Itoa(ConfigVersion), false, "heartbeat", myHostname,
		message,
		1, "netappdvp", 5)
	if emsErr != nil {
		log.Warnf("problem while logging ems message, error: %v", emsErr)
	}
}

const MSEC_PER_HOUR = 1000 * 60 * 60 // millis * seconds * minutes

func StartEmsHeartbeat(driverName string, api *ontap.Driver, config *OntapStorageDriverConfig) {

	heartbeatIntervalInHours := 24.0 // default to 24 hours
	if config.UsageHeartbeat != "" {
		f, errParsing := strconv.ParseFloat(config.UsageHeartbeat, 64)
		if errParsing != nil {
			log.Warnf("Problem parsing heartbeat interval: {%v} error: %v", config.UsageHeartbeat, errParsing)
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
func CreateOntapClone(name, source, snapshot string, api *ontap.Driver) error {
	log.Debugf("OntapCommon#CreateOntapClone(%v, %v, %v)", name, source, snapshot)

	// If the specified volume already exists, return an error
	response, err := api.VolumeSize(name)
	if err != nil {
		return fmt.Errorf("Error searching for existing volume: error: %v", err)
	}
	if isPassed(response.Result.ResultStatusAttr) {
		return fmt.Errorf("Volume %s already exists", name)
	}

	// If no specific snapshot was requested, create one
	if snapshot == "" {
		// This is golang being stupid: https://golang.org/pkg/time/#Time.Format
		snapshot = time.Now().UTC().Format("20060102T150405Z")
		response, err := api.SnapshotCreate(snapshot, source)
		if !isPassed(response.Result.ResultStatusAttr) || err != nil {
			return fmt.Errorf("Error creating snapshot: status: %v error: %v", response.Result.ResultStatusAttr, err)
		}
	}

	// Create the clone based on a snapshot
	response2, err2 := api.VolumeCloneCreate(name, source, snapshot)
	if !isPassed(response2.Result.ResultStatusAttr) || err2 != nil {
		if response2.Result.ResultErrnoAttr == azgo.EOBJECTNOTFOUND {
			return fmt.Errorf("Snapshot %s does not exist in volume %s", snapshot, source)
		} else {
			return fmt.Errorf("Error creating clone: status: %v error: %v", response2.Result.ResultStatusAttr, err2)
		}
	}

	// Mount the new volume
	response3, err3 := api.VolumeMount(name, "/"+name)
	if !isPassed(response3.Result.ResultStatusAttr) || err3 != nil {
		return fmt.Errorf("Error mounting volume to junction: status: %v error: %v", response3.Result.ResultStatusAttr, err3)
	}

	return nil
}

// Return the list of snapshots associated with the named volume
func GetSnapshotList(name string, api *ontap.Driver) ([]CommonSnapshot, error) {
	log.Debugf("OntapCommon#GetSnapshotList(%v)", name)

	response, err := api.SnapshotGetByVolume(name)
	if !isPassed(response.Result.ResultStatusAttr) || err != nil {
		return nil, fmt.Errorf("Error enumerating snapshots: status: %v error: %v", response.Result.ResultStatusAttr, err)
	}

	log.Debugf("Returned %v snapshots", response.Result.NumRecords())
	var snapshots []CommonSnapshot

	// AttributesList() returns []SnapshotInfoType
	for _, sit := range response.Result.AttributesList() {
		log.Debugf("Snapshot name: %v, date: %v", sit.Name(), sit.AccessTime())
		t := time.Unix(int64(sit.AccessTime()), 0)
		// Time format: yyyy-mm-ddThh:mm:ssZ
		tstr := t.UTC().Format("2006-01-02T15:04:05Z")
		snapshots = append(snapshots, CommonSnapshot{sit.Name(), tstr})
	}

	return snapshots, nil
}

// Return the list of volumes associated with the tenant
func GetVolumeList(api *ontap.Driver, config *OntapStorageDriverConfig) ([]string, error) {
	log.Debugf("OntapCommon#GetVolumeList()")

	prefix := *config.StoragePrefix

	response, err := api.VolumeList(prefix)
	if !isPassed(response.Result.ResultStatusAttr) || err != nil {
		return nil, fmt.Errorf("Error enumerating volumes: status: %v error: %v", response.Result.ResultStatusAttr, err)
	}

	var volumes []string

	// AttributesList() returns []VolumeAttributesType
	for _, vat := range response.Result.AttributesList() {
		vid := vat.VolumeIdAttributes()
		vol := string(vid.Name())[len(prefix):]
		volumes = append(volumes, vol)
	}

	return volumes, nil
}

// Test for the existence of a volume
func GetVolume(name string, api *ontap.Driver) error {
	response, err := api.VolumeSize(name)
	if !isPassed(response.Result.ResultStatusAttr) || err != nil {
		return fmt.Errorf("Error searching for existing volume: status: %v error: %v", response.Result.ResultStatusAttr, err)
	}

	return nil
}

func isPassed(s string) bool {
	const passed = "passed"
	return s == passed
}
