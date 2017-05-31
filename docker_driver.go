// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/netapp/netappdvp/storage_drivers"
	"github.com/netapp/netappdvp/utils"
)

type ndvpDriver struct {
	m      *sync.Mutex
	root   string
	config storage_drivers.CommonStorageDriverConfig
	sd     storage_drivers.StorageDriver
}

func (d *ndvpDriver) volumeName(name string) string {
	prefixToUse := *d.config.StoragePrefix
	if strings.HasPrefix(name, prefixToUse) {
		return name
	}
	return prefixToUse + name
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
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Create(%v)", r)

	opts := r.Options
	target := d.volumeName(r.Name)

	sizeBytes, err := utils.GetVolumeSizeBytes(opts, d.config.Size)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Error creating volume: %v", err)}
	}
	var createErr error

	// If 'from' is specified, create a snapshot and a clone rather than a new empty volume
	from := utils.GetV(opts, "from", "")
	if from != "" {
		source := d.volumeName(from)

		// If 'fromSnapshot' is specified, we use the existing snapshot instead
		snapshot := utils.GetV(opts, "fromSnapshot", "")
		createErr = d.sd.CreateClone(target, source, snapshot, opts)
	} else {
		createErr = d.sd.Create(target, sizeBytes, opts)
	}

	if createErr != nil {
		return volume.Response{Err: fmt.Sprintf("Error creating volume: %v", createErr)}
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

	// use the StorageDriver to destroy the storage objects
	destroyErr := d.sd.Destroy(target)
	if destroyErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem removing docker volume: %v error: %v", target, destroyErr)}
	}

	// Best effort removal of the mountpoint
	m := d.mountpoint(target)
	os.Remove(m)

	return volume.Response{}
}

func (d ndvpDriver) getPath(r volume.Request) (string, error) {
	// Currently, this returns the mountpoint based on whether the path exists.

	// Should it:
	// a) Also return what the mountpoint would be if it were mounted, even if it isn't?
	// b) Verify that the volume is actually mounted before returning it?
	// c) Stay as-is?

	target := d.volumeName(r.Name)
	m := d.mountpoint(target)

	log.Debugf("Getting path for volume '%s' as '%s'", target, m)

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
		return "", err
	}
	if fi == nil {
		return "", fmt.Errorf("Could not stat %v", m)
	}

	return m, nil
}

// Path is part of the core Docker API and is called to return the filesystem path to a docker volume
func (d ndvpDriver) Path(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Path(%v)", r)

	mountpoint, err := d.getPath(r)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{
		Mountpoint: mountpoint,
	}
}

// Get is part of the core Docker API and is called to return details about a docker volume
func (d ndvpDriver) Get(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Get(%v)", r)

	// Gather the target volume name as the storage sees it
	target := d.volumeName(r.Name)

	// Ask the storage driver whether the specified volume exists
	err := d.sd.Get(target)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	// Get the mountpoint, if this volume is mounted
	mountpoint, err := d.getPath(r)

	// Ask the storage driver for the list of snapshots associated with the volume
	snaps, err := d.sd.SnapshotList(target)

	// If we don't get any snapshots, that's fine. We'll return an empty list.
	status := map[string]interface{}{
		"Snapshots": snaps,
	}

	vol := &volume.Volume{
		Name:       r.Name,
		Mountpoint: mountpoint,
		Status:     status, // introduced in Docker 1.12, earlier versions ignore
	}

	return volume.Response{
		Volume: vol,
	}
}

// Mount is part of the core Docker API and is called to mount a docker volume
func (d ndvpDriver) Mount(r volume.MountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Mount(%v)", r)

	target := d.volumeName(r.Name)

	m := d.mountpoint(target)
	log.Debugf("Mounting volume %s on %s", target, m)

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

	// use the StorageDriver to unmount the storage objects
	detachErr := d.sd.Detach(target, m)
	if detachErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem unmounting docker volume: %v error: %v", target, detachErr)}
	}

	// Best effort removal of the mountpoint
	os.Remove(m)

	return volume.Response{}
}

// List is part of the core Docker API and is called to list all known docker volumes for this plugin
func (d ndvpDriver) List(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("List(%v)", r)

	var volumes []*volume.Volume
	vols, err := d.sd.List()
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Unable to retrieve volume list, error: %v", err)}
	}

	for _, vol := range vols {
		// What is the impact of leaving the mountpoints out of this response?
		v := &volume.Volume{Name: vol}
		volumes = append(volumes, v)
	}

	return volume.Response{Volumes: volumes}
}

func (d ndvpDriver) Capabilities(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Capabilities(%v)", r)

	return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}
