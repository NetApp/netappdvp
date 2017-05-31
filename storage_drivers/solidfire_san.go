// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/apis/sfapi"
	"github.com/netapp/netappdvp/utils"
)

// SolidfireSANStorageDriverName is the constant name for this Solidfire SAN storage driver
const SolidfireSANStorageDriverName = "solidfire-san"

func init() {
	san := &SolidfireSANStorageDriver{}
	san.Initialized = false
	Drivers[san.Name()] = san
	log.Debugf("registered driver '%s'", san.Name())
}

// SolidfireSANStorageDriver is for iSCSI storage provisioning
type SolidfireSANStorageDriver struct {
	Initialized      bool
	Config           SolidfireStorageDriverConfig
	Client           *sfapi.Client
	TenantID         int64
	AccessGroups     []int64
	LegacyNamePrefix string
	InitiatorIFace   string
}

// Name is for returning the name of this driver
func (d SolidfireSANStorageDriver) Name() string {
	return SolidfireSANStorageDriverName
}

// Initialize from the provided config
func (d *SolidfireSANStorageDriver) Initialize(configJSON string, commonConfig *CommonStorageDriverConfig) error {
	log.Debug("SolidfireSANStorageDriver#Initialize(...)")
	c := &SolidfireStorageDriverConfig{}
	c.CommonStorageDriverConfig = commonConfig

	// decode supplied configJSON string into SolidfireStorageDriverConfig object
	err := json.Unmarshal([]byte(configJSON), &c)
	if err != nil {
		log.Errorf("Cannot decode json configuration error: %v", err)
		return errors.New("json decode error reading config")
	}

	log.WithFields(log.Fields{
		"Version":           c.Version,
		"StorageDriverName": c.StorageDriverName,
		"DisableDelete":     c.DisableDelete,
	}).Debugf("Reparsed into solidfireConfig")

	// SF prefix is always empty
	prefix := ""
	c.StoragePrefix = &prefix

	log.Debugf("Decoded to %+v", c)
	d.Config = *c

	var tenantID int64

	// create a new sfapi.Config object from the read in json config file
	endpoint := c.EndPoint
	svip := c.SVIP
	cfg := sfapi.Config{
		TenantName:       c.TenantName,
		EndPoint:         c.EndPoint,
		SVIP:             c.SVIP,
		InitiatorIFace:   c.InitiatorIFace,
		Types:            c.Types,
		LegacyNamePrefix: c.LegacyNamePrefix,
		AccessGroups:     c.AccessGroups,
	}
	defaultTenantName := c.TenantName

	log.WithFields(log.Fields{
		"endpoint":          endpoint,
		"svip":              svip,
		"cfg":               cfg,
		"defaultTenantName": defaultTenantName,
	}).Debug("About to call NewFromParameters")

	// create a new sfapi.Client object for interacting with the SolidFire storage system
	client, _ := sfapi.NewFromParameters(endpoint, svip, cfg, defaultTenantName)
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

	if c.AccessGroups != nil {
		client.AccessGroups = c.AccessGroups
	}

	if c.DefaultVolSz != 0 {
		log.Warn("Configuration file setting DefaultVolSz is deprecated " +
			"and will be ignored.  Use defaults:{size} instead.")
	}

	d.TenantID = tenantID
	d.Client = client
	d.InitiatorIFace = iscsiInterface
	d.LegacyNamePrefix = legacyNamePrefix
	log.WithFields(log.Fields{
		"TenantID":       tenantID,
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
func (d *SolidfireSANStorageDriver) Create(name string, sizeBytes uint64, opts map[string]string) error {
	log.Debugf("SolidfireSANStorageDriver#Create(%s)", name)

	var req sfapi.CreateVolumeRequest
	var qos sfapi.QoS
	var meta = map[string]string{"platform": "Docker-NDVP",
		"ndvp-version": DriverVersion + " [" + ExtendedDriverVersion + "]",
		"docker-name":  name}

	v, err := d.GetVolume(name)
	if err == nil && v.VolumeID != 0 {
		log.Warningf("found existing Volume by name: %s", name)
		return errors.New("volume with requested name already exists")
	}

	qos_opt := utils.GetV(opts, "qos", "")
	if qos_opt != "" {
		iops := strings.Split(qos_opt, ",")
		qos.MinIOPS, _ = strconv.ParseInt(iops[0], 10, 64)
		qos.MaxIOPS, _ = strconv.ParseInt(iops[1], 10, 64)
		qos.BurstIOPS, _ = strconv.ParseInt(iops[2], 10, 64)
		req.Qos = qos
		log.Infof("received qos opts in Create: %+v", req.Qos)
	}

	type_opt := utils.GetV(opts, "type", "")
	if type_opt != "" {
		for _, t := range *d.Client.VolumeTypes {
			if strings.EqualFold(t.Type, type_opt) {
				req.Qos = t.QOS
				log.Infof("received Type opts in Create and set QoS: %+v", req.Qos)
				break
			}
		}
	}

	req.TotalSize = int64(sizeBytes)
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
func (d *SolidfireSANStorageDriver) CreateClone(name, source, snapshot string, opts map[string]string) error {
	log.Debugf("SolidfireSANStorageDriver#CreateClone(%s, %s, %s)", name, source, snapshot)

	var req sfapi.CloneVolumeRequest
	var meta = map[string]string{
		"platform":     "Docker-NDVP",
		"ndvp-version": DriverVersion + " [" + ExtendedDriverVersion + "]",
		"docker-name":  name,
	}

	// Check to see if the clone already exists
	v, err := d.GetVolume(name)
	if err == nil && v.VolumeID != 0 {
		log.Warningf("found existing Volume by name: %s", name)
		return errors.New("volume with requested name already exists")
	}

	// Get the volume ID for the source volume
	v, err = d.GetVolume(source)
	if err != nil || v.VolumeID == 0 {
		log.Errorf("unable to locate requested source volume: %+v", err)
		return errors.New("error performing clone operation, source volume not found")
	}

	// If a snapshot was specified, use that
	if snapshot != "" {
		s, err := d.Client.GetSnapshot(0, v.VolumeID, snapshot)
		if err != nil || s.SnapshotID == 0 {
			log.Errorf("unable to locate requested source snapshot: %+v", err)
			return errors.New("error performing clone operation, source snapshot not found")
		}
		req.SnapshotID = s.SnapshotID
	}

	// Create the clone of the source volume with the name specified
	req.VolumeID = v.VolumeID
	req.Name = MakeSolidFireName(name)
	req.Attributes = meta
	_, err = d.Client.CloneVolume(&req)
	if err != nil {
		log.Errorf("failed to create clone: %+v", err)
		return errors.New("error performing clone operation")
	}
	return nil
}

// Destroy the requested docker volume
func (d *SolidfireSANStorageDriver) Destroy(name string) error {
	log.Debugf("SolidfireSANStorageDriver#Destroy(%s)", name)

	v, err := d.GetVolume(name)
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
	v, err := d.GetVolume(name)
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

	v, err := d.GetVolume(name)
	if err != nil {
		log.Errorf("unable to locate volume: %+v", err)
		return errors.New("volume not found")
	}
	d.Client.DetachVolume(v)

	return nil
}

// Return the list of volumes according to backend device
func (d *SolidfireSANStorageDriver) List() (vols []string, err error) {
	var req sfapi.ListVolumesForAccountRequest
	req.AccountID = d.TenantID
	volumes, err := d.Client.ListVolumesForAccount(&req)
	for _, v := range volumes {
		if v.Status != "deleted" {
			attrs, _ := v.Attributes.(map[string]interface{})
			dName := strings.Replace(v.Name, d.LegacyNamePrefix, "", -1)
			if str, ok := attrs["docker-name"].(string); ok {
				dName = strings.Replace(str, d.LegacyNamePrefix, "", -1)
			}
			vols = append(vols, dName)
		}
	}
	return vols, err
}

// Return the list of snapshots associated with the named volume
func (d *SolidfireSANStorageDriver) SnapshotList(name string) ([]CommonSnapshot, error) {
	log.Debugf("SolidfireSANStorageDriver#SnapshotList(%s)", name)
	v, err := d.GetVolume(name)
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

// Test for the existence of a volume
func (d *SolidfireSANStorageDriver) Get(name string) error {
	_, err := d.GetVolume(name)
	return err
}

func (d *SolidfireSANStorageDriver) GetVolume(name string) (sfapi.Volume, error) {
	var vols []sfapi.Volume
	var req sfapi.ListVolumesForAccountRequest

	// I know, I know... just use V8 of the API and let the Cluster filter on
	// things like Name; trouble is we completely screwed up Name usage so we
	// can't trust it.  We now have a few possibilities including Name,
	// Name-With-Prefix and Attributes.  It could be any of the 3.  At some
	// point let's fix that and just use something efficient like Name and be
	// done with it. Otherwise, we just get all for the account and iterate
	// which isn't terrible.
	req.AccountID = d.TenantID
	volumes, err := d.Client.ListVolumesForAccount(&req)
	if err != nil {
		log.Errorf("error encountered requesting volumes in SolidFire:getVolume: %+v", err)
		return sfapi.Volume{}, errors.New("device reported API error")
	}

	legacyName := MakeSolidFireName(d.LegacyNamePrefix + name)
	baseSFName := MakeSolidFireName(name)

	for _, v := range volumes {
		attrs, _ := v.Attributes.(map[string]interface{})
		// We prefer attributes, so check that first, then pick up legacy
		// volumes using Volume Name
		if attrs["docker-name"] == name && v.Status == "active" {
			log.Debugf("found volume by attributes: %+v", v)
			vols = append(vols, v)
		} else if (v.Name == legacyName || v.Name == baseSFName) && v.Status == "active" {
			log.Warningf("found volume by name using deprecated Volume.Name mapping: %+v", v)
			vols = append(vols, v)
		}
	}
	if len(vols) == 0 {
		return sfapi.Volume{}, errors.New("volume not found")
	}
	return vols[0], nil
}
