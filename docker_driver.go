// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/netapp/netappdvp/storage_drivers"
	"github.com/netapp/netappdvp/utils"

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
			log.Debugf("storagePrefix is specified as \"\", using no prefix")
		} else {
			// trim quotes from start and end of string
			prefixToUse = s[1 : len(s)-1]
			log.Debugf("storagePrefix is specified, using prefix: %v", prefixToUse)
		}
	} else {
		prefixToUse = defaultPrefix
		log.Debugf("storagePrefix is unspecified, using default prefix: %v", prefixToUse)
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

	// use the StorageDriver to create the storage objects
	createErr := d.sd.Create(target, opts)
	if createErr != nil {
		os.Remove(m)
		return volume.Response{Err: fmt.Sprintf("Error creating storage: %v", createErr)}
	}

	return volume.Response{}
}

// Remove is part of the core Docker API and is called to remove a docker volume
func (d ndvpDriver) Remove(r volume.Request) volume.Response {
	target := d.volumeName(r.Name)

	// allow user to completely disable volume deletion
	if d.config.DisableDelete {
		log.Infof("Skipping removal of %s because of user preference to disable volume deletion", target)
		return volume.Response{}
	}

	log.Debugf("Removing docker volume %s", target)

	d.m.Lock()
	defer d.m.Unlock()
	m := d.mountpoint(target)

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
		return volume.Response{} // nothing to do
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v is not a directory", m)}
	}

	// use the StorageDriver to destroy the storage objects
	destroyErr := d.sd.Destroy(target)
	if destroyErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem removing docker volume: %v error: %v", target, destroyErr)}
	}

	log.Debugf("rmdir(%s)", m)
	err3 := os.Remove(m)
	if err3 != nil {
		return volume.Response{Err: err3.Error()}
	}

	return volume.Response{}
}

func (d ndvpDriver) getPath(r volume.Request) (*volume.Volume, error) {
	target := d.volumeName(r.Name)
	m := d.mountpoint(target)
	log.Debugf("Getting path for volume '%s' as '%s'", target, m)

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
		return nil, err
	}
	if fi == nil {
		return nil, fmt.Errorf("Could not stat %v", m)
	}

	volume := &volume.Volume{
		Name:       r.Name,
		Mountpoint: d.mountpoint(target)}
	return volume, nil
}

// Path is part of the core Docker API and is called to return the filesystem path to a docker volume
func (d ndvpDriver) Path(r volume.Request) volume.Response {
	v, err := d.getPath(r)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{
		Mountpoint: v.Mountpoint,
	}
}

// Get is part of the core Docker API and is called to return the filesystem path to a docker volume
func (d ndvpDriver) Get(r volume.Request) volume.Response {
	v, err := d.getPath(r)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{
		Volume: v,
	}
}

type nfsMount struct {
	IP         string
	Target     string
	MountPoint string
	MachineID  string
}

// Mount is part of the core Docker API and is called to mount a docker volume
func (d ndvpDriver) Mount(r volume.Request) volume.Response {
	target := d.volumeName(r.Name)

	d.m.Lock()
	defer d.m.Unlock()
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
		return volume.Response{Err: fmt.Sprintf("Problem attaching docker volume: %v mountpoint: %v error: %v", target, m, attachErr)}
	}

	return volume.Response{Mountpoint: m}
}

// Unmount is part of the core Docker API and is called to unmount a docker volume
func (d ndvpDriver) Unmount(r volume.Request) volume.Response {
	target := d.volumeName(r.Name)

	d.m.Lock()
	defer d.m.Unlock()
	m := d.mountpoint(target)
	log.Debugf("Unmounting docker volume %s", target)

	// use the StorageDriver to unmount the storage objects
	detachErr := d.sd.Detach(target, m)
	if detachErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem unmounting docker volume: %v error: %v", target, detachErr)}
	}

	return volume.Response{}
}

// List is part of the core Docker API and is called to list all known docker volumes for this plugin
func (d ndvpDriver) List(r volume.Request) volume.Response {

	d.m.Lock()
	defer d.m.Unlock()
	log.Debugf("Listing volumes...")

	// open directory ...
	volumeDir := d.root
	dir, err := os.Open(volumeDir)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Problem opening directory %v, error: %v", volumeDir, err)}
	}
	defer dir.Close()

	// stat the directory
	fi, err := dir.Stat()
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Problem stating directory %v, error: %v", volumeDir, err)}
	}
	if !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v is not a directory!", volumeDir)}
	}

	// finally, we spin through all the subdirectories (if any) and return them in our List response
	var vols []*volume.Volume
	dirs := make([]string, 0)   // lint complains to switch to this, but it doens't work -> var dirs []string
	fis, err := dir.Readdir(-1) // -1 means return all the FileInfos
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Problem reading directory %v, error: %v", volumeDir, err)}
	}
	for _, fileinfo := range fis {
		if fileinfo.IsDir() {
			dirs = append(dirs, fileinfo.Name())

			// removes the prefix based on prefix length, for instance [10:] to remove 'netappdvp_' from start of name
			volumePrefix := d.volumePrefix()

			// only trim if it matches the prefix
			if strings.HasPrefix(fileinfo.Name(), volumePrefix) {
				volumeName := fileinfo.Name()[len(volumePrefix):]
				log.Debugf("List() adding volume: %v from: %v", volumeName, filepath.Join(volumeDir, fileinfo.Name()))

				v := &volume.Volume{Name: volumeName, Mountpoint: filepath.Join(volumeDir, fileinfo.Name())}
				vols = append(vols, v)
			} else {
				log.Debugf("wrong prefix, skipping fileinfo.Name: %v", fileinfo.Name())
			}
		}
	}

	return volume.Response{Volumes: vols}
}

func (d ndvpDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}
