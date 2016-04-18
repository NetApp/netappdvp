// Copyright 2016 NetApp, Inc. All Rights Reserved.

package main

import (
	"flag"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestMain(m *testing.M) {

	// TODO validate pre-reqs for tests;  for instance, need a volume called 'v' for ontap_test.go
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	os.Exit(m.Run())
}
