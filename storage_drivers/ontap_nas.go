// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

func init() {
	nas := &OntapNASStorageDriver{}
	nas.initialized = false
	Drivers[nas.Name()] = nas
	log.Debugf("Registered driver '%v'", nas.Name())
}

// OntapNASStorageDriverName is the constant name for this Ontap NAS storage driver
const OntapNASStorageDriverName = "ontap-nas"

// OntapNASStorageDriver is for NFS storage provisioning
type OntapNASStorageDriver struct {
	initialized bool
	config      OntapStorageDriverConfig
	api         *ontap.Driver
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
	decoder := json.NewDecoder(strings.NewReader(configJSON))
	err := decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	d.config = *config
	d.api, err = InitializeOntapDriver(d.config)
	if err != nil {
		return fmt.Errorf("Problem while initializing, error: %v", err)
	}

	validationErr := d.Validate()
	if validationErr != nil {
		return fmt.Errorf("Problem validating OntapNASStorageDriver error: %v", validationErr)
	}

	// log an informational message when this plugin starts
	EmsInitialized(d.Name(), d.api)

	d.initialized = true
	return nil
}

// Validate the driver configuration and execution environment
func (d *OntapNASStorageDriver) Validate() error {
	log.Debugf("OntapNASStorageDriver#Validate()")

	zr := &azgo.ZapiRunner{
		ManagementLIF: d.config.ManagementLIF,
		SVM:           d.config.SVM,
		Username:      d.config.Username,
		Password:      d.config.Password,
		Secure:        true,
	}

	r0, err0 := azgo.NewSystemGetVersionRequest().ExecuteUsing(zr)
	if err0 != nil {
		return fmt.Errorf("Could not validate credentials for %v@%v, error: %v", d.config.Username, d.config.SVM, err0)
	}

	// Add system version validation, if needed, this is a sanity check right now
	systemVersion := r0.Result
	if systemVersion.VersionPtr == nil {
		return fmt.Errorf("Could not determine system version for %v@%v", d.config.Username, d.config.SVM)
	}

	r1, err1 := d.api.NetInterfaceGet()
	if err1 != nil {
		return fmt.Errorf("Problem checking network interfaces error: %v", err1)
	}

	// if they didn't set a lif to use in the config, we'll set it to the first nfs lif we happen to find
	if d.config.DataLIF == "" {
	loop:
		for _, attrs := range r1.Result.AttributesList() {
			for _, protocol := range attrs.DataProtocols() {
				if protocol == "nfs" {
					log.Debugf("Setting NFS protocol access to '%v'", attrs.Address())
					d.config.DataLIF = string(attrs.Address())
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
				log.Debugf("Comparing NFS protocol access on : '%v' vs '%v'", attrs.Address(), d.config.DataLIF)
				if string(attrs.Address()) == d.config.DataLIF {
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

	response, _ := d.api.VolumeSize(name)
	if isPassed(response.Result.ResultStatusAttr) {
		log.Debugf("%v already exists, skipping volume create...", name)
		return nil
	}

	// get options with default values if not specified in config file
	// TODO add to documentation
	volumeSize := utils.GetV(opts, "size", "1g")
	spaceReserve := utils.GetV(opts, "spaceReserve", "none")
	snapshotPolicy := utils.GetV(opts, "snapshotPolicy", "none")
	unixPermissions := utils.GetV(opts, "unixPermissions", "---rwxr-xr-x")
	snapshotDir := utils.GetV(opts, "snapshotDir", "true")

	log.WithFields(log.Fields{
		"name":            name,
		"volumeSize":      volumeSize,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
	}).Debug("Creating volume with values")

	// create the volume
	response1, error1 := d.api.VolumeCreate(name, d.config.Aggregate, volumeSize, spaceReserve, snapshotPolicy, unixPermissions)
	if !isPassed(response1.Result.ResultStatusAttr) || error1 != nil {
		return fmt.Errorf("Error creating volume\n%verror: %v", response1.Result, error1)
	}

	// disable '.snapshot' to allow official mysql container's chmod-in-init to work
	if snapshotDir != "true" {
		response2, error2 := d.api.VolumeDisableSnapshotDirectoryAccess(name)
		if !isPassed(response2.Result.ResultStatusAttr) || error2 != nil {
			return fmt.Errorf("Error disabling snapshot directory access\n%verror: %v", response2.Result, error2)
		}
	}

	// mount the volume at the specified junction
	response3, error3 := d.api.VolumeMount(name, "/"+name)
	if !isPassed(response3.Result.ResultStatusAttr) || error3 != nil {
		return fmt.Errorf("Error mounting volume to junction\n%verror: %v", response3.Result, error3)
	}

	return nil
}

// Destroy the volume
func (d *OntapNASStorageDriver) Destroy(name string) error {
	log.Debugf("OntapNASStorageDriver#Destroy(%v)", name)

	response, error := d.api.VolumeDestroy(name, true)
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

	ip := d.config.DataLIF

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
