// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"
)

// OntapNASStorageDriverName is the constant name for this Ontap NAS storage driver
const OntapNASStorageDriverName = "ontap-nas"

func init() {
	nas := &OntapNASStorageDriver{}
	nas.Initialized = false
	Drivers[nas.Name()] = nas
	log.Debugf("Registered driver '%v'", nas.Name())
}

// OntapNASStorageDriver is for NFS storage provisioning
type OntapNASStorageDriver struct {
	Initialized bool
	Config      OntapStorageDriverConfig
	API         *ontap.Driver
}

func (d *OntapNASStorageDriver) GetConfig() *OntapStorageDriverConfig {
	return &d.Config
}

func (d *OntapNASStorageDriver) GetAPI() *ontap.Driver {
	return d.API
}

// Name is for returning the name of this driver
func (d *OntapNASStorageDriver) Name() string {
	log.Debugf("OntapNASStorageDriver#Name()")
	return OntapNASStorageDriverName
}

// Initialize from the provided config
func (d *OntapNASStorageDriver) Initialize(configJSON string) error {
	log.Debugf("OntapNASStorageDriver#Initialize(...)")

	config := &OntapStorageDriverConfig{}

	// decode configJSON into OntapStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	log.WithFields(log.Fields{
		"Version":           config.Version,
		"StorageDriverName": config.StorageDriverName,
		"DisableDelete":     config.DisableDelete,
		"StoragePrefixRaw":  string(config.StoragePrefixRaw),
		"SnapshotPrefixRaw": string(config.SnapshotPrefixRaw),
	}).Debugf("Reparsed into ontapConfig")

	d.Config = *config
	d.API, err = InitializeOntapDriver(d.Config)
	if err != nil {
		return fmt.Errorf("Problem while initializing: %v", err)
	}

	defaultsErr := PopulateConfigurationDefaults(&d.Config)
	if defaultsErr != nil {
		return fmt.Errorf("Cannot populate configuration defaults: %v", defaultsErr)
	}

	validationErr := d.Validate()
	if validationErr != nil {
		return fmt.Errorf("Problem validating OntapNASStorageDriver: %v", validationErr)
	}

	// log an informational message on a heartbeat
	EmsInitialized(d.Name(), d.API, config)
	StartEmsHeartbeat(d.Name(), d.API, config)

	d.Initialized = true
	log.WithFields(log.Fields{
		"driverVersion":         DriverVersion,
		"extendedDriverVersion": ExtendedDriverVersion,
	}).Info("Initialized Ontap NAS storage driver.")
	return nil
}

// Validate the driver configuration and execution environment
func (d *OntapNASStorageDriver) Validate() error {
	log.Debugf("OntapNASStorageDriver#Validate()")

	zr := &azgo.ZapiRunner{
		ManagementLIF: d.Config.ManagementLIF,
		SVM:           d.Config.SVM,
		Username:      d.Config.Username,
		Password:      d.Config.Password,
		Secure:        true,
	}

	r0, err0 := azgo.NewSystemGetVersionRequest().ExecuteUsing(zr)
	if err0 != nil {
		return fmt.Errorf("Could not validate credentials for %v@%v, error: %v", d.Config.Username, d.Config.SVM, err0)
	}

	// Add system version validation, if needed, this is a sanity check right now
	systemVersion := r0.Result
	if systemVersion.VersionPtr == nil {
		return fmt.Errorf("Could not determine system version for %v@%v", d.Config.Username, d.Config.SVM)
	}

	r1, err1 := d.API.NetInterfaceGet()
	if err1 != nil {
		return fmt.Errorf("Problem checking network interfaces error: %v", err1)
	}

	// if they didn't set a lif to use in the config, we'll set it to the first nfs lif we happen to find
	if d.Config.DataLIF == "" {
	loop:
		for _, attrs := range r1.Result.AttributesList() {
			for _, protocol := range attrs.DataProtocols() {
				if protocol == "nfs" {
					log.Debugf("Setting NFS protocol access to '%v'", attrs.Address())
					d.Config.DataLIF = string(attrs.Address())
					break loop
				}
			}
		}
	}

	foundNfs := false
loop2:
	for _, attrs := range r1.Result.AttributesList() {
		for _, protocol := range attrs.DataProtocols() {
			if protocol == "nfs" {
				log.Debugf("Comparing NFS protocol access on : '%v' vs '%v'", attrs.Address(), d.Config.DataLIF)
				if string(attrs.Address()) == d.Config.DataLIF {
					foundNfs = true
					break loop2
				}
			}
		}
	}

	if !foundNfs {
		return fmt.Errorf("Could not find NFS DataLIF")
	}

	return nil
}

// Create a volume with the specified options
func (d *OntapNASStorageDriver) Create(name string, opts map[string]string) error {
	log.Debugf("OntapNASStorageDriver#Create(%v)", name)

	// If the volume already exists, bail out
	response, _ := d.API.VolumeSize(name)
	if isPassed(response.Result.ResultStatusAttr) {
		return fmt.Errorf("Volume already exists")
	}

	// get options with default fallback values
	// see also: ontap_common.go#PopulateConfigurationDefaults
	volumeSize := utils.GetV(opts, "size", d.Config.VolumeSize)
	spaceReserve := utils.GetV(opts, "spaceReserve", d.Config.SpaceReserve)
	snapshotPolicy := utils.GetV(opts, "snapshotPolicy", d.Config.SnapshotPolicy)
	unixPermissions := utils.GetV(opts, "unixPermissions", d.Config.UnixPermissions)
	snapshotDir := utils.GetV(opts, "snapshotDir", d.Config.SnapshotDir)
	exportPolicy := utils.GetV(opts, "exportPolicy", d.Config.ExportPolicy)
	aggregate := utils.GetV(opts, "aggregate", d.Config.Aggregate)
	securityStyle := utils.GetV(opts, "securityStyle", d.Config.SecurityStyle)

	log.WithFields(log.Fields{
		"name":            name,
		"volumeSize":      volumeSize,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
		"snapshotDir":     snapshotDir,
		"exportPolicy":    exportPolicy,
		"aggregate":       aggregate,
		"securityStyle":   securityStyle,
	}).Debug("Creating volume with values")

	// create the volume
	response1, error1 := d.API.VolumeCreate(name, aggregate, volumeSize, spaceReserve, snapshotPolicy, unixPermissions, exportPolicy, securityStyle)
	if !isPassed(response1.Result.ResultStatusAttr) || error1 != nil {
		if response1.Result.ResultErrnoAttr != azgo.EAPIERROR {
			return fmt.Errorf("Error creating volume\n%verror: %v", response1.Result, error1)
		} else {
			if !strings.HasSuffix(strings.TrimSpace(response1.Result.ResultReasonAttr), "Job exists") {
				return fmt.Errorf("Error creating volume\n%verror: %v", response1.Result, error1)
			} else {
				log.Warnf("%v volume create job already exists, skipping volume create on this node...", name)
				return nil
			}
		}
	}

	// disable '.snapshot' to allow official mysql container's chmod-in-init to work
	if snapshotDir != "true" {
		response2, error2 := d.API.VolumeDisableSnapshotDirectoryAccess(name)
		if !isPassed(response2.Result.ResultStatusAttr) || error2 != nil {
			return fmt.Errorf("Error disabling snapshot directory access\n%verror: %v", response2.Result, error2)
		}
	}

	// mount the volume at the specified junction
	response3, error3 := d.API.VolumeMount(name, "/"+name)
	if !isPassed(response3.Result.ResultStatusAttr) || error3 != nil {
		return fmt.Errorf("Error mounting volume to junction\n%verror: %v", response3.Result, error3)
	}

	return nil
}

// Create a volume clone
func (d *OntapNASStorageDriver) CreateClone(name, source, snapshot, newSnapshotPrefix string) error {
	return CreateOntapClone(name, source, snapshot, newSnapshotPrefix, d.API)
}

// Destroy the volume
func (d *OntapNASStorageDriver) Destroy(name string) error {
	log.Debugf("OntapNASStorageDriver#Destroy(%v)", name)

	// TODO: If this is the parent of one or more clones, those clones have to split from this
	// volume before it can be deleted, which means separate copies of those volumes.
	// If there are a lot of clones on this volume, that could seriously balloon the amount of
	// utilized space. Is that what we want? Or should we just deny the delete, and force the
	// user to keep the volume around until all of the clones are gone? If we do that, need a
	// way to list the clones. Maybe volume inspect.

	response, error := d.API.VolumeDestroy(name, true)
	if !isPassed(response.Result.ResultStatusAttr) || error != nil {
		if response.Result.ResultErrnoAttr != azgo.EVOLUMEDOESNOTEXIST {
			return fmt.Errorf("Error destroying volume: %v\n%verror: %v", name, response.Result, error)
		} else {
			log.Warnf("Volume already deleted while destroying volume: %v\n%verror: %v", name, response.Result, error)
		}

	}
	return nil
}

// Attach the volume
func (d *OntapNASStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("OntapNASStorageDriver#Attach(%v, %v, %v)", name, mountpoint, opts)

	ip := d.Config.DataLIF

	var cmd string
	switch runtime.GOOS {
	case utils.Linux:
		cmd = fmt.Sprintf("mount -o nfsvers=3 %s:/%s %s", ip, name, mountpoint)
	case utils.Darwin:
		cmd = fmt.Sprintf("mount -o rw -t nfs %s:/%s %s", ip, name, mountpoint)
	default:
		return fmt.Errorf("Unsupported operating system: %v", runtime.GOOS)
	}
	log.Debugf("mount cmd==%s", cmd)

	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.Debugf("out==%v", string(out))
		return fmt.Errorf("Problem mounting volume: %v mountpoint: %v error: %v", name, mountpoint, err)
	}

	return nil
}

// Detach the volume
func (d *OntapNASStorageDriver) Detach(name, mountpoint string) error {
	log.Debugf("OntapNASStorageDriver#Detach(%v, %v)", name, mountpoint)

	cmd := fmt.Sprintf("umount %s", mountpoint)
	log.Debugf("cmd==%s", cmd)
	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.Debugf("out==%v", string(out))
		return fmt.Errorf("Problem unmounting docker volume: %v mountpoint: %v error: %v", name, mountpoint, err)
	}

	return nil
}

// DefaultStoragePrefix is the driver specific prefix for created storage, can be overridden in the config file
func (d *OntapNASStorageDriver) DefaultStoragePrefix() string {
	return "netappdvp_"
}

// DefaultSnapshotPrefix is the driver specific prefix for created snapshots, can be overridden in the config file
func (d *OntapNASStorageDriver) DefaultSnapshotPrefix() string {
	return "netappdvp_"
}

// Return the list of snapshots associated with the named volume
func (d *OntapNASStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {
	return GetSnapshotList(name, d.API)
}

// Return the list of volumes associated with this tenant
func (d *OntapNASStorageDriver) List(prefix string) ([]string, error) {
	return GetVolumeList(prefix, d.API)
}

// Test for the existence of a volume
func (d *OntapNASStorageDriver) Get(name string) error {
	return GetVolume(name, d.API)
}
