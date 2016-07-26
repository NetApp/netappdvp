// Copyright 2016 NetApp, Inc. All Rights Reserved.

package eseries

import ()

var GenericResponseOkay int = 200
var GenericResponseSuccess int = 201
var GenericResponseNoContent int = 204
var GenericResponseNotFound int = 404
var GenericResponseOffline int = 424
var GenericResponseMalformed int = 422

//Add array to Web Services Proxy
type MsgConnect struct {
	ControllerAddresses []string `json:"controllerAddresses"`
	Password            string   `json:"password,omitempty"`
}

type MsgConnectResponse struct {
	ArrayID       string `json:"id"`
	AlreadyExists bool   `json:"alreadyExists"`
}

//Obtain volume group information
type VolumeGroupExResponse struct {
	SequenceNumber int    `json:"sequenceNum"`
	IsOffline      bool   `json:"offline"`
	WorldWideName  string `json:"worldWideName"`
	VolumeGroupRef string `json:"volumeGroupRef"`
	VolumeLabel    string `json:"label"`
	FreeSpace      string `json:"freeSpace"` //Documentation says this is an int but really it is a string!
}

//Create a volume
type MsgVolumeEx struct {
	VolumeGroupRef   string      `json:"poolId"`
	Name             string      `json:"name"`
	SizeUnit         string      `json:"sizeUnit"` //bytes, b, kb, mb, gb, tb, pb, eb, zb, yb
	Size             int         `json:"size"`
	SegmentSize      int         `json:"segSize"`
	DataAssurance    bool        `json:"dataAssuranceEnabled,omitempty"`
	OwningController string      `json:"owningControllerId,omitempty"`
	VolumeTags       []VolumeTag `json:"metaTags,omitempty"`
}

//Volume Metadata Tag
type VolumeTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MsgVolumeExResponse struct {
	IsOffline      bool         `json:"offline"`
	Label          string       `json:"label"`
	VolumeSize     string       `json:"capacity"`
	SegmentSize    int          `json:"segmentSize"`
	VolumeRef      string       `json:"volumeRef"`
	VolumeGroupRef string       `json:"volumeGroupRef"`
	ListOfMappings []LUNMapping `json:"listOfMappings"`
	IsMapped       bool         `json:"mapped"`
}

//Obtain information about all hosts on array
type HostExResponse struct {
	HostRef    string            `json:"hostRef"`
	Label      string            `json:"label"`
	Initiators []HostExInitiator `json:"initiators"`
}

type HostExInitiator struct {
	InitiatorRef string             `json:"initiatorRef"`
	NodeName     HostExScsiNodeName `json:"nodeName"`
	Label        string             `json:"label"`
}

type HostExScsiNodeName struct {
	IoInterfaceType string `json:"ioInterfaceType"` //scsi, fc, sata, iscsi, ib, fcoe, __UNDEFINED
	IscsiNodeName   string `json:"iscsiNodeName"`   //IQN from host
}

type HostExHostPort struct {
	//TODO - I think this would be used to support Fiber Channel
}

//Request to map a created volume to a host
type VolumeMappingCreateRequest struct {
	MappableObjectId string `json:"mappableObjectId"`
	TargetID         string `json:"targetId"`
	LunNumber        int    `json:"lun,omitempty"`
}

//Structure that reflects LUN information
type LUNMapping struct {
	LunMappingRef string `json:"lunMappingRef"`
	LunNumber     int    `json:"lun"`
	VolumeRef     string `json:"volumeRef"`
	HostRef       string `json:"mapRef"`
}

//Used for errors on RESTful calls to return what went wrong
type CallResponseError struct {
	ErrorMsg     string `json:"errorMessage"`
	LocalizedMsg string `json:"localizedMessage"`
	ReturnCode   string `json:"retcode"`
	CodeType     string `json:"codeType"` //'symbol', 'webservice', 'systemerror', 'devicemgrerror'
}
