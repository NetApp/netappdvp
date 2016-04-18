// Copyright 2016 NetApp, Inc. All Rights Reserved.

package ontap

import (
	"flag"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestMain(m *testing.M) {

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	os.Exit(m.Run())
}
