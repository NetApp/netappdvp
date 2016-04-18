// Copyright 2016 NetApp, Inc. All Rights Reserved.

package utils

import (
	"runtime"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestGetDFOutput(t *testing.T) {
	log.Debug("Running TestGetDFOutput...")

	if runtime.GOOS != "linux" {
		return
	}

	output, err := GetDFOutput()
	if output == nil || len(output) <= 0 || err != nil {
		t.Error("Could not validate `df` output.", err)
	}

	foundRootFS := false
	for _, e := range output {
		if e.Target == "/" {
			foundRootFS = true
		}
	}
	if !foundRootFS {
		t.Error("Could not find root filesystem, something's wrong.")
	}
}

func TestGetInitiatorIqns(t *testing.T) {
	log.Debug("Running TestGetInitiatorIqns...")

	if runtime.GOOS != "linux" {
		return
	}

	isIscsiSupported := IscsiSupported()
	if isIscsiSupported {
		iqns, err := GetInitiatorIqns()
		if iqns == nil || len(iqns) <= 0 || err != nil {
			t.Error("Could not lookup iqns.", err)
		}
		if iqns[0][0:4] != "iqn." {
			t.Error("Unexpected iqn found:", iqns[0][0:4])
		}
	} else {
		t.Skip("iSCSI support not found.")
	}
}

func TestIscsiSupported(t *testing.T) {
	log.Debug("Running TestIscsiSupported...")

	if runtime.GOOS != "linux" {
		return
	}

	isIscsiSupported := IscsiSupported()
	if !isIscsiSupported {
		t.Error("iSCSI support not found")
	}
}

func TestIscsiSessionExists(t *testing.T) {
	log.Debug("Running TestIscsiSessionExists...")

	if runtime.GOOS != "linux" {
		return
	}

	sessionExists, err := IscsiSessionExists("127.0.0.1")
	if err != nil {
		t.Errorf("Unexpected iSCSI session error: %v", err)
	}
	if sessionExists {
		t.Error("Unexpected iSCSI session found")
	}

	sessionExists, err = IscsiSessionExists("10.0.207.7")
	if err != nil {
		t.Errorf("Unexpected iSCSI session error: %v", err)
	}
	if !sessionExists {
		t.Error("Expected iSCSI session NOT found")
	}
}

func TestGetFSType(t *testing.T) {
	log.Debug("Running TestGetFSType...")

	if runtime.GOOS != "linux" {
		return
	}

	rootFSType := GetFSType("/dev/sda1")
	if rootFSType == "" {
		t.Error("Could not determine root fs type")
	}
}

func TestGetDeviceInfoForLuns(t *testing.T) {
	log.Debug("Running TestGetDeviceInfoForLuns...")

	if runtime.GOOS != "linux" {
		return
	}

	info, err := GetDeviceInfoForLuns()

	if err != nil {
		t.Errorf("Unexpected problem getting device info for LUNs, error: %v", err)
	}

	for _, e1 := range info {
		log.Debugf("Found %v", e1)
	}
}
