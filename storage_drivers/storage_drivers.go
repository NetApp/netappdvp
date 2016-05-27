// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CurrentDriverVersion is the expected version in the config file
const CurrentDriverVersion = 1

// CommonStorageDriverConfig holds settings in common across all StorageDrivers
type CommonStorageDriverConfig struct {
	Version           int    `json:"version"`
	StorageDriverName string `json:"storageDriverName"`
	Debug             bool   `json:"debug"`
	DisableDelete     bool   `json:"disableDelete"`
}

// ValidateCommonSettings attempts to "partially" decode the JSON into just the settings in CommonStorageDriverConfig
func ValidateCommonSettings(configJSON string) (*CommonStorageDriverConfig, error) {
	config := &CommonStorageDriverConfig{}

	// decode configJSON into config object
	decoder := json.NewDecoder(strings.NewReader(configJSON))
	err := decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	// load storage drivers and validate the one specified actually exists
	if config.StorageDriverName == "" {
		return nil, fmt.Errorf("Missing storage driver name in configuration file")
	}

	// validate config file version information
	if config.Version != CurrentDriverVersion {
		return nil, fmt.Errorf("Unexpected config file version;  found %v expected %v", config.Version, CurrentDriverVersion)
	}

	return config, nil
}

// OntapStorageDriverConfig holds settings for OntapStorageDrivers
type OntapStorageDriverConfig struct {
	CommonStorageDriverConfig        // embedded types replicate all fields
	ManagementLIF             string `json:"managementLIF"`
	DataLIF                   string `json:"dataLIF"`
	IgroupName                string `json:"igroupName"`
	SVM                       string `json:"svm"`
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	Aggregate                 string `json:"aggregate"`
}

// ESeriesStorageDriverConfig holds settings for ESeriesStorageDriver
type ESeriesStorageDriverConfig struct {
	CommonStorageDriverConfig

	//Web Proxy Services Info
	WebProxy_Hostname string `json:"webProxyHostname"`
	Username          string `json:"username"` //rw
	Password          string `json:"password"` //rw

	//Array Info
	Controller_A     string `json:"controllerA"`
	Controller_B     string `json:"controllerB"`
	Password_Array   string `json:"passwordArray"`   //optional
	Array_Registered bool   `json:"arrayRegistered"` //optional

	//Host Networking
	HostData_IP string `json:"hostData_IP"` //for iSCSI can be either port if multipathing is setup
}

// Drivers is a map of driver names -> object
var Drivers = make(map[string]StorageDriver)

// StorageDriver provides a common interface for storage related operations
type StorageDriver interface {
	Name() string
	Initialize(string) error
	Validate() error
	Create(name string, opts map[string]string) error
	Destroy(name string) error
	Attach(name, mountpoint string, opts map[string]string) error
	Detach(name, mountpoint string) error
}
