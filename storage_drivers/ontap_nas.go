// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"fmt"
	"strconv"
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
	return OntapNASStorageDriverName
}

// Initialize from the provided config
func (d *OntapNASStorageDriver) Initialize(
	context DriverContext, configJSON string, commonConfig *CommonStorageDriverConfig,
) error {

	if commonConfig.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Initialize", "Type": "OntapNASStorageDriver"}
		log.WithFields(fields).Debug(">>>> Initialize")
		defer log.WithFields(fields).Debug("<<<< Initialize")
	}

	// Parse the config
	config, err := InitializeOntapConfig(configJSON, commonConfig)
	if err != nil {
		return fmt.Errorf("Error initializing %s driver. %v", d.Name(), err)
	}

	d.Config = *config
	d.API, err = InitializeOntapDriver(&d.Config)
	if err != nil {
		return fmt.Errorf("Error initializing %s driver. %v", d.Name(), err)
	}

	err = d.Validate(context)
	if err != nil {
		return fmt.Errorf("Error validating %s driver. %v", d.Name(), err)
	}

	// log an informational message on a heartbeat
	EmsInitialized(d.Name(), d.API, &d.Config)
	StartEmsHeartbeat(d.Name(), d.API, &d.Config)

	d.Initialized = true
	return nil
}

// Validate the driver configuration and execution environment
func (d *OntapNASStorageDriver) Validate(context DriverContext) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Validate", "Type": "OntapNASStorageDriver", "context": context}
		log.WithFields(fields).Debug(">>>> Validate")
		defer log.WithFields(fields).Debug("<<<< Validate")
	}

	err := ValidateNASDriver(context, d.API, &d.Config)
	if err != nil {
		return fmt.Errorf("Driver validation failed. %v", err)
	}

	return nil
}

// Create a volume with the specified options
func (d *OntapNASStorageDriver) Create(name string, sizeBytes uint64, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":    "Create",
			"Type":      "OntapNASStorageDriver",
			"name":      name,
			"sizeBytes": sizeBytes,
			"opts":      opts,
		}
		log.WithFields(fields).Debug(">>>> Create")
		defer log.WithFields(fields).Debug("<<<< Create")
	}

	// If the volume already exists, bail out
	volExists, err := d.API.VolumeExists(name)
	if err != nil {
		return fmt.Errorf("Error checking for existing volume. %v", err)
	}
	if volExists {
		return fmt.Errorf("Volume %s already exists.", name)
	}

	if sizeBytes < OntapMinimumVolumeSizeBytes {
		return fmt.Errorf("Requested volume size (%d bytes) is too small.  The minimum volume size is %d bytes.",
			sizeBytes, OntapMinimumVolumeSizeBytes)
	}

	// get options with default fallback values
	// see also: ontap_common.go#PopulateConfigurationDefaults
	size := strconv.FormatUint(sizeBytes, 10)
	spaceReserve := utils.GetV(opts, "spaceReserve", d.Config.SpaceReserve)
	snapshotPolicy := utils.GetV(opts, "snapshotPolicy", d.Config.SnapshotPolicy)
	unixPermissions := utils.GetV(opts, "unixPermissions", d.Config.UnixPermissions)
	snapshotDir := utils.GetV(opts, "snapshotDir", d.Config.SnapshotDir)
	exportPolicy := utils.GetV(opts, "exportPolicy", d.Config.ExportPolicy)
	aggregate := utils.GetV(opts, "aggregate", d.Config.Aggregate)
	securityStyle := utils.GetV(opts, "securityStyle", d.Config.SecurityStyle)
	encryption := utils.GetV(opts, "encryption", d.Config.Encryption)

	enableSnapshotDir, err := strconv.ParseBool(snapshotDir)
	if err != nil {
		return fmt.Errorf("Invalid boolean value for snapshotDir. %v", err)
	}

	encrypt, err := ValidateEncryptionAttribute(encryption, d.API)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"name":            name,
		"size":            size,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
		"snapshotDir":     enableSnapshotDir,
		"exportPolicy":    exportPolicy,
		"aggregate":       aggregate,
		"securityStyle":   securityStyle,
		"encryption":      encryption,
	}).Debug("Creating Flexvol.")

	// Create the volume
	volCreateResponse, err := d.API.VolumeCreate(
		name, aggregate, size, spaceReserve, snapshotPolicy,
		unixPermissions, exportPolicy, securityStyle, encrypt)

	if err = ontap.GetError(volCreateResponse, err); err != nil {
		if zerr, ok := err.(ontap.ZapiError); ok {
			// Handle case where the Create is passed to every Docker Swarm node
			if zerr.Code() == azgo.EAPIERROR && strings.HasSuffix(strings.TrimSpace(zerr.Reason()), "Job exists") {
				log.WithField("volume", name).Warn("Volume create job already exists, skipping volume create on this node.")
				return nil
			}
		}
		return fmt.Errorf("Error creating volume. %v", err)
	}

	// Disable '.snapshot' to allow official mysql container's chmod-in-init to work
	if !enableSnapshotDir {
		snapDirResponse, err := d.API.VolumeDisableSnapshotDirectoryAccess(name)
		if err = ontap.GetError(snapDirResponse, err); err != nil {
			return fmt.Errorf("Error disabling snapshot directory access. %v", err)
		}
	}

	// Mount the volume at the specified junction
	mountResponse, err := d.API.VolumeMount(name, "/"+name)
	if err = ontap.GetError(mountResponse, err); err != nil {
		return fmt.Errorf("Error mounting volume to junction. %v", err)
	}

	// If LS mirrors are present on the SVM root volume, update them
	UpdateLoadSharingMirrors(d.API)

	return nil
}

// Create a volume clone
func (d *OntapNASStorageDriver) CreateClone(name, source, snapshot string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":   "CreateClone",
			"Type":     "OntapNASStorageDriver",
			"name":     name,
			"source":   source,
			"snapshot": snapshot,
			"opts":     opts,
		}
		log.WithFields(fields).Debug(">>>> CreateClone")
		defer log.WithFields(fields).Debug("<<<< CreateClone")
	}

	split, err := strconv.ParseBool(utils.GetV(opts, "splitOnClone", d.Config.SplitOnClone))
	if err != nil {
		return fmt.Errorf("Invalid boolean value for splitOnClone. %v", err)
	}

	log.WithField("splitOnClone", split).Debug("Creating volume clone.")
	return CreateOntapClone(name, source, snapshot, split, &d.Config, d.API)
}

// Destroy the volume
func (d *OntapNASStorageDriver) Destroy(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Destroy",
			"Type":   "OntapNASStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> Destroy")
		defer log.WithFields(fields).Debug("<<<< Destroy")
	}

	// TODO: If this is the parent of one or more clones, those clones have to split from this
	// volume before it can be deleted, which means separate copies of those volumes.
	// If there are a lot of clones on this volume, that could seriously balloon the amount of
	// utilized space. Is that what we want? Or should we just deny the delete, and force the
	// user to keep the volume around until all of the clones are gone? If we do that, need a
	// way to list the clones. Maybe volume inspect.

	volDestroyResponse, err := d.API.VolumeDestroy(name, true)
	if err != nil {
		return fmt.Errorf("Error destroying volume %v. %v", name, err)
	}
	if zerr := ontap.NewZapiError(volDestroyResponse); !zerr.IsPassed() {

		// It's not an error if the volume no longer exists
		if zerr.Code() == azgo.EVOLUMEDOESNOTEXIST {
			log.WithField("volume", name).Warn("Volume already deleted.")
		} else {
			return fmt.Errorf("Error destroying volume %v. %v", name, zerr)
		}
	}

	return nil
}

// Attach the volume
func (d *OntapNASStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Attach",
			"Type":       "OntapNASStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
			"opts":       opts,
		}
		log.WithFields(fields).Debug(">>>> Attach")
		defer log.WithFields(fields).Debug("<<<< Attach")
	}

	exportPath := fmt.Sprintf("%s:/%s", d.Config.DataLIF, name)

	return MountVolume(exportPath, mountpoint, &d.Config)
}

// Detach the volume
func (d *OntapNASStorageDriver) Detach(name, mountpoint string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Detach",
			"Type":       "OntapNASStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
		}
		log.WithFields(fields).Debug(">>>> Detach")
		defer log.WithFields(fields).Debug("<<<< Detach")
	}

	return UnmountVolume(mountpoint, &d.Config)
}

// Return the list of snapshots associated with the named volume
func (d *OntapNASStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "SnapshotList",
			"Type":   "OntapNASStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> SnapshotList")
		defer log.WithFields(fields).Debug("<<<< SnapshotList")
	}

	return GetSnapshotList(name, &d.Config, d.API)
}

// Return the list of volumes associated with this tenant
func (d *OntapNASStorageDriver) List() ([]string, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "List", "Type": "OntapNASStorageDriver"}
		log.WithFields(fields).Debug(">>>> List")
		defer log.WithFields(fields).Debug("<<<< List")
	}

	return GetVolumeList(d.API, &d.Config)
}

// Test for the existence of a volume
func (d *OntapNASStorageDriver) Get(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Get", "Type": "OntapNASStorageDriver"}
		log.WithFields(fields).Debug(">>>> Get")
		defer log.WithFields(fields).Debug("<<<< Get")
	}

	return GetVolume(name, d.API, &d.Config)
}
