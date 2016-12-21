// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

// encodeConfig takes the supplied config object, encodes to json, and returns the result
func encodeConfig(c *CommonStorageDriverConfig) (string, error) {
	sb := bytes.NewBufferString("")
	encoder := json.NewEncoder(sb)
	if err := encoder.Encode(c); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func TestOntap_ValidateCommonSettings(t *testing.T) {
	log.Debug("Running storage_drivers.TestOntap_ValidateCommonSettings...")

	badConfig, err := ValidateCommonSettings("")
	if err == nil || badConfig != nil {
		t.Errorf("Exepcted an error for an empty JSON configuration object")
	}

	// test with missing storage driver name
	configMissingStorageDriverName := &CommonStorageDriverConfig{
		Version:           1,
		StorageDriverName: "", // emtpy on purpose
		Debug:             false,
		DisableDelete:     false,
		StoragePrefixRaw:  []byte(`""`),
		SnapshotPrefixRaw: []byte(`""`),
	}

	json, err := encodeConfig(configMissingStorageDriverName)
	_, err = ValidateCommonSettings(json)
	if !strings.HasPrefix(err.Error(), "Missing storage driver name") {
		t.Errorf("Did not validate storage driver name")
	}

	// test with unexpcted version
	configWithWrongVersion := &CommonStorageDriverConfig{
		Version:           -1, // invalid number on purpose
		StorageDriverName: OntapNASStorageDriverName,
		Debug:             false,
		DisableDelete:     false,
		StoragePrefixRaw:  []byte(`""`),
		SnapshotPrefixRaw: []byte(`""`),
	}

	json, err = encodeConfig(configWithWrongVersion)
	_, err = ValidateCommonSettings(json)
	if !strings.HasPrefix(err.Error(), "Unexpected config file version") {
		t.Errorf("Did not validate version number")
	}

	// test with an expcted good config version
	config := &CommonStorageDriverConfig{
		Version:           ConfigVersion,
		StorageDriverName: OntapNASStorageDriverName,
		Debug:             false,
		DisableDelete:     false,
		StoragePrefixRaw:  []byte(`""`),
		SnapshotPrefixRaw: []byte(`""`),
	}

	json, err = encodeConfig(config)
	_, err = ValidateCommonSettings(json)
	if err != nil || config == nil {
		t.Errorf("Expected to have a good config object, unexpected error: %v", err)
	}

	if config.Version != ConfigVersion {
		t.Errorf("Unexpcted config version found: %v", config.Version)
	}
}

func TestOntapNas_Init(t *testing.T) {
	log.Debug("Running storage_drivers.TestOntapNas_Init...")

	if Drivers == nil {
		t.Error("Expected storageDrivers to be a valid object")
	}

	if Drivers[OntapNASStorageDriverName] == nil {
		t.Error("Expected to find a valid object")
	}

	if Drivers[OntapNASStorageDriverName].Name() != OntapNASStorageDriverName {
		t.Errorf("Unexpected object found for key %v", OntapNASStorageDriverName)
	}
}

func TestOntapNas_Initialize(t *testing.T) {
	log.Debug("Running storage_drivers.TestOntapNas_Initialize...")

	nas := &OntapNASStorageDriver{}

	if err := nas.Initialize(""); err == nil {
		t.Errorf("Exepcted an error for an empty JSON configuration object")
	}
}

func TestOntapSan_Init(t *testing.T) {
	log.Debug("Running storage_drivers.TestOntapSan_Init...")

	if Drivers == nil {
		t.Error("Expected storageDrivers to be a valid object")
	}

	if Drivers[OntapSANStorageDriverName] == nil {
		t.Error("Expected to find a valid object")
	}

	if Drivers[OntapSANStorageDriverName].Name() != OntapSANStorageDriverName {
		t.Errorf("Unexpected object found for key %v", OntapSANStorageDriverName)
	}
}

func TestOntapSan_Initialize(t *testing.T) {
	log.Debug("Running storage_drivers.TestOntapSan_Initialize...")

	san := &OntapSANStorageDriver{}

	if err := san.Initialize(""); err == nil {
		t.Errorf("Exepcted an error for an empty JSON configuration object")
	}
}
