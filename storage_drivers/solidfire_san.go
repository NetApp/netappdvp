// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/units"
	"github.com/ebalduf/netappdvp/apis/sfapi"
	"github.com/ebalduf/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

// SolidfireSANStorageDriverName is the constant name for this Solidfire SAN storage driver
const SolidfireSANStorageDriverName = "solidfire-san"

func init() {
	san := &SolidfireSANStorageDriver{}
	san.Initialized = false
	Drivers[san.Name()] = san
	log.Debugf("Registered driver '%v'", san.Name())
}

func formatOpts(opts map[string]string) {
	// NOTE(jdg): For now we just want to minimize issues like case usage for
	// the two basic opts most used (size and type).  Going forward we can add
	// all sorts of things here based on what we decide to add as valid opts
	// during create and even other calls
	for k, v := range opts {
		if strings.EqualFold(k, "size") {
			opts["size"] = v
		} else if strings.EqualFold(k, "type") {
			opts["type"] = v
		} else if strings.EqualFold(k, "qos") {
			opts["qos"] = v
		}
	}
}

// SolidfireSANStorageDriver is for iSCSI storage provisioning
type SolidfireSANStorageDriver struct {
	Initialized    bool
	Config         SolidfireStorageDriverConfig
	Client         *sfapi.Client
	TenantID       int64
	DefaultVolSz   int64
	VagID          int64
	InitiatorIFace string
}

// Name is for returning the name of this driver
func (d SolidfireSANStorageDriver) Name() string {
	return SolidfireSANStorageDriverName
}

// Initialize from the provided config
func (d *SolidfireSANStorageDriver) Initialize(configJSON string) error {
	log.Debugf("SolidfireSANStorageDriver#Initialize(...)")

	c := &SolidfireStorageDriverConfig{}

	// decode supplied configJSON string into SolidfireStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &c)
	if err != nil {
		return fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	log.WithFields(log.Fields{
		"Version":           c.Version,
		"StorageDriverName": c.StorageDriverName,
		"Debug":             c.Debug,
		"DisableDelete":     c.DisableDelete,
		"StoragePrefixRaw":  string(c.StoragePrefixRaw),
		"SnapshotPrefixRaw": string(c.SnapshotPrefixRaw),
	}).Debugf("Reparsed into solidfireConfig")

	c.DefaultVolSz = c.DefaultVolSz * int64(units.GiB)
	log.Debugf("Decoded to %v", c)
	d.Config = *c

	var tenantID int64

	// create a new sfapi.Config object from the read in json config file
	endpoint := c.EndPoint
	defaultSizeGiB := c.DefaultVolSz
	svip := c.SVIP
	cfg := sfapi.Config{
		TenantName:     c.TenantName,
		EndPoint:       c.EndPoint,
		DefaultVolSz:   c.DefaultVolSz,
		SVIP:           c.SVIP,
		InitiatorIFace: c.InitiatorIFace,
		Types:          c.Types,
	}
	defaultTenantName := c.TenantName

	log.WithFields(log.Fields{
		"endpoint":          endpoint,
		"defaultSizeGiB":    defaultSizeGiB,
		"svip":              svip,
		"cfg":               cfg,
		"defaultTenantName": defaultTenantName,
	}).Debug("About to call NewFromParameters")

	// creaet a new sfapi.Client object for interacting with the SolidFire storage system
	client, _ := sfapi.NewFromParameters(endpoint, defaultSizeGiB, svip, cfg, defaultTenantName)
	req := sfapi.GetAccountByNameRequest{
		Name: c.TenantName,
	}

	// lookup the specified account; if not found, dynamically create it
	account, err := client.GetAccountByName(&req)
	if err != nil {
		req := sfapi.AddAccountRequest{
			Username: c.TenantName,
		}
		tenantID, err = client.AddAccount(&req)
		if err != nil {
			log.Fatal("Failed to initialize solidfire driver while creating tenant: ", err)
		}
	} else {
		tenantID = account.AccountID
	}

	iscsiInterface := "default"
	if c.InitiatorIFace != "" {
		iscsiInterface = c.InitiatorIFace
	}

	if c.Types != nil {
		client.VolumeTypes = c.Types
	}

	defaultVolSize := int64(1)
	if c.DefaultVolSz != 0 {
		defaultVolSize = c.DefaultVolSz
	}

	d.TenantID = tenantID
	d.Client = client
	d.DefaultVolSz = defaultVolSize
	d.InitiatorIFace = iscsiInterface
	log.WithFields(log.Fields{
		"TenantID":       tenantID,
		"DefaultVolSz":   defaultVolSize,
		"InitiatorIFace": iscsiInterface,
	}).Debug("Driver initialized with the following settings")

	validationErr := d.Validate()
	if validationErr != nil {
		return fmt.Errorf("Problem validating SolidfireSANStorageDriver error: %v", validationErr)
	}

	// log an informational message when this plugin starts
	// TODO how does solidfire do this?
	//EmsInitialized(d.Name(), d.api)

	d.Initialized = true
	log.Infof("Successfully initialized SolidFire Docker driver version %v", DriverVersion)
	return nil
}

// Validate the driver configuration and execution environment
func (d *SolidfireSANStorageDriver) Validate() error {
	log.Debugf("SolidfireSANStorageDriver#Validate()")

	// We want to verify we have everything we need to run the Docker driver
	if d.Config.TenantName == "" {
		log.Fatal("TenantName required in SolidFire Docker config")
	}
	if d.Config.EndPoint == "" {
		log.Fatal("EndPoint required in SolidFire Docker config")
	}
	if d.Config.DefaultVolSz == 0 {
		log.Fatal("DefaultVolSz required in SolidFire Docker config")
	}
	if d.Config.SVIP == "" {
		log.Fatal("SVIP required in SolidFire Docker config")
	}

	// Validate the environment
	isIscsiSupported := utils.IscsiSupported()
	if !isIscsiSupported {
		return fmt.Errorf("iSCSI support not detected")
	}

	return nil
}

// Create a SolidFire volume
func (d *SolidfireSANStorageDriver) Create(name string, opts map[string]string) error {
	log.Debugf("SolidfireSANStorageDriver#Create(%v)", name)

	var req sfapi.CreateVolumeRequest
	var qos sfapi.QoS
	var vsz int64
	var meta = map[string]string{"platform": "Docker-NDVP"}

	log.Debugf("GetVolumeByName: %s, %d", name, d.TenantID)
	log.Debugf("Options passed in to create: %+v", opts)
	v, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err == nil && v.VolumeID != 0 {
		log.Infof("Found existing Volume by name: %s", name)
		return nil
	}

	formatOpts(opts)
	log.Debugf("Options after conversion: %+v", opts)
	if opts["size"] != "" {
		s, _ := strconv.ParseInt(opts["size"], 10, 64)
		log.Info("Received size request in Create: ", s)
		vsz = int64(units.GiB) * s
	} else {
		// NOTE(jdg): We need to cleanup the conversions and such when we read
		// in from the config file, it's sort of ugly.  BUT, just remember that
		// when we pull the value from d.DefaultVolSz it's already been
		// multiplied
		vsz = d.DefaultVolSz
		log.Info("Creating with default size of: ", vsz)
	}

	if opts["qos"] != "" {
		iops := strings.Split(opts["qos"], ",")
		qos.MinIOPS, _ = strconv.ParseInt(iops[0], 10, 64)
		qos.MaxIOPS, _ = strconv.ParseInt(iops[1], 10, 64)
		qos.BurstIOPS, _ = strconv.ParseInt(iops[2], 10, 64)
		req.Qos = qos
		log.Infof("Received qos opts in Create: %+v", req.Qos)
	}

	if opts["type"] != "" {
		for _, t := range *d.Client.VolumeTypes {
			if strings.EqualFold(t.Type, opts["type"]) {
				req.Qos = t.QOS
				log.Infof("Received Type opts in Create and set QoS: %+v", req.Qos)
				break
			}
		}
	}

	req.TotalSize = vsz
	req.AccountID = d.TenantID
	req.Name = name
	req.Attributes = meta
	_, err = d.Client.CreateVolume(&req)
	if err != nil {
		return err
	}
	return nil
}

// Create a volume clone
func (d *SolidfireSANStorageDriver) CreateClone(name, source, snapshot, newSnapshotPrefix string) error {
	log.Debugf("SolidfireSANStorageDriver#CreateClone(%v, %v, %v, %v)", name, source, snapshot, newSnapshotPrefix)

	var req sfapi.CloneVolumeRequest

	// Check to see if the clone already exists
	v, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err == nil && v.VolumeID != 0 {
		// The clone already exists; skip and call it a success
		return nil
	}

	// If a snapshot was specified, use that
	if snapshot != "" {
		s, err := d.Client.GetSnapshot(0, v.VolumeID, snapshot)
		if err != nil || s.SnapshotID == 0 {
			return fmt.Errorf("Failed to find snapshot specified: error: %v", err)
		}
		req.SnapshotID = s.SnapshotID
	}

	// Get the volume ID for the source volume
	v, err = d.Client.GetVolumeByName(source, d.TenantID)
	if err != nil || v.VolumeID == 0 {
		return fmt.Errorf("Failed to find source volume: error: %v", err)
	}

	// Create the clone of the source volume with the name specified
	req.VolumeID = v.VolumeID
	req.Name = name
	_, err = d.Client.CloneVolume(&req)
	if err != nil {
		return fmt.Errorf("Failed to create clone: error: %v", err)
	}
	return nil
}

// Destroy the requested docker volume
func (d *SolidfireSANStorageDriver) Destroy(name string) error {
	log.Debugf("SolidfireSANStorageDriver#Destroy(%v)", name)

	v, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve volume named %s during Remove operation;  error: %v", name, err)
	}
	d.Client.DetachVolume(v)
	err = d.Client.DeleteVolume(v.VolumeID)
	if err != nil {
		// FIXME(jdg): Check if it's a "DNE" error in that case we're golden
		log.Error("Error encountered during delete: ", err)
	}

	// perform rediscovery to remove the deleted LUN
	//utils.MultipathFlush() // flush unused paths
	//utils.IscsiRescan()

	return nil
}

// Attach the lun
func (d *SolidfireSANStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("SolidfireSANStorageDriver#Attach(%v, %v, %v)", name, mountpoint, opts)

	v, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve volume by name in mount operation;  name: %v error: %v", name, err)
	}
	path, device, err := d.Client.AttachVolume(&v, d.InitiatorIFace)
	if path == "" || device == "" && err == nil {
		return fmt.Errorf("Problem attaching docker volume but err is nil;  path: %v device: %v", path, device)
	}
	if err != nil {
		return fmt.Errorf("Failed to perform iscsi attach;  volume: %s error: %v", name, err)
	}
	log.Debugf("Attached volume at (path, devfile): %s, %s", path, device)
	if utils.GetFSType(device) == "" {
		//TODO(jdg): Enable selection of *other* fs types
		err := utils.FormatVolume(device, "ext4")
		if err != nil {
			return fmt.Errorf("Failed to format device: %v error: %v", device, err)
		}
	}
	if mountErr := utils.Mount(device, mountpoint); mountErr != nil {
		return fmt.Errorf("Problem mounting docker volume: %v device: %v mountpoint: %v error: %v", name, device, mountpoint, mountErr)
	}

	return nil
}

// Detach the volume
func (d *SolidfireSANStorageDriver) Detach(name, mountpoint string) error {
	log.Debugf("SolidfireSANStorageDriver#Detach(%v, %v)", name, mountpoint)

	umountErr := utils.Umount(mountpoint)
	if umountErr != nil {
		return fmt.Errorf("Problem unmounting docker volume: %v mountpoint: %v error: %v", name, mountpoint, umountErr)
	}

	v, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err != nil {
		return fmt.Errorf("Problem looking up volume name: %v TenantID: %v error: %v", name, d.TenantID, err)
	}
	d.Client.DetachVolume(v)

	return nil
}

// DefaultStoragePrefix is the driver specific prefix for created storage, can be overridden in the config file
func (d *SolidfireSANStorageDriver) DefaultStoragePrefix() string {
	return "netappdvp-"
}

// DefaultSnapshotPrefix is the driver specific prefix for created snapshots, can be overridden in the config file
func (d *SolidfireSANStorageDriver) DefaultSnapshotPrefix() string {
	return "netappdvp-"
}

// Return the list of snapshots associated with the named volume
func (d *SolidfireSANStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {
	log.Debugf("SolidfireSANStorageDriver#SnapshotList(%v)", name)

	v, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve volume by name in snapshot list operation; name: %v error: %v", name, err)
	}

	var req sfapi.ListSnapshotsRequest
	req.VolumeID = v.VolumeID

	s, err := d.Client.ListSnapshots(&req)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve snapshots for volume; name: %v error: %v", name, err)
	}

	log.Debugf("Returned %v snapshots", len(s))
	var snapshots []CommonSnapshot

	for _, snap := range s {
		log.Debugf("Snapshot name: %v, date: %v", snap.Name, snap.CreateTime)
		snapshots = append(snapshots, CommonSnapshot{snap.Name, snap.CreateTime})
	}

	return snapshots, nil

// List volumes
func (d *SolidfireSANStorageDriver) ListVolumes() (vols []string, err error) {
	log.Debugf("SolidfireSANStorageDriver#ListVolumes")

	var listReq sfapi.ListVolumesForAccountRequest

	listReq.AccountID = d.TenantID
	volList, err := d.Client.ListVolumesForAccount(&listReq)
	if err != nil {
		return nil, fmt.Errorf("Problem looking up volumes for TenantID: %v", d.TenantID)
	}

	for _, v := range volList {
		if v.Status == "active" && v.AccountID == d.TenantID {
			vols = append(vols, v.Name)
		}
	}
	return vols, nil
}

// get a volume
func (d *SolidfireSANStorageDriver) VolGet(name string) (volID int64, err error) {
	log.Debugf("SolidfireSANStorageDriver#VolGet(%v)", name)
	vol, err := d.Client.GetVolumeByName(name, d.TenantID)
	if err != nil {
		return 0, fmt.Errorf("Problem looking up volumes for TenantID: >%v", d.TenantID)
	}
	return vol.VolumeID, nil
}
