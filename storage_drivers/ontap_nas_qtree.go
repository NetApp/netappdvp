// Copyright 2017 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"
)

// OntapNASQtreeStorageDriverName is the constant name for this Ontap qtree-based NAS storage driver
const OntapNASQtreeStorageDriverName = "ontap-nas-economy"
const deletedQtreeNamePrefix = "deleted_"
const maxQtreeNameLength = 64
const maxQtreesPerFlexvol = 200
const defaultPruneFlexvolsPeriodSecs = uint64(600) // default to 10 minutes
const defaultResizeQuotasPeriodSecs = uint64(60)   // default to 1 minute

func init() {
	nas := &OntapNASQtreeStorageDriver{}
	nas.Initialized = false
	Drivers[nas.Name()] = nas
}

// OntapNASQtreeStorageDriver is for NFS storage provisioning of qtrees
type OntapNASQtreeStorageDriver struct {
	Initialized         bool
	Config              OntapStorageDriverConfig
	API                 *ontap.Driver
	quotaResizeMap      map[string]bool
	provMutex           *sync.Mutex
	flexvolNamePrefix   string
	flexvolExportPolicy string
}

func (d *OntapNASQtreeStorageDriver) GetConfig() *OntapStorageDriverConfig {
	return &d.Config
}

func (d *OntapNASQtreeStorageDriver) GetAPI() *ontap.Driver {
	return d.API
}

// Name is for returning the name of this driver
func (d *OntapNASQtreeStorageDriver) Name() string {
	return OntapNASQtreeStorageDriverName
}

func (d *OntapNASQtreeStorageDriver) FlexvolNamePrefix() string {
	return d.flexvolNamePrefix
}

// Initialize from the provided config
func (d *OntapNASQtreeStorageDriver) Initialize(
	context DriverContext, configJSON string, commonConfig *CommonStorageDriverConfig,
) error {

	if commonConfig.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Initialize", "Type": "OntapNASQtreeStorageDriver"}
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

	// Set up internal driver state
	d.quotaResizeMap = make(map[string]bool)
	d.provMutex = &sync.Mutex{}
	d.flexvolNamePrefix = fmt.Sprintf("%s_qtree_pool_%s_", string(context), *d.Config.StoragePrefix)
	d.flexvolNamePrefix = strings.Replace(d.flexvolNamePrefix, "__", "_", -1)
	d.flexvolExportPolicy = fmt.Sprintf("%s_qtree_pool_export_policy", string(context))

	log.WithFields(log.Fields{
		"FlexvolNamePrefix":   d.flexvolNamePrefix,
		"FlexvolExportPolicy": d.flexvolExportPolicy,
	}).Debugf("Qtree driver settings.")

	err = d.Validate(context)
	if err != nil {
		return fmt.Errorf("Error validating %s driver. %v", d.Name(), err)
	}

	// Log an informational message on a heartbeat
	EmsInitialized(d.Name(), d.API, &d.Config)

	// Ensure all quotas are in force after a driver restart
	d.queueAllFlexvolsForQuotaResize()

	// Do periodic housekeeping like cleaning up unused Flexvols
	d.startHousekeepingTasks()

	d.Initialized = true
	return nil
}

// Validate the driver configuration and execution environment
func (d *OntapNASQtreeStorageDriver) Validate(context DriverContext) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Validate", "Type": "OntapNASQtreeStorageDriver", "context": context}
		log.WithFields(fields).Debug(">>>> Validate")
		defer log.WithFields(fields).Debug("<<<< Validate")
	}

	err := ValidateNASDriver(context, d.API, &d.Config)
	if err != nil {
		return fmt.Errorf("Driver validation failed. %v", err)
	}

	// Make sure we have an export policy for all the Flexvols we create
	err = d.ensureDefaultExportPolicy()
	if err != nil {
		return fmt.Errorf("Error configuring export policy. %v", err)
	}

	return nil
}

func (d *OntapNASQtreeStorageDriver) startHousekeepingTasks() {

	// Send EMS message on a configurable schedule
	StartEmsHeartbeat(d.Name(), d.API, &d.Config)

	// Read background task timings from config file, use defaults if missing or invalid
	pruneFlexvolsPeriodSecs := defaultPruneFlexvolsPeriodSecs
	if d.Config.QtreePruneFlexvolsPeriod != "" {
		i, err := strconv.ParseUint(d.Config.QtreePruneFlexvolsPeriod, 10, 64)
		if err != nil {
			log.WithField("interval", d.Config.QtreePruneFlexvolsPeriod).Warnf(
				"Invalid Flexvol pruning interval. %v", err)
		} else {
			pruneFlexvolsPeriodSecs = i
		}
	}
	log.WithFields(log.Fields{
		"IntervalSeconds": pruneFlexvolsPeriodSecs,
	}).Debug("Configured Flexvol pruning period.")

	resizeQuotasPeriodSecs := defaultResizeQuotasPeriodSecs
	if d.Config.QtreeQuotaResizePeriod != "" {
		i, err := strconv.ParseUint(d.Config.QtreeQuotaResizePeriod, 10, 64)
		if err != nil {
			log.WithField("interval", d.Config.QtreeQuotaResizePeriod).Warnf(
				"Invalid quota resize interval. %v", err)
		} else {
			resizeQuotasPeriodSecs = i
		}
	}
	log.WithFields(log.Fields{
		"IntervalSeconds": resizeQuotasPeriodSecs,
	}).Debug("Configured quota resize period.")

	// Keep the system devoid of Flexvols with no qtrees
	d.pruneUnusedFlexvols()
	d.reapDeletedQtrees()
	pruneTicker := time.NewTicker(time.Duration(pruneFlexvolsPeriodSecs) * time.Second)
	go func() {
		for range pruneTicker.C {
			d.pruneUnusedFlexvols()
			d.reapDeletedQtrees()
		}
	}()

	// Keep the quotas current
	d.resizeQuotas()
	resizeTicker := time.NewTicker(time.Duration(resizeQuotasPeriodSecs) * time.Second)
	go func() {
		for range resizeTicker.C {
			d.resizeQuotas()
		}
	}()
}

// Create a qtree-backed volume with the specified options
func (d *OntapNASQtreeStorageDriver) Create(name string, sizeBytes uint64, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":    "Create",
			"Type":      "OntapNASQtreeStorageDriver",
			"name":      name,
			"sizeBytes": sizeBytes,
			"opts":      opts,
		}
		log.WithFields(fields).Debug(">>>> Create")
		defer log.WithFields(fields).Debug("<<<< Create")
	}

	// Generic user-facing message
	createError := errors.New("Volume creation failed.")

	// Ensure volume doesn't already exist
	exists, existsInFlexvol, err := d.API.QtreeExists(name, d.FlexvolNamePrefix())
	if err != nil {
		log.Errorf("Error checking for existing volume. %v", err)
		return createError
	}
	if exists {
		log.WithFields(log.Fields{"qtree": name, "flexvol": existsInFlexvol}).Debug("Qtree already exists.")
		return fmt.Errorf("Volume %s already exists.", name)
	}

	if sizeBytes < OntapMinimumVolumeSizeBytes {
		return fmt.Errorf("Requested volume size (%d bytes) is too small.  The minimum volume size is %d bytes.",
			sizeBytes, OntapMinimumVolumeSizeBytes)
	}

	// Ensure qtree name isn't too long
	if len(name) > maxQtreeNameLength {
		return fmt.Errorf("Volume %s name exceeds the limit of %d characters.", name, maxQtreeNameLength)
	}

	// Get Flexvol options with default fallback values
	// see also: ontap_common.go#PopulateConfigurationDefaults
	size := strconv.FormatUint(sizeBytes, 10)
	aggregate := utils.GetV(opts, "aggregate", d.Config.Aggregate)
	spaceReserve := utils.GetV(opts, "spaceReserve", d.Config.SpaceReserve)
	snapshotPolicy := utils.GetV(opts, "snapshotPolicy", d.Config.SnapshotPolicy)
	snapshotDir := utils.GetV(opts, "snapshotDir", d.Config.SnapshotDir)
	encryption := utils.GetV(opts, "encryption", d.Config.Encryption)

	enableSnapshotDir, err := strconv.ParseBool(snapshotDir)
	if err != nil {
		return fmt.Errorf("Invalid boolean value for snapshotDir. %v", err)
	}

	encrypt, err := ValidateEncryptionAttribute(encryption, d.API)
	if err != nil {
		return err
	}

	// Ensure any Flexvol we use won't be pruned before we place a qtree on it
	d.provMutex.Lock()
	defer d.provMutex.Unlock()

	// Make sure we have a Flexvol for the new qtree
	flexvol, err := d.ensureFlexvolForQtree(
		aggregate, spaceReserve, snapshotPolicy, enableSnapshotDir, encrypt)
	if err != nil {
		log.Errorf("Flexvol location/creation failed. %v", err)
		return createError
	}

	// Grow or shrink the Flexvol as needed
	flexvolSizeBytes, err := d.getOptimalSizeForFlexvol(flexvol, sizeBytes)
	if err != nil {
		log.Warnf("Could not calculate optimal Flexvol size. %v", err)

		// Lacking the optimal size, just grow the Flexvol to contain the new qtree
		resizeResponse, err := d.API.SetVolumeSize(flexvol, "+"+size)
		if err = ontap.GetError(resizeResponse.Result, err); err != nil {
			log.Errorf("Flexvol resize failed. %v", err)
			return createError
		}
	} else {

		// Got optimal size, so just set the Flexvol to that value
		flexvolSizeStr := strconv.FormatUint(flexvolSizeBytes, 10)
		resizeResponse, err := d.API.SetVolumeSize(flexvol, flexvolSizeStr)
		if err = ontap.GetError(resizeResponse.Result, err); err != nil {
			log.Errorf("Flexvol resize failed. %v", err)
			return createError
		}
	}

	// Get qtree options with default fallback values
	unixPermissions := utils.GetV(opts, "unixPermissions", d.Config.UnixPermissions)
	exportPolicy := utils.GetV(opts, "exportPolicy", d.Config.ExportPolicy)
	securityStyle := utils.GetV(opts, "securityStyle", d.Config.SecurityStyle)

	// Create the qtree
	qtreeResponse, err := d.API.QtreeCreate(name, flexvol, unixPermissions, exportPolicy, securityStyle)
	if err = ontap.GetError(qtreeResponse, err); err != nil {
		log.Errorf("Qtree creation failed. %v", err)
		return createError
	}

	// Add the quota
	d.addQuotaForQtree(name, flexvol, sizeBytes)
	if err != nil {
		log.Errorf("Qtree quota definition failed. %v", err)
		return createError
	}

	return nil
}

// Create a volume clone
func (d *OntapNASQtreeStorageDriver) CreateClone(name, source, snapshot string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":   "CreateClone",
			"Type":     "OntapNASQtreeStorageDriver",
			"name":     name,
			"source":   source,
			"snapshot": snapshot,
			"opts":     opts,
		}
		log.WithFields(fields).Debug(">>>> CreateClone")
		defer log.WithFields(fields).Debug("<<<< CreateClone")
	}

	return errors.New("Cloning with the ONTAP NAS Qtree driver is not supported.")
}

// Destroy the volume
func (d *OntapNASQtreeStorageDriver) Destroy(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Destroy",
			"Type":   "OntapNASQtreeStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> Destroy")
		defer log.WithFields(fields).Debug("<<<< Destroy")
	}

	// Generic user-facing message
	deleteError := errors.New("Volume deletion failed.")

	exists, flexvol, err := d.API.QtreeExists(name, d.FlexvolNamePrefix())
	if err != nil {
		log.Errorf("Error checking for existing qtree. %v", err)
		return deleteError
	}
	if !exists {
		log.WithField("qtree", name).Warn("Qtree not found.")
		return nil
	}

	// Rename qtree so it doesn't show up in lists while ONTAP is deleting it in the background.
	// Ensure the deleted name doesn't exceed the qtree name length limit of 64 characters.
	path := fmt.Sprintf("/vol/%s/%s", flexvol, name)
	deletedName := deletedQtreeNamePrefix + name + "_" + utils.RandomString(5)
	if len(deletedName) > maxQtreeNameLength {
		trimLength := len(deletedQtreeNamePrefix) + 10
		deletedName = deletedQtreeNamePrefix + name[trimLength:] + "_" + utils.RandomString(5)
	}
	deletedPath := fmt.Sprintf("/vol/%s/%s", flexvol, deletedName)

	renameResponse, err := d.API.QtreeRename(path, deletedPath)
	if err = ontap.GetError(renameResponse, err); err != nil {
		log.Errorf("Qtree rename failed. %v", err)
		return deleteError
	}

	// Destroy the qtree in the background.  If this fails, try to restore the original qtree name.
	destroyResponse, err := d.API.QtreeDestroyAsync(deletedPath, true)
	if err = ontap.GetError(destroyResponse, err); err != nil {
		log.Errorf("Qtree async delete failed. %v", err)
		defer d.API.QtreeRename(deletedPath, path)
		return deleteError
	}

	return nil
}

// Attach the volume
func (d *OntapNASQtreeStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Attach",
			"Type":       "OntapNASQtreeStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
			"opts":       opts,
		}
		log.WithFields(fields).Debug(">>>> Attach")
		defer log.WithFields(fields).Debug("<<<< Attach")
	}

	// Check if qtree exists, and find its Flexvol so we can build the export location
	exists, flexvol, err := d.API.QtreeExists(name, d.FlexvolNamePrefix())
	if err != nil {
		log.Errorf("Error checking for existing qtree. %v", err)
		return errors.New("Volume mount failed.")
	}
	if !exists {
		log.WithField("qtree", name).Debug("Qtree not found.")
		return fmt.Errorf("Volume %s not found.", name)
	}

	exportPath := fmt.Sprintf("%s:/%s/%s", d.Config.DataLIF, flexvol, name)

	return MountVolume(exportPath, mountpoint, &d.Config)
}

// Detach the volume
func (d *OntapNASQtreeStorageDriver) Detach(name, mountpoint string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Detach",
			"Type":       "OntapNASQtreeStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
		}
		log.WithFields(fields).Debug(">>>> Detach")
		defer log.WithFields(fields).Debug("<<<< Detach")
	}

	exists, _, err := d.API.QtreeExists(name, d.FlexvolNamePrefix())
	if err != nil {
		log.Warnf("Error checking for existing qtree. %v", err)
	}
	if !exists {
		log.WithField("qtree", name).Warn("Qtree not found, attempting unmount anyway.")
	}

	return UnmountVolume(mountpoint, &d.Config)
}

// Return the list of snapshots associated with the named volume
func (d *OntapNASQtreeStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "SnapshotList",
			"Type":   "OntapNASQtreeStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> SnapshotList")
		defer log.WithFields(fields).Debug("<<<< SnapshotList")
	}

	// Qtrees can't have snapshots, so return an empty list
	return []CommonSnapshot{}, nil
}

// Return the list of volumes associated with this tenant
func (d *OntapNASQtreeStorageDriver) List() ([]string, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "List", "Type": "OntapNASQtreeStorageDriver"}
		log.WithFields(fields).Debug(">>>> List")
		defer log.WithFields(fields).Debug("<<<< List")
	}

	// Generic user-facing message
	listError := errors.New("Volume list failed.")

	prefix := *d.Config.StoragePrefix
	volumes := []string{}

	// Get all qtrees in all Flexvols managed by this driver
	listResponse, err := d.API.QtreeList(prefix, d.FlexvolNamePrefix())
	if err = ontap.GetError(listResponse, err); err != nil {
		log.Errorf("Qtree list failed. %v", err)
		return volumes, listError
	}

	// AttributesList() returns []QtreeInfoType
	for _, qtree := range listResponse.Result.AttributesList() {
		vol := qtree.Qtree()[len(prefix):]
		volumes = append(volumes, vol)
	}

	return volumes, nil
}

// Test for the existence of a volume
func (d *OntapNASQtreeStorageDriver) Get(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Get", "Type": "OntapNASQtreeStorageDriver"}
		log.WithFields(fields).Debug(">>>> Get")
		defer log.WithFields(fields).Debug("<<<< Get")
	}

	// Generic user-facing message
	getError := fmt.Errorf("Volume %s not found.", name)

	exists, flexvol, err := d.API.QtreeExists(name, d.FlexvolNamePrefix())
	if err != nil {
		log.Errorf("Error checking for existing qtree. %v", err)
		return getError
	}
	if !exists {
		log.WithField("qtree", name).Debug("Qtree not found.")
		return getError
	}

	log.WithFields(log.Fields{"qtree": name, "flexvol": flexvol}).Debug("Qtree found.")

	return nil
}

// ensureFlexvolForQtree accepts a set of Flexvol characteristics and either finds one to contain a new
// qtree or it creates a new Flexvol with the needed attributes.
func (d *OntapNASQtreeStorageDriver) ensureFlexvolForQtree(
	aggregate, spaceReserve, snapshotPolicy string, enableSnapshotDir bool, encrypt *bool,
) (string, error) {

	// Check if a suitable Flexvol already exists
	flexvol, err := d.getFlexvolForQtree(aggregate, spaceReserve, snapshotPolicy, enableSnapshotDir, encrypt)
	if err != nil {
		return "", fmt.Errorf("Error finding Flexvol for qtree. %v", err)
	}

	// Found one!
	if flexvol != "" {
		return flexvol, nil
	}

	// Nothing found, so create a suitable Flexvol
	flexvol, err = d.createFlexvolForQtree(aggregate, spaceReserve, snapshotPolicy, enableSnapshotDir, encrypt)
	if err != nil {
		return "", fmt.Errorf("Error creating Flexvol for qtree. %v", err)
	}

	return flexvol, nil
}

// createFlexvolForQtree creates a new Flexvol matching the specified attributes for
// the purpose of containing qtrees supplied as container volumes by this driver.
// Once this method returns, the Flexvol exists, is mounted, and has a default tree
// quota.
func (d *OntapNASQtreeStorageDriver) createFlexvolForQtree(
	aggregate, spaceReserve, snapshotPolicy string, enableSnapshotDir bool, encrypt *bool,
) (string, error) {

	flexvol := d.FlexvolNamePrefix() + utils.RandomString(10)
	size := "1g"
	unixPermissions := "0700"
	exportPolicy := d.flexvolExportPolicy
	securityStyle := "unix"

	encryption := false
	if encrypt != nil {
		encryption = *encrypt
	}

	log.WithFields(log.Fields{
		"name":            flexvol,
		"aggregate":       aggregate,
		"size":            size,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
		"snapshotDir":     enableSnapshotDir,
		"exportPolicy":    exportPolicy,
		"securityStyle":   securityStyle,
		"encryption":      encryption,
	}).Debug("Creating Flexvol for qtrees.")

	// Create the Flexvol
	createResponse, err := d.API.VolumeCreate(
		flexvol, aggregate, size, spaceReserve, snapshotPolicy,
		unixPermissions, exportPolicy, securityStyle, encrypt)
	if err = ontap.GetError(createResponse, err); err != nil {
		return "", fmt.Errorf("Error creating Flexvol. %v", err)
	}

	// Disable '.snapshot' as needed
	if !enableSnapshotDir {
		snapDirResponse, err := d.API.VolumeDisableSnapshotDirectoryAccess(flexvol)
		if err = ontap.GetError(snapDirResponse, err); err != nil {
			defer d.API.VolumeDestroy(flexvol, true)
			return "", fmt.Errorf("Error disabling snapshot directory access. %v", err)
		}
	}

	// Mount the volume at the specified junction
	mountResponse, err := d.API.VolumeMount(flexvol, "/"+flexvol)
	if err = ontap.GetError(mountResponse, err); err != nil {
		defer d.API.VolumeDestroy(flexvol, true)
		return "", fmt.Errorf("Error mounting Flexvol. %v", err)
	}

	// If LS mirrors are present on the SVM root volume, update them
	UpdateLoadSharingMirrors(d.API)

	// Create the default quota rule so we can use quota-resize for new qtrees
	err = d.addDefaultQuotaForFlexvol(flexvol)
	if err != nil {
		defer d.API.VolumeDestroy(flexvol, true)
		return "", fmt.Errorf("Error adding default quota to Flexvol. %v", err)
	}

	return flexvol, nil
}

// getFlexvolForQtree returns a Flexvol (from the set of existing Flexvols) that
// matches the specified Flexvol attributes and does not already contain more
// than the maximum configured number of qtrees.  No matching Flexvols is not
// considered an error.  If more than one matching Flexvol is found, one of those
// is returned at random.
func (d *OntapNASQtreeStorageDriver) getFlexvolForQtree(
	aggregate, spaceReserve, snapshotPolicy string, enableSnapshotDir bool, encrypt *bool,
) (string, error) {

	// Get all volumes matching the specified attributes
	volListResponse, err := d.API.VolumeListByAttrs(
		d.FlexvolNamePrefix(), aggregate, spaceReserve, snapshotPolicy, enableSnapshotDir, encrypt)

	if err = ontap.GetError(volListResponse, err); err != nil {
		return "", fmt.Errorf("Error enumerating Flexvols. %v", err)
	}

	// Weed out the Flexvols already having too many qtrees
	var volumes []string
	for _, volAttrs := range volListResponse.Result.AttributesList() {
		volIdAttrs := volAttrs.VolumeIdAttributes()
		volName := string(volIdAttrs.Name())

		count, err := d.API.QtreeCount(volName)
		if err != nil {
			return "", fmt.Errorf("Error enumerating qtrees. %v", err)
		}

		if count < maxQtreesPerFlexvol {
			volumes = append(volumes, volName)
		}
	}

	// Pick a Flexvol.  If there are multiple matches, pick one at random.
	switch len(volumes) {
	case 0:
		return "", nil
	case 1:
		return volumes[0], nil
	default:
		rand.Seed(time.Now().UnixNano())
		return volumes[rand.Intn(len(volumes))], nil
	}
}

// getOptimalSizeForFlexvol sums up all the disk limit quota rules on a Flexvol and adds the size of
// the new qtree being added as well as the current Flexvol snapshot reserve.  This value may be used
// to grow (or shrink) the Flexvol as new qtrees are being added.
func (d *OntapNASQtreeStorageDriver) getOptimalSizeForFlexvol(
	flexvol string, newQtreeSizeBytes uint64,
) (uint64, error) {

	// Get more info about the Flexvol
	volAttrs, err := d.API.VolumeGet(flexvol)
	if err != nil {
		return 0, err
	}
	volSpaceAttrs := volAttrs.VolumeSpaceAttributes()
	snapReserveMultiplier := 1.0 + (float64(volSpaceAttrs.PercentageSnapshotReserve()) / 100.0)

	totalDiskLimitBytes, err := d.getTotalHardDiskLimitQuota(flexvol)
	if err != nil {
		return 0, err
	}

	usableSpaceBytes := float64(newQtreeSizeBytes + totalDiskLimitBytes)
	flexvolSizeBytes := uint64(usableSpaceBytes * snapReserveMultiplier)

	log.WithFields(log.Fields{
		"flexvol":               flexvol,
		"snapReserveMultiplier": snapReserveMultiplier,
		"totalDiskLimitBytes":   totalDiskLimitBytes,
		"newQtreeSizeBytes":     newQtreeSizeBytes,
		"flexvolSizeBytes":      flexvolSizeBytes,
	}).Debug("Calculated optimal size for Flexvol with new qtree.")

	return flexvolSizeBytes, nil
}

// addDefaultQuotaForFlexvol adds a default quota rule to a Flexvol so that quotas for
// new qtrees may be added on demand with simple quota resize instead of a heavyweight
// quota reinitialization.
func (d *OntapNASQtreeStorageDriver) addDefaultQuotaForFlexvol(flexvol string) error {

	response, err := d.API.QuotaSetEntry("", flexvol, "", "tree", "-")
	if err = ontap.GetError(response, err); err != nil {
		return fmt.Errorf("Error adding default quota. %v", err)
	}

	d.disableQuotas(flexvol, true)
	if err != nil {
		return fmt.Errorf("Error adding default quota. %v", err)
	}

	d.enableQuotas(flexvol, true)
	if err != nil {
		return fmt.Errorf("Error adding default quota. %v", err)
	}

	return nil
}

// addQuotaForQtree adds a tree quota to a Flexvol/qtree with a hard disk size limit.
func (d *OntapNASQtreeStorageDriver) addQuotaForQtree(qtree, flexvol string, sizeBytes uint64) error {

	target := fmt.Sprintf("/vol/%s/%s", flexvol, qtree)
	sizeKB := strconv.FormatUint(sizeBytes/1024, 10)

	response, err := d.API.QuotaSetEntry("", flexvol, target, "tree", sizeKB)
	if err = ontap.GetError(response, err); err != nil {
		return fmt.Errorf("Error adding qtree quota. %v", err)
	}

	// Mark this Flexvol as needing a quota resize
	d.quotaResizeMap[flexvol] = true

	return nil
}

// enableQuotas disables quotas on a Flexvol, optionally waiting for the operation to finish.
func (d *OntapNASQtreeStorageDriver) disableQuotas(flexvol string, wait bool) error {

	status, err := d.getQuotaStatus(flexvol)
	if err != nil {
		return fmt.Errorf("Error disabling quotas. %v.", err)
	}
	if status == "corrupt" {
		return fmt.Errorf("Error disabling quotas. Quotas are corrupt on Flexvol %s.", flexvol)
	}

	if status != "off" {
		offResponse, err := d.API.QuotaOff(flexvol)
		if err = ontap.GetError(offResponse, err); err != nil {
			return fmt.Errorf("Error disabling quotas. %v.", err)
		}
	}

	if wait {
		for status != "off" {
			time.Sleep(1 * time.Second)

			status, err = d.getQuotaStatus(flexvol)
			if err != nil {
				return fmt.Errorf("Error disabling quotas. %v.", err)
			}
			if status == "corrupt" {
				return fmt.Errorf("Error disabling quotas. Quotas are corrupt on Flexvol %s.", flexvol)
			}
		}
	}

	return nil
}

// enableQuotas enables quotas on a Flexvol, optionally waiting for the operation to finish.
func (d *OntapNASQtreeStorageDriver) enableQuotas(flexvol string, wait bool) error {

	status, err := d.getQuotaStatus(flexvol)
	if err != nil {
		return fmt.Errorf("Error enabling quotas. %v.", err)
	}
	if status == "corrupt" {
		return fmt.Errorf("Error enabling quotas. Quotas are corrupt on Flexvol %s.", flexvol)
	}

	if status == "off" {
		onResponse, err := d.API.QuotaOn(flexvol)
		if err = ontap.GetError(onResponse, err); err != nil {
			return fmt.Errorf("Error enabling quotas. %v.", err)
		}
	}

	if wait {
		for status != "on" {
			time.Sleep(1 * time.Second)

			status, err = d.getQuotaStatus(flexvol)
			if err != nil {
				return fmt.Errorf("Error enabling quotas. %v.", err)
			}
			if status == "corrupt" {
				return fmt.Errorf("Error enabling quotas. Quotas are corrupt on Flexvol %s.", flexvol)
			}
		}
	}

	return nil
}

// queueAllFlexvolsForQuotaResize flags every Flexvol managed by this driver as
// needing a quota resize.  This is called once on driver startup to handle the
// case where the driver was shut down with pending quota resize operations.
func (d *OntapNASQtreeStorageDriver) queueAllFlexvolsForQuotaResize() {

	// Get list of Flexvols managed by this driver
	volumeListResponse, err := d.API.VolumeList(d.FlexvolNamePrefix())
	if err = ontap.GetError(volumeListResponse, err); err != nil {
		log.Errorf("Error listing Flexvols. %v", err)
	}

	for _, volAttrs := range volumeListResponse.Result.AttributesList() {
		volIdAttrs := volAttrs.VolumeIdAttributes()
		flexvol := string(volIdAttrs.Name())
		d.quotaResizeMap[flexvol] = true
	}
}

// resizeQuotas may be called by a background task, or by a method that changed
// the qtree population on a Flexvol.  Flexvols needing an update must be flagged
// in quotaResizeMap.  Any failures that occur are simply logged, and the resize
// operation will be attempted each time this method is called until it succeeds.
func (d *OntapNASQtreeStorageDriver) resizeQuotas() {

	// Ensure we don't forget any Flexvol that is involved in a qtree provisioning workflow
	d.provMutex.Lock()
	defer d.provMutex.Unlock()

	log.Debug("Housekeeping, resizing quotas.")

	for flexvol, resize := range d.quotaResizeMap {

		if resize {
			resizeResponse, err := d.API.QuotaResize(flexvol)
			if err != nil {
				log.WithFields(log.Fields{"flexvol": flexvol, "error": err}).Debug("Error resizing quotas.")
				continue
			}
			if zerr := ontap.NewZapiError(resizeResponse); !zerr.IsPassed() {

				if zerr.Code() == azgo.EVOLUMEDOESNOTEXIST {
					// Volume gone, so no need to try again
					log.WithField("flexvol", flexvol).Debug("Volume does not exist.")
					delete(d.quotaResizeMap, flexvol)
				} else {
					log.WithFields(log.Fields{"flexvol": flexvol, "error": zerr}).Debug("Error resizing quotas.")
				}

				continue
			}

			log.WithField("flexvol", flexvol).Debug("Started quota resize.")

			// Resize start succeeded, so no need to try again
			delete(d.quotaResizeMap, flexvol)
		}
	}
}

// getQuotaStatus returns the status of the quotas on a Flexvol
func (d *OntapNASQtreeStorageDriver) getQuotaStatus(flexvol string) (string, error) {

	statusResponse, err := d.API.QuotaStatus(flexvol)
	if err = ontap.GetError(statusResponse, err); err != nil {
		return "", fmt.Errorf("Error getting quota status for Flexvol %s. %v", flexvol, err)
	}

	return statusResponse.Result.Status(), nil

}

// getTotalHardDiskLimitQuota returns the sum of all disk limit quota rules on a Flexvol
func (d *OntapNASQtreeStorageDriver) getTotalHardDiskLimitQuota(flexvol string) (uint64, error) {

	listResponse, err := d.API.QuotaEntryList(flexvol)
	if err != nil {
		return 0, err
	}

	var totalDiskLimitKB uint64 = 0

	for _, rule := range listResponse.Result.AttributesList() {
		diskLimitKB, err := strconv.ParseUint(rule.DiskLimit(), 10, 64)
		if err != nil {
			continue
		}
		totalDiskLimitKB += diskLimitKB
	}

	return totalDiskLimitKB * 1024, nil
}

// pruneUnusedFlexvols is called periodically by a background task.  Any Flexvols
// that are managed by this driver (discovered by virtue of having a well-known
// hardcoded prefix on their names) that have no qtrees are deleted.
func (d *OntapNASQtreeStorageDriver) pruneUnusedFlexvols() {

	// Ensure we don't prune any Flexvol that is involved in a qtree provisioning workflow
	d.provMutex.Lock()
	defer d.provMutex.Unlock()

	log.Debug("Housekeeping, checking for managed Flexvols with no qtrees.")

	// Get list of Flexvols managed by this driver
	volumeListResponse, err := d.API.VolumeList(d.FlexvolNamePrefix())
	if err = ontap.GetError(volumeListResponse, err); err != nil {
		log.Errorf("Error listing Flexvols. %v", err)
	}

	var flexvols []string
	for _, volAttrs := range volumeListResponse.Result.AttributesList() {
		volIdAttrs := volAttrs.VolumeIdAttributes()
		volName := string(volIdAttrs.Name())
		flexvols = append(flexvols, volName)
	}

	// Destroy any Flexvol if it is devoid of qtrees
	for _, flexvol := range flexvols {
		qtreeCount, err := d.API.QtreeCount(flexvol)
		if err == nil && qtreeCount == 0 {
			log.WithField("flexvol", flexvol).Debug("Housekeeping, deleting managed Flexvol with no qtrees.")
			d.API.VolumeDestroy(flexvol, true)
		}
	}
}

// reapDeletedQtrees is called periodically by a background task.  Any qtrees
// that have been deleted (discovered by virtue of having a well-known hardcoded
// prefix on their names) are destroyed.  This is only needed for the exceptional case
// in which a qtree was renamed (prior to being destroyed) but the subsequent
// destroy call failed or was never made due to a process interruption.
func (d *OntapNASQtreeStorageDriver) reapDeletedQtrees() {

	log.Debug("Housekeeping, checking for deleted qtrees.")

	// Get all deleted qtrees in all FlexVols managed by this driver
	prefix := deletedQtreeNamePrefix + *d.Config.StoragePrefix
	listResponse, err := d.API.QtreeList(prefix, d.FlexvolNamePrefix())
	if err = ontap.GetError(listResponse, err); err != nil {
		log.Errorf("Error listing deleted qtrees. %v", err)
	}

	// AttributesList() returns []QtreeInfoType
	for _, qtree := range listResponse.Result.AttributesList() {
		qtreePath := fmt.Sprintf("/vol/%s/%s", qtree.Volume(), qtree.Qtree())
		log.WithField("qtree", qtreePath).Debug("Housekeeping, reaping deleted qtree.")
		d.API.QtreeDestroyAsync(qtreePath, true)
	}
}

// ensureDefaultExportPolicy checks for an export policy with a well-known name that will be suitable
// for setting on a Flexvol and will enable access to all qtrees therein.  If the policy exists, the
// method assumes it created the policy itself and that all is good.  If the policy does not exist,
// it is created and populated with a rule that allows access to NFS qtrees.  This method should be
// called once during driver initialization.
func (d *OntapNASQtreeStorageDriver) ensureDefaultExportPolicy() error {

	policyResponse, err := d.API.ExportPolicyCreate(d.flexvolExportPolicy)
	if err != nil {
		return fmt.Errorf("Error creating export policy %s. %v", d.flexvolExportPolicy, err)
	}
	if zerr := ontap.NewZapiError(policyResponse); !zerr.IsPassed() {
		if zerr.Code() == azgo.EDUPLICATEENTRY {
			log.WithField("exportPolicy", d.flexvolExportPolicy).Debug("Export policy already exists.")
		} else {
			return fmt.Errorf("Error creating export policy %s. %v", d.flexvolExportPolicy, zerr)
		}
	}

	return d.ensureDefaultExportPolicyRule()
}

// ensureDefaultExportPolicyRule guarantees that the export policy used on Flexvols managed by this
// driver has at least one rule, which is necessary (but not always sufficient) to enable qtrees
// to be mounted by clients.
func (d *OntapNASQtreeStorageDriver) ensureDefaultExportPolicyRule() error {

	ruleListResponse, err := d.API.ExportRuleGetIterRequest(d.flexvolExportPolicy)
	if err = ontap.GetError(ruleListResponse, err); err != nil {
		return fmt.Errorf("Error listing export policy rules. %v", err)
	}

	if ruleListResponse.Result.NumRecords() == 0 {

		// No rules, so create one
		ruleResponse, err := d.API.ExportRuleCreate(
			d.flexvolExportPolicy, "0.0.0.0/0",
			[]string{"nfs"}, []string{"any"}, []string{"any"}, []string{"any"})
		if err = ontap.GetError(ruleResponse, err); err != nil {
			return fmt.Errorf("Error creating export rule. %v", err)
		}
	} else {
		log.WithField("exportPolicy", d.flexvolExportPolicy).Debug("Export policy has at least one rule.")
	}

	return nil
}
