// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/sfapi"
	"github.com/netapp/netappdvp/utils"
)

// ConfigVersion is the expected version specified in the config file
const ConfigVersion = 1

// DriverVersion is the actual release version number
const DriverVersion = "17.07.0"

// FullDriverVersion is the DriverVersion as well as any pre-release tags
var FullDriverVersion = DriverVersion

// BuildVersion is the extended release version with build information
var BuildVersion = "unknown"

// BuildTime is the date and time the binary was built, if known
var BuildTime = "unknown"

// ExtendedDriverVersion can be overridden by embeddors such as Trident to uniquify the version string
var ExtendedDriverVersion = "native"

// DefaultStoragePrefix can be overridden by Trident too. God.
var DefaultStoragePrefix = "netappdvp_"

// CommonStorageDriverConfig holds settings in common across all StorageDrivers
type CommonStorageDriverConfig struct {
	Version                           int             `json:"version"`
	StorageDriverName                 string          `json:"storageDriverName"`
	Debug                             bool            `json:"debug"`           // Unsupported!
	DebugTraceFlags                   map[string]bool `json:"debugTraceFlags"` // Example: {"api":false, "method":true}
	DisableDelete                     bool            `json:"disableDelete"`
	StoragePrefix                     *string         `json:"storagePrefix"`
	CommonStorageDriverConfigDefaults `json:"defaults"`
}

type CommonStorageDriverConfigDefaults struct {
	Size string `json:"size"`
}

// ValidateCommonSettings attempts to "partially" decode the JSON into just the settings in CommonStorageDriverConfig
func ValidateCommonSettings(configJSON string) (*CommonStorageDriverConfig, error) {
	config := &CommonStorageDriverConfig{}

	// decode configJSON into config object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	// load storage drivers and validate the one specified actually exists
	if config.StorageDriverName == "" {
		return nil, errors.New("Missing storage driver name in configuration file")
	}

	// validate config file version information
	if config.Version != ConfigVersion {
		return nil, fmt.Errorf("Unexpected config file version;  found %v expected %v", config.Version, ConfigVersion)
	}

	// warn if using debug in the config file
	if config.Debug {
		log.Warnf("The debug setting in the configuration file is now ignored; use the command " +
			"line --debug switch instead.")
	}

	// ensure the default volume size is valid, using a default of 1g if not set
	if config.Size == "" {
		config.Size = "1G"
		log.WithField("size", config.Size).Debug("Setting default volume size.")
	} else {
		_, err = utils.ConvertSizeToBytes(config.Size)
		if err != nil {
			return nil, fmt.Errorf("Invalid config value for default volume size: %v", err)
		}
	}

	return config, nil
}

// OntapStorageDriverConfig holds settings for OntapStorageDrivers
type OntapStorageDriverConfig struct {
	*CommonStorageDriverConfig              // embedded types replicate all fields
	ManagementLIF                    string `json:"managementLIF"`
	DataLIF                          string `json:"dataLIF"`
	IgroupName                       string `json:"igroupName"`
	SVM                              string `json:"svm"`
	Username                         string `json:"username"`
	Password                         string `json:"password"`
	Aggregate                        string `json:"aggregate"`
	UsageHeartbeat                   string `json:"usageHeartbeat"` // in hours, default to 24.0
	NfsMountOptions                  string `json:"nfsMountOptions"`
	OntapStorageDriverConfigDefaults `json:"defaults"`
}

type OntapStorageDriverConfigDefaults struct {
	SpaceReserve    string `json:"spaceReserve"`
	SnapshotPolicy  string `json:"snapshotPolicy"`
	UnixPermissions string `json:"unixPermissions"`
	SnapshotDir     string `json:"snapshotDir"`
	ExportPolicy    string `json:"exportPolicy"`
	SecurityStyle   string `json:"securityStyle"`
}

// ESeriesStorageDriverConfig holds settings for ESeriesStorageDriver
type ESeriesStorageDriverConfig struct {
	*CommonStorageDriverConfig

	// Web Proxy Services Info
	WebProxyHostname  string `json:"webProxyHostname"`
	WebProxyPort      string `json:"webProxyPort"`      // optional
	WebProxyUseHTTP   bool   `json:"webProxyUseHTTP"`   // optional
	WebProxyVerifyTLS bool   `json:"webProxyVerifyTLS"` // optional
	Username          string `json:"username"`
	Password          string `json:"password"`

	// Array Info
	ControllerA   string `json:"controllerA"`
	ControllerB   string `json:"controllerB"`
	PasswordArray string `json:"passwordArray"` //optional

	// Options
	PoolNameSearchPattern string `json:"poolNameSearchPattern"` //optional

	// Host Networking
	HostData_IP string `json:"hostData_IP,omitempty"` // for backward compatibility only
	HostDataIP  string `json:"hostDataIP"`            // for iSCSI can be either port if multipathing is setup
	AccessGroup string `json:"accessGroupName"`       // name for host group, default is 'netappdvp'
	HostType    string `json:"hostType"`              // host type, default is 'linux_dm_mp'
}

// SolidfireStorageDriverConfig holds settings for SolidfireStorageDrivers
type SolidfireStorageDriverConfig struct {
	*CommonStorageDriverConfig // embedded types replicate all fields
	TenantName                 string
	EndPoint                   string
	DefaultVolSz               int64 //Default volume size in GiB (deprecated)
	SVIP                       string
	InitiatorIFace             string //iface to use of iSCSI initiator
	Types                      *[]sfapi.VolType
	LegacyNamePrefix           string //name prefix used in earlier ndvp versions
	AccessGroups               []int64
}

// CommonSnapshot contains the normalized volume snapshot format we report to Docker
type CommonSnapshot struct {
	Name    string // The snapshot name or other identifier you would use to reference it
	Created string // The UTC time that the snapshot was created, in RFC3339 format
}

// Drivers is a map of driver names -> object
var Drivers = make(map[string]StorageDriver)

// StorageDriver provides a common interface for storage related operations
type StorageDriver interface {
	Name() string
	Initialize(string, *CommonStorageDriverConfig) error
	Validate() error
	Create(name string, sizeBytes uint64, opts map[string]string) error
	CreateClone(name, source, snapshot string) error
	Destroy(name string) error
	Attach(name, mountpoint string, opts map[string]string) error
	Detach(name, mountpoint string) error
	SnapshotList(name string) ([]CommonSnapshot, error)
	List() ([]string, error)
	Get(name string) error
}
