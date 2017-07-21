// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/netapp/netappdvp/storage_drivers"
)

var (
	debug        = flag.Bool("debug", false, "Enable debugging output")
	logLevel     = flag.String("log-level", "info", "Logging level (debug, info, warn, error, fatal)")
	configFile   = flag.String("config", "config.json", "Path to configuration file")
	driverID     = flag.String("volume-driver", "netapp", "Register as a docker volume plugin with this driver name")
	port         = flag.String("port", "", "Listen on this port instead of using a bsd socket")
	printVersion = flag.Bool("version", false, "Print version and exit")
)

func main() {

	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(int64(time.Now().Nanosecond()))

	if *printVersion == true {
		fmt.Printf("NetApp Docker Volume Plugin\nVersion: %v\n", storage_drivers.FullDriverVersion)
		fmt.Printf("Build: %v\n", storage_drivers.BuildVersion)
		fmt.Printf("Built: %v\n", storage_drivers.BuildTime)
		os.Exit(0)
	}

	// Set up logging
	err := initLogging()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	// Read config file
	fileContents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"configFile": *configFile,
			"error":      err,
		}).Fatal("Problem reading configuration file.")
	}

	// Validate common config settings
	configJSON := string(fileContents)
	commonConfig, err := storage_drivers.ValidateCommonSettings(configJSON)
	if err != nil {
		log.WithFields(log.Fields{
			"configFile": *configFile,
			"error":      err,
		}).Fatal("Problem validating configuration file.")
	}

	log.Debugf("Parsed commonConfig: %+v", *commonConfig)

	// Lookup the specified storage driver
	storageDriver := storage_drivers.Drivers[commonConfig.StorageDriverName]
	if storageDriver == nil {
		log.WithFields(log.Fields{
			"configFile":        *configFile,
			"storageDriverName": commonConfig.StorageDriverName,
		}).Fatal("Unknown storage driver in configuration file.")
	}

	// Determine path to Docker volumes
	volumePath := filepath.Join(volume.DefaultDockerRootDirectory, *driverID)

	log.WithFields(log.Fields{
		"version":       storage_drivers.DriverVersion,
		"mode":          storage_drivers.ExtendedDriverVersion,
		"storageDriver": commonConfig.StorageDriverName,
		"volumeDriver":  *driverID,
		"port":          *port,
	}).Info("Initializing storage driver")

	// Initialize the specified storage driver which also triggers a call to Validate
	if err := storageDriver.Initialize(storage_drivers.ContextNDVP, configJSON, commonConfig); err != nil {
		log.WithFields(log.Fields{
			"storageDriverName": commonConfig.StorageDriverName,
			"error":             err,
		}).Fatal("Problem initializing storage driver.")
	}

	// Plugin connection registered in /var/run/docker/plugins
	d, err := newNetAppDockerVolumePlugin(volumePath, *commonConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Initalized driver; plugin ready!")

	h := volume.NewHandler(d)

	if *port != "" {
		log.Error(h.ServeTCP(*driverID, ":"+*port, nil))
	} else {
		log.Error(h.ServeUnix(*driverID, 0)) // 0 is the unix group to start as (root gid)
	}

	log.Fatal("Unexpected exit")
}
