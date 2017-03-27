// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

func newNfsTestConfig() (c *OntapStorageDriverConfig) {
	c = &OntapStorageDriverConfig{}
	c.Username = "user"
	c.Password = "password"
	c.ManagementLIF = "1.2.3.4"
	c.DataLIF = "1.2.3.5"
	c.SVM = "svm"
	c.Aggregate = "aggr"
	return
}

func TestStorageDrivers(t *testing.T) {
	log.Debug("Running TestStorageDrivers...")

	if Drivers == nil {
		t.Error("Expected storageDrivers to be a valid object")
	}

	// make sure we can find our storage driver objects
	if len(Drivers) < 2 {
		t.Error("Expected to have at least OntapNAS and OntapSAN in the list of storage drivers")
	}
}
