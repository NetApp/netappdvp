// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/units"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/netapp/netappdvp/apis/sfapi"
	"github.com/netapp/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

// SolidfireSANStorageDriverName is the constant name for this Solidfire SAN storage driver
const SolidfireSANStorageDriverName = "solidfire-san"

func init() {
	san := &SolidfireSANStorageDriver{}
	san.Initialized = false
	Drivers[san.Name()] = san
	log.Debugf("registered driver '%s'", san.Name())
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
	Initialized      bool
	Config           SolidfireStorageDriverConfig
	Client           *sfapi.Client
	TenantID         int64
	DefaultVolSz     int64
	VagID            int64
	LegacyNamePrefix string
	InitiatorIFace   string
}

// Name is for returning the name of this driver
func (d SolidfireSANStorageDriver) Name() string {
	return SolidfireSANStorageDriverName
}

// Initialize from the provided config
func (d *SolidfireSANStorageDriver) Initialize(configJSON string) error {
	log.Debug("SolidfireSANStorageDriver#Initialize(...)")
	c := &SolidfireStorageDriverConfig{}

	// decode supplied configJSON string into SolidfireStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &c)
	if err != nil {
		log.Errorf("Cannot decode json configuration error: %v", err)
		return errors.New("json decode error reading config")
	}

	log.WithFields(log.Fields{
		"Version":           c.Version,
		"StorageDriverName": c.StorageDriverName,
		"Debug":             c.Debug,
		"DisableDelete":     c.DisableDelete,
		"StoragePrefixRaw":  string(c.StoragePrefixRaw),
		"SnapshotPrefixRaw": string(c.SnapshotPrefixRaw),
	}).Debugf("Reparsed into solidfireConfig")

	log.Debugf("Decoded to %+v", c)
	d.Config = *c

	var tenantID int64

	// create a new sfapi.Config object from the read in json config file
	endpoint := c.EndPoint
	defaultSizeGiB := c.DefaultVolSz * int64(units.GiB)
	svip := c.SVIP
	cfg := sfapi.Config{
		TenantName:       c.TenantName,
		EndPoint:         c.EndPoint,
		DefaultVolSz:     defaultSizeGiB,
		SVIP:             c.SVIP,
		InitiatorIFace:   c.InitiatorIFace,
		Types:            c.Types,
		LegacyNamePrefix: c.LegacyNamePrefix,
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
			log.Fatalf("failed to initialize solidfire driver while creating tenant: %+v", err)
		}
	} else {
		tenantID = account.AccountID
	}

	legacyNamePrefix := "netappdvp-"
	if c.LegacyNamePrefix != "" {
		legacyNamePrefix = c.LegacyNamePrefix
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
		defaultVolSize = defaultSizeGiB
	}

	d.TenantID = tenantID
	d.Client = client
	d.DefaultVolSz = defaultVolSize
	d.InitiatorIFace = iscsiInterface
	d.LegacyNamePrefix = legacyNamePrefix
	log.WithFields(log.Fields{
		"TenantID":       tenantID,
		"DefaultVolSz":   defaultVolSize,
		"InitiatorIFace": iscsiInterface,
	}).Debug("Driver initialized with the following settings")

	validationErr := d.Validate()
	if validationErr != nil {
		log.Errorf("problem validating SolidfireSANStorageDriver error: %+v", validationErr)
		return errors.New("error encountered validating SolidFire driver on init")
	}

	// log an informational message when this plugin starts
	// TODO how does solidfire do this?
	//EmsInitialized(d.Name(), d.api)

	d.Initialized = true
	log.Infof("successfully initialized SolidFire Docker driver version %s [%s]", DriverVersion, ExtendedDriverVersion)
	return nil
}

// Validate the driver configuration and execution environment
func (d *SolidfireSANStorageDriver) Validate() error {
	log.Debug("SolidfireSANStorageDriver#Validate()")

	// We want to verify we have everything we need to run the Docker driver
	if d.Config.TenantName == "" {
		log.Fatal("missing required TenantName in config")
	}
	if d.Config.EndPoint == "" {
		log.Fatal("missing required EndPoint in config")
	}
	if d.Config.SVIP == "" {
		log.Fatal("missing required SVIP in config")
	}

	// Validate the environment
	isIscsiSupported := utils.IscsiSupported()
	if !isIscsiSupported {
		log.Errorf("host doesn't appear to support iSCSI")
		return errors.New("no iSCSI support on this host")
	}

	return nil
}

// Make SolidFire name
func MakeSolidFireName(name string) string {
	return strings.Replace(name, "_", "-", -1)
}

// Create a SolidFire volume
func (d *SolidfireSANStorageDriver) Create(name string, opts map[string]string) error {
	log.Debugf("SolidfireSANStorageDriver#Create(%s)", name)

	var req sfapi.CreateVolumeRequest
	var qos sfapi.QoS
	var vsz int64
	var meta = map[string]string{"platform": "Docker-NDVP",
		"ndvp-version": DriverVersion + " [" + ExtendedDriverVersion + "]",
		"docker-name":  name}

	v, err := d.getVolume(name)
	if err == nil && v.VolumeID != 0 {
		log.Infof("found existing Volume by name: %s", name)
		return nil
	}

	formatOpts(opts)
	log.Debugf("options after conversion: %+v", opts)
	if opts["size"] != "" {
		s, _ := strconv.ParseInt(opts["size"], 10, 64)
		log.Infof("received size request in Create: %s ", s)
		vsz = int64(units.GiB) * s
	} else {
		// NOTE(jdg): We need to cleanup the conversions and such when we read
		// in from the config file, it's sort of ugly.  BUT, just remember that
		// when we pull the value from d.DefaultVolSz it's already been
		// multiplied
		vsz = d.DefaultVolSz
		log.Infof("creating with default size of: %s", vsz)
	}

	if opts["qos"] != "" {
		iops := strings.Split(opts["qos"], ",")
		qos.MinIOPS, _ = strconv.ParseInt(iops[0], 10, 64)
		qos.MaxIOPS, _ = strconv.ParseInt(iops[1], 10, 64)
		qos.BurstIOPS, _ = strconv.ParseInt(iops[2], 10, 64)
		req.Qos = qos
		log.Infof("received qos opts in Create: %+v", req.Qos)
	}

	if opts["type"] != "" {
		for _, t := range *d.Client.VolumeTypes {
			if strings.EqualFold(t.Type, opts["type"]) {
				req.Qos = t.QOS
				log.Infof("received Type opts in Create and set QoS: %+v", req.Qos)
				break
			}
		}
	}

	req.TotalSize = vsz
	req.AccountID = d.TenantID
	req.Name = MakeSolidFireName(name)
	req.Attributes = meta
	_, err = d.Client.CreateVolume(&req)
	if err != nil {
		return err
	}
	return nil
}

// Create a volume clone
func (d *SolidfireSANStorageDriver) CreateClone(name, source, snapshot, newSnapshotPrefix string) error {
	log.Debugf("SolidfireSANStorageDriver#CreateClone(%s, %s, %s, %s)", name, source, snapshot, newSnapshotPrefix)

	var req sfapi.CloneVolumeRequest
	var meta = map[string]string{"platform": "Docker-NDVP",
		"ndvp-version": DriverVersion + " [" + ExtendedDriverVersion + "]",
		"docker-name":  name}

	// Check to see if the clone already exists
	v, err := d.getVolume(name)
	if err == nil && v.VolumeID != 0 {
		// The clone already exists; skip and call it a success
		return nil
	}

	// If a snapshot was specified, use that
	if snapshot != "" {
		s, err := d.Client.GetSnapshot(0, v.VolumeID, snapshot)
		if err != nil || s.SnapshotID == 0 {
			log.Errorf("unable to locate requested snapshot: %+v")
			return errors.New("no iSCSI support on this host")
		}
		req.SnapshotID = s.SnapshotID
	}

	// Get the volume ID for the source volume
	v, err = d.getVolume(source)
	if err != nil || v.VolumeID == 0 {
		log.Errorf("unable to locate requested source volume: %+v", err)
		return errors.New("volume not found")
	}

	// Create the clone of the source volume with the name specified
	req.VolumeID = v.VolumeID
	req.Name = MakeSolidFireName(name)
	req.Attributes = meta
	_, err = d.Client.CloneVolume(&req)
	if err != nil {
		log.Errorf("failed to create clone: error: %+v", err)
		return errors.New("error performaing clone operation")
	}
	return nil
}

// Destroy the requested docker volume
func (d *SolidfireSANStorageDriver) Destroy(name string) error {
	log.Debugf("SolidfireSANStorageDriver#Destroy(%s)", name)

	v, err := d.getVolume(name)
	if err != nil {
		log.Errorf("unable to locate volume for delete operation: %+v", err)
		return errors.New("volume not found")
	}
	d.Client.DetachVolume(v)
	err = d.Client.DeleteVolume(v.VolumeID)
	if err != nil {
		// FIXME(jdg): Check if it's a "DNE" error in that case we're golden
		log.Errorf("error during delete operation: %+v", err)
	}

	// perform rediscovery to remove the deleted LUN
	//utils.MultipathFlush() // flush unused paths
	//utils.IscsiRescan()

	return nil
}

// Attach the lun
func (d *SolidfireSANStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("SolidfireSANStorageDriver#Attach(%s, %s, %+v)", name, mountpoint, opts)
	v, err := d.getVolume(name)
	if err != nil {
		log.Errorf("unable to locate volume for mount operation: %+v", err)
		return errors.New("volume not found")
	}
	path, device, err := d.Client.AttachVolume(&v, d.InitiatorIFace)
	if path == "" || device == "" && err == nil {
		log.Errorf("path not found on attach: (path: %s, device: %s)", path, device)
		return errors.New("path not found")
	}
	if err != nil {
		log.Errorf("error on iSCSI attach: %+v", err)
		return errors.New("iSCSI attach error")
	}
	log.Debugf("Attached volume at (path, devfile): %s, %s", path, device)
	if utils.GetFSType(device) == "" {
		//TODO(jdg): Enable selection of *other* fs types
		err := utils.FormatVolume(device, "ext4")
		if err != nil {
			log.Errorf("error on formatting volume: %+v", err)
			return errors.New("format (mkfs) error")
		}
	}
	if mountErr := utils.Mount(device, mountpoint); mountErr != nil {
		log.Errorf("unable to mount device: (device: %s, mountpoint: %s, error: %+v", device, mountpoint, err)
		return errors.New("unable to mount device")
	}

	return nil
}

// Detach the volume
func (d *SolidfireSANStorageDriver) Detach(name, mountpoint string) error {
	log.Debugf("SolidfireSANStorageDriver#Detach(%s, %s)", name, mountpoint)
	umountErr := utils.Umount(mountpoint)
	if umountErr != nil {
		log.Errorf("unable to unmount device: (name: %s, mountpoint: %s, error: %+v", name, mountpoint, umountErr)
		return errors.New("unable to unmount device")
	}

	v, err := d.getVolume(name)
	if err != nil {
		log.Errorf("unable to locate volume: %+v", err)
		return errors.New("volume not found")
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

// Return the list of volumes according to backend device
func (d *SolidfireSANStorageDriver) VolumeList(vDir string) ([]*volume.Volume, error) {
	log.Info("List volumes from SolidFire backend")
	var req sfapi.ListVolumesForAccountRequest
	var vols []*volume.Volume
	req.AccountID = d.TenantID
	volumes, err := d.Client.ListVolumesForAccount(&req)
	for _, v := range volumes {
		if v.Status != "deleted" {
			attrs, _ := v.Attributes.(map[string]interface{})
			dName := strings.Replace(v.Name, d.LegacyNamePrefix, "", -1)
			if str, ok := attrs["docker-name"].(string); ok {
				dName = strings.Replace(str, d.LegacyNamePrefix, "", -1)
			}
			vols = append(vols, &volume.Volume{Name: dName, Mountpoint: filepath.Join(vDir, v.Name)})
		}
	}
	return vols, err
}

// Return the list of snapshots associated with the named volume
func (d *SolidfireSANStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {
	log.Debugf("SolidfireSANStorageDriver#SnapshotList(%s)", name)
	v, err := d.getVolume(name)
	if err != nil {
		log.Errorf("unable to locate parent volume in snapshotlist: %+v", err)
		return nil, errors.New("volume not found")
	}

	var req sfapi.ListSnapshotsRequest
	req.VolumeID = v.VolumeID

	s, err := d.Client.ListSnapshots(&req)
	if err != nil {
		log.Errorf("unable to locate snapshot: %+v", err)
		return nil, errors.New("snapshot not found")
	}

	log.Debugf("returned %d snapshots", len(s))
	var snapshots []CommonSnapshot

	for _, snap := range s {
		log.Debugf("snapshot name: %s, date: %s", snap.Name, snap.CreateTime)
		snapshots = append(snapshots, CommonSnapshot{snap.Name, snap.CreateTime})
	}

	return snapshots, nil
}

// Get volume using a couple of different methods to support changes after
// upgrades
func (d *SolidfireSANStorageDriver) getVolume(name string) (v sfapi.Volume, err error) {
	// By default we now use the attributes which is the raw docker name so we
	// can ignore the prefix/translation nonsense for the first go around
	v, err = d.Client.GetVolumeByDockerName(name, d.TenantID)

	// Ok, the volume may not exist which we expect in a number of cases, but
	// now we have the annoying challenge of determining, was this in fact
	// because it DNE, or is it a problem with legacy volume-name munging?
	if (err != nil) && (d.LegacyNamePrefix != "") {
		// We'll allow a user to specify the old naming convention in their SF
		// config and use that to try and find by name using the old
		// translation method, note that we default this prefix to 'netappdvp-'
		legacyName := d.LegacyNamePrefix + name
		legacyName = MakeSolidFireName(legacyName)
		log.Debugf("Attempting failed search using legacy-name: %s", legacyName)
		v, err = d.Client.GetVolumeByName(legacyName, d.TenantID)
	}
	if err != nil {
		log.Errorf("failed to retrieve volume: %s (%+v)", name, err)
		return v, errors.New("volume not found")
	}
	return v, nil
}
