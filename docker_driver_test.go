// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
  "testing"
  "sync"
  "github.com/netapp/netappdvp/storage_drivers"
)

func newNdvpDriverWithPrefix(prefix string) (*ndvpDriver) {
  volumeDir := "/path/to/root"

  commonConfig := &storage_drivers.CommonStorageDriverConfig {}
  commonConfig.Version = 1
  commonConfig.StorageDriverName = "ontap-nas"
  commonConfig.Debug = false
  commonConfig.DisableDelete = false
  commonConfig.StoragePrefixRaw = []byte(prefix)

  d := &ndvpDriver {
    root: volumeDir,
    config: *commonConfig,
    m: &sync.Mutex{},
    sd: storage_drivers.Drivers[commonConfig.StorageDriverName],
  }

  return d
}

func TestStoragePrefix(t *testing.T) {
  prefix_cases := []struct {
    in, expected_prefix string
  } {
    {``, "netappdvp_"},
    {`"myprefix_"`, "myprefix_"},
    {`""`, ""},
  }

  for _, c := range prefix_cases {
    driver := newNdvpDriverWithPrefix(c.in)
    got := driver.volumePrefix()
    if got != c.expected_prefix {
      t.Errorf("ndvpDriver.volumePrefix() == %q, expected %q", got, c.expected_prefix)
    }
  }
}

func TestVolumeNames(t *testing.T) {
  volume_name_cases := []struct {
    prefix, volume, expected_volume_name string
  } {
    {``, "vol1", "netappdvp_vol1"},
    {`"myprefix_"`, "vol2", "myprefix_vol2"},
    {`""`, "vol3", "vol3"},
  }

  for _, c := range volume_name_cases {
    driver := newNdvpDriverWithPrefix(c.prefix)
    got := driver.volumeName(c.volume)
    if got != c.expected_volume_name {
      t.Errorf("ndvpDriver.volumeName(%q) == %q, expected %q", c.volume, got, c.expected_volume_name)
    }
  }
}
