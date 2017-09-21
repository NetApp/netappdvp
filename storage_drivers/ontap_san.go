// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/ontap"
	"github.com/netapp/netappdvp/azgo"
	"github.com/netapp/netappdvp/utils"
)

// OntapSANStorageDriverName is the constant name for this Ontap NAS storage driver
const OntapSANStorageDriverName = "ontap-san"
const OntapLUNAttributeFstype = "com.netapp.ndvp.fstype"

func init() {
	san := &OntapSANStorageDriver{}
	san.Initialized = false
	Drivers[san.Name()] = san
}

func lunPath(name string) string {
	return fmt.Sprintf("/vol/%v/lun0", name)
}

// OntapSANStorageDriver is for iSCSI storage provisioning
type OntapSANStorageDriver struct {
	Initialized bool
	Config      OntapStorageDriverConfig
	API         *ontap.Driver
}

func (d *OntapSANStorageDriver) GetConfig() *OntapStorageDriverConfig {
	return &d.Config
}

func (d *OntapSANStorageDriver) GetAPI() *ontap.Driver {
	return d.API
}

// Name is for returning the name of this driver
func (d OntapSANStorageDriver) Name() string {
	return OntapSANStorageDriverName
}

// Initialize from the provided config
func (d *OntapSANStorageDriver) Initialize(
	context DriverContext, configJSON string, commonConfig *CommonStorageDriverConfig,
) error {

	if commonConfig.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Initialize", "Type": "OntapSANStorageDriver"}
		log.WithFields(fields).Debug(">>>> Initialize")
		defer log.WithFields(fields).Debug("<<<< Initialize")
	}

	// Parse the config
	config, err := InitializeOntapConfig(configJSON, commonConfig)
	if err != nil {
		return fmt.Errorf("Error initializing %s driver. %v", d.Name(), err)
	}

	if config.IgroupName == "" {
		config.IgroupName = "netappdvp"
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
func (d *OntapSANStorageDriver) Validate(context DriverContext) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Validate", "Type": "OntapSANStorageDriver", "context": context}
		log.WithFields(fields).Debug(">>>> Validate")
		defer log.WithFields(fields).Debug("<<<< Validate")
	}

	lifResponse, err := d.API.NetInterfaceGet()
	if err = ontap.GetError(lifResponse, err); err != nil {
		return fmt.Errorf("Error checking network interfaces. %v", err)
	}

	// If they didn't set a LIF to use in the config, we'll set it to the first iSCSI LIF we happen to find
	if d.Config.DataLIF == "" {
	loop:
		for _, attrs := range lifResponse.Result.AttributesList() {
			for _, protocol := range attrs.DataProtocols() {
				if protocol == "iscsi" {
					log.WithField("address", attrs.Address()).Debug("Choosing LIF for iSCSI.")
					d.Config.DataLIF = string(attrs.Address())
					break loop
				}
			}
		}
	}

	// Validate our settings
	foundIscsi := false
	iscsiLifCount := 0
	for _, attrs := range lifResponse.Result.AttributesList() {
		for _, protocol := range attrs.DataProtocols() {
			if protocol == "iscsi" {
				log.Debugf("Comparing iSCSI protocol access on %v vs. %v", attrs.Address(), d.Config.DataLIF)
				if string(attrs.Address()) == d.Config.DataLIF {
					foundIscsi = true
					iscsiLifCount++
				}
			}
		}
	}

	if iscsiLifCount > 1 {
		log.Debugf("Found multiple iSCSI LIFs.")
	}

	if !foundIscsi {
		return fmt.Errorf("Could not find iSCSI data LIF.")
	}

	if context == ContextNDVP {
		// Make sure this host is logged into the ONTAP iSCSI target
		err := utils.EnsureIscsiSession(d.Config.DataLIF)
		if err != nil {
			return fmt.Errorf("Error establishing iSCSI session. %v", err)
		}

		// Make sure the configured aggregate is available
		err = ValidateAggregate(d.API, &d.Config)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create a volume+LUN with the specified options
func (d *OntapSANStorageDriver) Create(name string, sizeBytes uint64, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":    "Create",
			"Type":      "OntapSANStorageDriver",
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

	// Get options with default fallback values
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

	encrypt, err := ValidateEncryptionAttribute(encryption, d.API)
	if err != nil {
		return err
	}

	// Check for a supported file system type
	fstype := strings.ToLower(utils.GetV(opts, "fstype|fileSystemType", d.Config.FileSystemType))
	switch fstype {
	case "xfs", "ext3", "ext4":
		log.WithFields(log.Fields{"fileSystemType": fstype, "name": name}).Debug("Filesystem format.")
	default:
		return fmt.Errorf("Unsupported fileSystemType option: %s.", fstype)
	}

	log.WithFields(log.Fields{
		"name":            name,
		"size":            size,
		"spaceReserve":    spaceReserve,
		"snapshotPolicy":  snapshotPolicy,
		"unixPermissions": unixPermissions,
		"snapshotDir":     snapshotDir,
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

	lunPath := lunPath(name)
	osType := "linux"
	spaceReserved := false

	// Create the LUN
	lunCreateResponse, err := d.API.LunCreate(lunPath, int(sizeBytes), osType, spaceReserved)
	if err = ontap.GetError(lunCreateResponse, err); err != nil {
		return fmt.Errorf("Error creating LUN. %v", err)
	}

	// Save the fstype in a LUN attribute so we know what to do in Attach
	attrResponse, err := d.API.LunSetAttribute(lunPath, OntapLUNAttributeFstype, fstype)
	if err = ontap.GetError(attrResponse, err); err != nil {
		defer d.API.LunDestroy(lunPath)
		return fmt.Errorf("Error saving file system type for LUN. %v", err)
	}

	return nil
}

// Create a volume clone
func (d *OntapSANStorageDriver) CreateClone(name, source, snapshot string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":   "CreateClone",
			"Type":     "OntapSANStorageDriver",
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
		return fmt.Errorf("Invalid boolean value for splitOnClone: %v", err)
	}

	log.WithField("splitOnClone", split).Debug("Creating volume clone.")
	return CreateOntapClone(name, source, snapshot, split, &d.Config, d.API)
}

// Destroy the requested (volume,lun) storage tuple
func (d *OntapSANStorageDriver) Destroy(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Destroy",
			"Type":   "OntapSANStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> Destroy")
		defer log.WithFields(fields).Debug("<<<< Destroy")
	}

	lunPath := lunPath(name)

	// Validate Flexvol exists before trying to destroy
	volExists, err := d.API.VolumeExists(name)
	if err != nil {
		return fmt.Errorf("Error checking for existing volume. %v", err)
	}
	if !volExists {
		log.WithField("volume", name).Debug("Volume already deleted, skipping destroy.")
		return nil
	}

	// Set LUN offline
	offlineResponse, err := d.API.LunOffline(lunPath)
	if err := ontap.GetError(offlineResponse, err); err != nil {
		log.Warnf("Error attempting to offline LUN %v. %v", lunPath, err)
	}

	// Destroy LUN
	lunDestroyResponse, err := d.API.LunDestroy(lunPath)
	if err := ontap.GetError(lunDestroyResponse, err); err != nil {
		log.Warnf("Error destroying LUN %v. %v", lunPath, err)
	}

	// Perform rediscovery to remove the deleted LUN
	utils.MultipathFlush() // flush unused paths
	utils.IscsiRescan(true)

	// Delete the Flexvol
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

// Attach the lun
func (d *OntapSANStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Attach",
			"Type":       "OntapSANStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
			"opts":       opts,
		}
		log.WithFields(fields).Debug(">>>> Attach")
		defer log.WithFields(fields).Debug("<<<< Attach")
	}

	// Error if no iSCSI session exists for the specified iscsi portal
	sessionExists, err := utils.IscsiSessionExists(d.Config.DataLIF)
	if err != nil {
		return fmt.Errorf("Unexpected iSCSI session error. %v", err)
	}
	if !sessionExists {
		// TODO automatically login for the user if no session detected?
		return fmt.Errorf("Expected iSCSI session %v not found, please login to the iSCSI portal.", d.Config.DataLIF)
	}

	igroupName := d.Config.IgroupName
	lunPath := lunPath(name)

	// Get the fstype
	fstype := DefaultFileSystemType
	attrResponse, err := d.API.LunGetAttribute(lunPath, OntapLUNAttributeFstype)
	if err = ontap.GetError(attrResponse, err); err != nil {
		log.WithFields(log.Fields{
			"LUN":    lunPath,
			"fstype": fstype,
		}).Warn("LUN attribute fstype not found, using default.")
	} else {
		fstype = attrResponse.Result.Value()
		log.WithFields(log.Fields{"LUN": lunPath, "fstype": fstype}).Debug("Found LUN attribute fstype.")
	}

	// Create igroup
	igroupResponse, err := d.API.IgroupCreate(igroupName, "iscsi", "linux")
	if err != nil {
		return fmt.Errorf("Error creating igroup. %v", err)
	}
	if zerr := ontap.NewZapiError(igroupResponse); !zerr.IsPassed() {
		// Handle case where the igroup already exists
		if zerr.Code() != azgo.EVDISK_ERROR_INITGROUP_EXISTS {
			return fmt.Errorf("Error creating igroup %v. %v", igroupName, zerr)
		}
	}

	// Lookup host IQNs
	iqns, err := utils.GetInitiatorIqns()
	if err != nil {
		return fmt.Errorf("Error determining host initiator IQNs. %v", err)
	}

	// Add each IQN found to group
	for _, iqn := range iqns {
		igroupAddResponse, err := d.API.IgroupAdd(igroupName, iqn)
		if err := ontap.GetError(igroupAddResponse, err); err != nil {
			if zerr, ok := err.(ontap.ZapiError); ok {
				if zerr.Code() == azgo.EVDISK_ERROR_INITGROUP_HAS_NODE {
					continue
				}
			}
			return fmt.Errorf("Error adding IQN %v to igroup %v. %v", iqn, igroupName, err)
		}
	}

	// Map LUN
	lunID, err := d.API.LunMapIfNotMapped(igroupName, lunPath)
	if err != nil {
		return err
	}

	// Perform discovery to see the created/mapped LUN
	utils.IscsiRescan(false)

	cmd := fmt.Sprintf("dmesg | tail -n 50")
	log.Debugf("running 'sh -c %v'", cmd)
	out2, _ := exec.Command("sh", "-c", cmd).CombinedOutput()
	log.Debug(string(out2))

	// Lookup all the SCSI device information
	info, err := utils.GetDeviceInfoForLuns()
	if err != nil {
		return fmt.Errorf("Error getting SCSI device information. %v", err)
	}

	// Lookup all the iSCSI session information
	sessionInfo, err := utils.GetIscsiSessionInfo()
	if err != nil {
		return fmt.Errorf("Error getting iSCSI session information. %v", err)
	}

	sessionInfoToUse := utils.IscsiSessionInfo{}
	for i, e := range sessionInfo {
		if e.PortalIP == d.Config.DataLIF {
			sessionInfoToUse = sessionInfo[i]
		}
	}

	for i, e := range info {
		log.WithFields(log.Fields{
			"i":                i,
			"scsiHost":         e.Host,
			"scsiChannel":      e.Channel,
			"scsiTarget":       e.Target,
			"scsiLun":          e.LUN,
			"multipathDevFile": e.MultipathDevice,
			"devFile":          e.Device,
			"fsType":           e.Filesystem,
			"iqn":              e.IQN,
		}).Debug("Found")
	}

	// look for the expected mapped lun
	for i, e := range info {

		log.WithFields(log.Fields{
			"i":                i,
			"scsiHost":         e.Host,
			"scsiChannel":      e.Channel,
			"scsiTarget":       e.Target,
			"scsiLun":          e.LUN,
			"multipathDevFile": e.MultipathDevice,
			"devFile":          e.Device,
			"fsType":           e.Filesystem,
			"iqn":              e.IQN,
		}).Debug("Checking")

		if e.LUN != strconv.Itoa(lunID) {
			log.Debugf("Skipping... lun id %v != %v", e.LUN, lunID)
			continue
		}

		if !strings.HasPrefix(e.IQN, sessionInfoToUse.TargetName) {
			log.Debugf("Skipping... %v doesn't start with %v", e.IQN, sessionInfoToUse.TargetName)
			continue
		}

		// If we're here then, we should be on the right info element:
		// *) we have the expected LUN ID
		// *) we have the expected iscsi session target
		log.Debugf("Using... %v", e)

		// Make sure we use the proper device (multipath if in use)
		deviceToUse := e.Device
		if e.MultipathDevice != "" {
			deviceToUse = e.MultipathDevice
		}

		if deviceToUse == "" {
			return fmt.Errorf("Could not determine device to use for %v.", name)
		}

		// Put a filesystem on it if there isn't one already there
		if e.Filesystem == "" {
			log.WithFields(log.Fields{"LUN": lunPath, "fstype": fstype}).Debug("Formatting LUN.")
			err := utils.FormatVolume(deviceToUse, fstype)
			if err != nil {
				return fmt.Errorf("Error formatting LUN %v, device %v. %v", name, deviceToUse, err)
			}
		} else if e.Filesystem != fstype {
			log.WithFields(log.Fields{
				"LUN":             lunPath,
				"existingFstype":  e.Filesystem,
				"requestedFstype": fstype,
			}).Warn("LUN already formatted with a different file system type.")
		} else {
			log.WithFields(log.Fields{"LUN": lunPath, "fstype": e.Filesystem}).Debug("LUN already formatted.")
		}

		// Mount it
		err := utils.Mount(deviceToUse, mountpoint)
		if err != nil {
			return fmt.Errorf("Error mounting LUN %v, device %v, mountpoint %v. %v", name, deviceToUse, mountpoint, err)
		}
		return nil
	}

	return fmt.Errorf("Attach failed, device not found. %v", name)
}

// Detach the volume
func (d *OntapSANStorageDriver) Detach(name, mountpoint string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "Detach",
			"Type":       "OntapSANStorageDriver",
			"name":       name,
			"mountpoint": mountpoint,
		}
		log.WithFields(fields).Debug(">>>> Detach")
		defer log.WithFields(fields).Debug("<<<< Detach")
	}

	cmd := fmt.Sprintf("umount %s", mountpoint)
	log.WithField("command", cmd).Debug("Unmounting volume.")

	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		log.WithField("output", string(out)).Debug("Unmount failed.")
		return fmt.Errorf("Error unmounting volume %v, mountpoint %v. %v", name, mountpoint, err)
	}

	return nil
}

// Return the list of snapshots associated with the named volume
func (d *OntapSANStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "SnapshotList",
			"Type":   "OntapSANStorageDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> SnapshotList")
		defer log.WithFields(fields).Debug("<<<< SnapshotList")
	}

	return GetSnapshotList(name, &d.Config, d.API)
}

// Return the list of volumes associated with this tenant
func (d *OntapSANStorageDriver) List() ([]string, error) {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "List", "Type": "OntapSANStorageDriver"}
		log.WithFields(fields).Debug(">>>> List")
		defer log.WithFields(fields).Debug("<<<< List")
	}

	return GetVolumeList(d.API, &d.Config)
}

// Test for the existence of a volume
func (d *OntapSANStorageDriver) Get(name string) error {

	if d.Config.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "Get", "Type": "OntapSANStorageDriver"}
		log.WithFields(fields).Debug(">>>> Get")
		defer log.WithFields(fields).Debug("<<<< Get")
	}

	return GetVolume(name, d.API, &d.Config)
}
