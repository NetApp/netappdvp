// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

// TestConstants is a simple sanity test
func TestConstants(t *testing.T) {
	log.Debug("Running TestSystemGetVersion...")

	if EONTAPI_EEXIST != "17" {
		t.Error("Unexpected constant value found for EONTAPI_EEXIST")
	}
}
