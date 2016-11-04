// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/netapp/netappdvp/azgo"
	"github.com/ebalduf/netappdvp/apis/ontap"

	log "github.com/Sirupsen/logrus"
)

// InitializeOntapDriver will attempt to derive the SVM to use if not provided
func InitializeOntapDriver(config OntapStorageDriverConfig) (*ontap.Driver, error) {
	api := ontap.NewDriver(ontap.DriverConfig{
		ManagementLIF: config.ManagementLIF,
		SVM:           config.SVM,
		Username:      config.Username,
		Password:      config.Password,
	})

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
		return nil, fmt.Errorf("Cannot derive SVM to use, please specify SVM in config file.")
	}

	// update everything to use our derived svm
	config.SVM = response1.Result.AttributesList()[0].VserverName()
	api = ontap.NewDriver(ontap.DriverConfig{
		ManagementLIF: config.ManagementLIF,
		SVM:           config.SVM,
		Username:      config.Username,
		Password:      config.Password,
	})
	log.Debugf("Using derived SVM: %v", config.SVM)
	return api, nil
}

// EmsInitialized logs an ASUP message that this docker volume plugin has been initialized
// view them via filer::> event log show
func EmsInitialized(driverName string, api *ontap.Driver) {

	// log an informational message when this plugin starts
	myHostname, hostlookupErr := os.Hostname()
	if hostlookupErr != nil {
		log.Warnf("problem while looking up hostname, error: %v", hostlookupErr)
		myHostname = "unknown"
	}

	_, emsErr := api.EmsAutosupportLog(strconv.Itoa(CurrentDriverVersion), false, "initialized", myHostname, driverName+" docker volume plugin initialized, version "+DriverVersion, 1, "netappdvp", 6)
	if emsErr != nil {
		log.Warnf("problem while logging ems message, error: %v", emsErr)
	}
}

// Create a volume clone
func CreateOntapClone(name, source, snapshot, newSnapshotPrefix string, api *ontap.Driver) error {
	log.Debugf("OntapCommon#CreateOntapClone(%v, %v, %v, %v)", name, source, snapshot, newSnapshotPrefix)

	// If the specified volume already exists, skip creation and call it a success
	response, err := api.VolumeSize(name)
	if err != nil {
		return fmt.Errorf("Error searching for existing volume: error: %v", err)
	}
	if isPassed(response.Result.ResultStatusAttr) {
		return nil
	}

	// If no specific snapshot was requested, create one
	if snapshot == "" {
		// This is golang being stupid: https://golang.org/pkg/time/#Time.Format
		snapshot = newSnapshotPrefix + time.Now().UTC().Format("20060102T150405Z")
		response, err := api.SnapshotCreate(snapshot, source)
		if !isPassed(response.Result.ResultStatusAttr) || err != nil {
			return fmt.Errorf("Error creating snapshot: status: %v error: %v", response.Result.ResultStatusAttr, err)
		}
	}

	// Create the clone based on a snapshot
	response2, err2 := api.VolumeCloneCreate(name, source, snapshot)
	if !isPassed(response2.Result.ResultStatusAttr) || err2 != nil {
		if response2.Result.ResultErrnoAttr == azgo.EOBJECTNOTFOUND {
			return fmt.Errorf("Snapshot does not exist in volume")
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

func isPassed(s string) bool {
	const passed = "passed"
	return s == passed
}
