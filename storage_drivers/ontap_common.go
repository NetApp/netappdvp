// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/netapp/netappdvp/apis/ontap"

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
		return nil, fmt.Errorf("Error ennumerating SVMs:  status: %v error: %v", response1.Result.ResultStatusAttr, err)
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

func isPassed(s string) bool {
	const passed = "passed"
	return s == passed
}
