package main

import (
	"fmt"

	"github.com/netapp/netappdvp/storage_drivers"
)

// GitHash is the git hash the binary was built from
var GitHash string

// BuildType is the type of build: custom, beta or stable
var BuildType string

// BuildTypeRev is the revision of the build
var BuildTypeRev string

// BuildTime is the time the binary was built
var BuildTime string

func init() {
	if GitHash == "" {
		GitHash = "unknown"
	}
	if BuildType == "" {
		BuildType = "custom"
	}
	if BuildTypeRev == "" {
		BuildTypeRev = "0"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if BuildType != "stable" {
		if BuildType == "custom" {
			storage_drivers.FullDriverVersion = fmt.Sprintf("%v-%v", storage_drivers.DriverVersion, BuildType)
		} else {
			storage_drivers.FullDriverVersion = fmt.Sprintf("%v-%v.%v", storage_drivers.DriverVersion, BuildType, BuildTypeRev)
		}
	}
	storage_drivers.BuildVersion = fmt.Sprintf("%v+%v", storage_drivers.FullDriverVersion, GitHash)
	storage_drivers.BuildTime = BuildTime
}
