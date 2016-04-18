// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/netapp/netappdvp/storage_drivers"
	"github.com/netapp/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

var (
	debug      = flag.Bool("debug", true, "Enable debugging output")
	configFile = flag.String("config", "config.json", "Path to configuration file")
	driverID   = flag.String("volume-driver", "netapp", "Register as a docker volume plugin with this driver name")
	port       = flag.String("port", "", "Listen on this port instead of using a bsd socket")
)

func initLogging(logName string) *os.File {
	logRoot := "/var/log/netappdvp"

	// if logdir doesn't exist, make it
	dir, err := os.Lstat(logRoot)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logRoot, 0755); err != nil {
			fmt.Printf("problem creating log directory: '%v' error: %v", logRoot, err)
			os.Exit(1)
		}
	}
	// if logRoot isn't a directory, error
	if dir != nil && !dir.IsDir() {
		fmt.Printf("log directory '%v' exists and it's not a directory, please remove", logRoot)
		os.Exit(1)
	}

	// open a file for logging
	logFileLocation := ""
	switch runtime.GOOS {
	case utils.Linux:
		logFileLocation = logRoot + "/" + logName + ".log"
		break
	case utils.Darwin:
		logFileLocation = logRoot + "/" + logName + ".log"
		break
	}

	logFile, err := os.OpenFile(logFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening log file: %v error: %v", logFileLocation, err)
		os.Exit(1)
	}

	// assign it to the standard logger
	log.SetFormatter(&log.TextFormatter{}) // default for logrus
	log.SetOutput(logFile)                 // os.Stderr
	if *debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return logFile
}

func main() {

	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// open config file and read contents in to configJson
	fileContents, fileErr := ioutil.ReadFile(*configFile)
	if fileErr != nil {
		log.Error("Error reading configuration file: ", fileErr)
		os.Exit(1)
	}

	// validate the common settings and, if successful, we can continue and validate driver specific configuration
	configJSON := string(fileContents)
	commonConfig, commonConfigErr := storage_drivers.ValidateCommonSettings(configJSON)
	if commonConfigErr != nil {
		log.Errorf("Problem while validating configuration file: %v error: %v", *configFile, commonConfigErr)
		os.Exit(1)
	}

	// lookup the specified storageDriver
	storageDriver := storage_drivers.Drivers[commonConfig.StorageDriverName]
	if storageDriver == nil {
		log.Errorf("Unknown storage driver '%v' in configuration file: %v", commonConfig.StorageDriverName, *configFile)
		os.Exit(1)
	}

	// initialize the specified storageDriver which also triggers a call to Validate
	if initializeErr := storageDriver.Initialize(configJSON); initializeErr != nil {
		log.Errorf("Problem initializing storage driver: '%v' error: %v", commonConfig.StorageDriverName, initializeErr)
		os.Exit(1)
	}

	logFile := initLogging(*driverID)
	defer logFile.Close() // don't forget to close it
	log.Infof("Using storage driver: %v", commonConfig.StorageDriverName)
	log.Infof("Using config: %v", *commonConfig)

	volumeDir := filepath.Join(volume.DefaultDockerRootDirectory, *driverID)
	log.WithFields(log.Fields{
		"volumeDir":     volumeDir,
		"volume-driver": *driverID,
		"port":          *port,
	}).Info("Starting docker volume plugin with the following options:")

	// plugin connection registered in /var/run/docker/plugins
	d, err := newNetAppDockerVolumePlugin(volumeDir, *commonConfig)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	h := volume.NewHandler(d)
	if *port != "" {
		log.Info(h.ServeTCP(*driverID, ":"+*port))
	} else {
		log.Info(h.ServeUnix("root", *driverID)) // 'root' here is the unix group to start as
	}
}
