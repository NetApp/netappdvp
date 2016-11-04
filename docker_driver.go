// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/ebalduf/netappdvp/storage_drivers"
	"github.com/ebalduf/netappdvp/utils"

	log "github.com/Sirupsen/logrus"
)

type ndvpDriver struct {
	m      *sync.Mutex
	root   string
	config storage_drivers.CommonStorageDriverConfig
	sd     storage_drivers.StorageDriver
}

func (d *ndvpDriver) volumePrefix() string {
	defaultPrefix := d.sd.DefaultStoragePrefix()
	prefixToUse := defaultPrefix
	storagePrefixRaw := d.config.StoragePrefixRaw // this is a raw version of the json value, we will get quotes in it
	if len(storagePrefixRaw) >= 2 {
		s := string(storagePrefixRaw)
		if s == "\"\"" || s == "" {
			prefixToUse = ""
			//log.Debugf("storagePrefix is specified as \"\", using no prefix")
		} else {
			// trim quotes from start and end of string
			prefixToUse = s[1 : len(s)-1]
			//log.Debugf("storagePrefix is specified, using prefix: %v", prefixToUse)
		}
	} else {
		prefixToUse = defaultPrefix
		//log.Debugf("storagePrefix is unspecified, using default prefix: %v", prefixToUse)
	}

	return prefixToUse
}

func (d *ndvpDriver) volumeName(name string) string {
	prefixToUse := d.volumePrefix()
	if strings.HasPrefix(name, prefixToUse) {
		return name
	}
	return prefixToUse + name
}

func (d *ndvpDriver) snapshotPrefix() string {
	defaultPrefix := d.sd.DefaultSnapshotPrefix()
	prefixToUse := defaultPrefix
	snapshotPrefixRaw := d.config.SnapshotPrefixRaw // this is a raw version of the json value, we will get quotes in it
	if len(snapshotPrefixRaw) >= 2 {
		s := string(snapshotPrefixRaw)
		if s == "\"\"" || s == "" {
			prefixToUse = ""
			//log.Debugf("snapshotPrefix is specified as \"\", using no prefix")
		} else {
			// trim quotes from start and end of string
			prefixToUse = s[1 : len(s)-1]
			//log.Debugf("snapshotPrefix is specified, using prefix: %v", prefixToUse)
		}
	} else {
		prefixToUse = defaultPrefix
		//log.Debugf("snapshotPrefix is unspecified, using default prefix: %v", prefixToUse)
	}

	return prefixToUse
}

func (d *ndvpDriver) mountpoint(name string) string {
	return filepath.Join(d.root, name)
}

func newNetAppDockerVolumePlugin(root string, config storage_drivers.CommonStorageDriverConfig) (*ndvpDriver, error) {
	// if root (volumeDir) doesn't exist, make it
	dir, err := os.Lstat(root)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(root, 0755); err != nil {
			return nil, err
		}
	}
	// if root (volumeDir) isn't a directory, error
	if dir != nil && !dir.IsDir() {
		return nil, fmt.Errorf("Volume directory '%v' exists and it's not a directory", root)
	}

	d := &ndvpDriver{
		root:   root,
		config: config,
		m:      &sync.Mutex{},
		sd:     storage_drivers.Drivers[config.StorageDriverName],
	}
	return d, nil
}

// Create is part of the core Docker API and is called to create a docker volume
func (d ndvpDriver) Create(r volume.Request) volume.Response {
	log.Debugf("Docker volume Create %s ", r.Name)
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Create(%v)", r)

	opts := r.Options
	target := d.volumeName(r.Name)
	m := d.mountpoint(target)

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(m, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v already exists and it's not a directory", m)}
	}

	var createErr error

	// If 'from' is specified, create a snapshot and a clone rather than a new empty volume
	if from, ok := opts["from"]; ok {
		source := d.volumeName(from)

		// If 'fromSnapshot' is specified, we use the existing snapshot instead
		snapshot := opts["fromSnapshot"]
		createErr = d.sd.CreateClone(target, source, snapshot, d.snapshotPrefix())
	} else {
		createErr = d.sd.Create(target, opts)
	}

	if createErr != nil {
		os.Remove(m)
		return volume.Response{Err: fmt.Sprintf("Error creating storage: %v", createErr)}
	}

	return volume.Response{}
}

// Remove is part of the core Docker API and is called to remove a docker volume
func (d ndvpDriver) Remove(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Remove(%v)", r)

	target := d.volumeName(r.Name)

	// allow user to completely disable volume deletion
	if d.config.DisableDelete {
		log.Infof("Skipping removal of %s because of user preference to disable volume deletion", target)
		return volume.Response{}
	}

        // NOTE:  I choose to leave the mount points around instead of the
        // problem of tryin to clean them and then error exiting before we delete

	log.Debugf("Removing docker volume %s", target)

	// use the StorageDriver to destroy the storage objects
	destroyErr := d.sd.Destroy(target)
	if destroyErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem removing docker volume: %v error: %v", target, destroyErr)}
	}

	return volume.Response{}
}

func (d ndvpDriver) getPath(r volume.Request) (*volume.Volume, error) {
	target := d.volumeName(r.Name)

	volume := &volume.Volume{
		Name:       r.Name,
		Mountpoint: d.mountpoint(target)}
	return volume, nil
}

// Path is part of the core Docker API and is called to return the filesystem path to a docker volume
func (d ndvpDriver) Path(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Path(%v)", r)

	v, err := d.getPath(r)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{
		Mountpoint: v.Mountpoint,
	}
}

// Get is part of the core Docker API and is called to return details about a docker volume
func (d ndvpDriver) Get(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Get(%v)", r)

	// Gather the target volume name as the storage sees it
	target := d.volumeName(r.Name)

	v, err := d.getPath(r)

        // Check the array, is the volume still there, or just the path.
	_, err := d.sd.VolGet(d.volumeName(r.Name))
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	v, err := d.getPath(r)

	// Ask the storage driver for the list of snapshots associated with the volume
	snaps, err := d.sd.SnapshotList(target)

	// If we don't get any snapshots, that's fine. We'll return an empty list.
	status := map[string]interface{}{
		"Snapshots": snaps,
	}

	v2 := &volume.Volume{
		Name:       v.Name,
		Mountpoint: v.Mountpoint,
		Status:     status, // introduced in Docker 1.12, earlier versions ignore
	}

	return volume.Response{
		Volume: v2,
	}
}

type nfsMount struct {
	IP         string
	Target     string
	MountPoint string
	MachineID  string
}

// Mount is part of the core Docker API and is called to mount a docker volume
func (d ndvpDriver) Mount(r volume.MountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Mount(%v)", r)

	target := d.volumeName(r.Name)

	m := d.mountpoint(target)
	log.Debugf("Docker volume Mount %s on %s", target, m)

	fi, err := os.Lstat(m)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(m, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v already exists and it's not a directory", m)}
	}

	// check if already mounted before we do anything...
	dfOuput, dfOuputErr := utils.GetDFOutput()
	if dfOuputErr != nil {
		return volume.Response{Err: fmt.Sprintf("Error checking if %v is already mounted: %v", m, dfOuputErr)}
	}
	for _, e := range dfOuput {
		if e.Target == m {
			log.Debugf("%v already mounted, returning existing mount", m)
			return volume.Response{Mountpoint: m}
		}
	}

	// use the StorageDriver to attach the storage objects, place any extra options in this map
	attachOptions := make(map[string]string)

	attachErr := d.sd.Attach(target, m, attachOptions)
	if attachErr != nil {
		log.Error(attachErr)
		return volume.Response{Err: fmt.Sprintf("Problem attaching docker volume: %v mountpoint: %v error: %v", target, m, attachErr)}
	}

	return volume.Response{Mountpoint: m}
}

// Unmount is part of the core Docker API and is called to unmount a docker volume
func (d ndvpDriver) Unmount(r volume.UnmountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Unmount(%v)", r)

	target := d.volumeName(r.Name)

	m := d.mountpoint(target)
	log.Debugf("Unmounting docker volume %s", target)

	// use the StorageDriver to unmount the storage objects
	detachErr := d.sd.Detach(target, m)
	if detachErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem unmounting docker volume: %v error: %v", target, detachErr)}
	}
	// Remove the mountpoint
	log.Debugf("Removing mountpoint %s", m)
	if err := os.Remove(m); err != nil {
		log.Error("Failed to remove Mount directory: %v", err)
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{}
}

// List is part of the core Docker API and is called to list all known docker volumes for this plugin
func (d ndvpDriver) List(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("List(%v)", r)

	volumeDir := d.root
	volumePrefix := d.volumePrefix()
	var volumes []*volume.Volume
	// get list of volumes from array
	vols, err := d.sd.ListVolumes()
	if err != nil {
		log.Error("Failed to List volumes: %v", err)
		return volume.Response{Err: err.Error()}
	}

	// the list from the array doesn't have the mount point, so add it.
	for _, vName := range vols {
		// removes the prefix based on prefix length, only trim if it matches the prefix
		if strings.HasPrefix(vName, volumePrefix) {
			volumeName := vName[len(volumePrefix):]
			log.Debugf("List() adding volume: %v ", volumeName)

			v := &volume.Volume{Name: volumeName, Mountpoint: filepath.Join(volumeDir, vName)}
			volumes = append(volumes, v)
		} else {
			log.Debugf("wrong prefix, skipping Name: %v", vName)
		}
	}
	return volume.Response{Volumes: volumes}
}

func (d ndvpDriver) Capabilities(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Capabilities(%v)", r)

	return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}
