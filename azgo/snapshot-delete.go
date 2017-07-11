// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

// SnapshotDeleteRequest is a structure to represent a snapshot-delete ZAPI request object
type SnapshotDeleteRequest struct {
	XMLName xml.Name `xml:"snapshot-delete"`

	IgnoreOwnersPtr         *bool   `xml:"ignore-owners"`
	SnapshotPtr             *string `xml:"snapshot"`
	SnapshotInstanceUUIDPtr *string `xml:"snapshot-instance-uuid"`
	VolumePtr               *string `xml:"volume"`
}

// ToXML converts this object into an xml string representation
func (o *SnapshotDeleteRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

// NewSnapshotDeleteRequest is a factory method for creating new instances of SnapshotDeleteRequest objects
func NewSnapshotDeleteRequest() *SnapshotDeleteRequest { return &SnapshotDeleteRequest{} }

// ExecuteUsing converts this object to a ZAPI XML representation and uses the supplied ZapiRunner to send to a filer
func (o *SnapshotDeleteRequest) ExecuteUsing(zr *ZapiRunner) (SnapshotDeleteResponse, error) {
	resp, err := zr.SendZapi(o)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n SnapshotDeleteResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("snapshot-delete result:\n%s", n.Result)

	return n, err
}

// String returns a string representation of this object's fields and implements the Stringer interface
func (o SnapshotDeleteRequest) String() string {
	var buffer bytes.Buffer
	if o.IgnoreOwnersPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "ignore-owners", *o.IgnoreOwnersPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("ignore-owners: nil\n"))
	}
	if o.SnapshotPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot", *o.SnapshotPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot: nil\n"))
	}
	if o.SnapshotInstanceUUIDPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-instance-uuid", *o.SnapshotInstanceUUIDPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-instance-uuid: nil\n"))
	}
	if o.VolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume", *o.VolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume: nil\n"))
	}
	return buffer.String()
}

// IgnoreOwners is a fluent style 'getter' method that can be chained
func (o *SnapshotDeleteRequest) IgnoreOwners() bool {
	r := *o.IgnoreOwnersPtr
	return r
}

// SetIgnoreOwners is a fluent style 'setter' method that can be chained
func (o *SnapshotDeleteRequest) SetIgnoreOwners(newValue bool) *SnapshotDeleteRequest {
	o.IgnoreOwnersPtr = &newValue
	return o
}

// Snapshot is a fluent style 'getter' method that can be chained
func (o *SnapshotDeleteRequest) Snapshot() string {
	r := *o.SnapshotPtr
	return r
}

// SetSnapshot is a fluent style 'setter' method that can be chained
func (o *SnapshotDeleteRequest) SetSnapshot(newValue string) *SnapshotDeleteRequest {
	o.SnapshotPtr = &newValue
	return o
}

// SnapshotInstanceUUID is a fluent style 'getter' method that can be chained
func (o *SnapshotDeleteRequest) SnapshotInstanceUUID() string {
	r := *o.SnapshotInstanceUUIDPtr
	return r
}

// SetSnapshotInstanceUUID is a fluent style 'setter' method that can be chained
func (o *SnapshotDeleteRequest) SetSnapshotInstanceUUID(newValue string) *SnapshotDeleteRequest {
	o.SnapshotInstanceUUIDPtr = &newValue
	return o
}

// Volume is a fluent style 'getter' method that can be chained
func (o *SnapshotDeleteRequest) Volume() string {
	r := *o.VolumePtr
	return r
}

// SetVolume is a fluent style 'setter' method that can be chained
func (o *SnapshotDeleteRequest) SetVolume(newValue string) *SnapshotDeleteRequest {
	o.VolumePtr = &newValue
	return o
}

// SnapshotDeleteResponse is a structure to represent a snapshot-delete ZAPI response object
type SnapshotDeleteResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result SnapshotDeleteResponseResult `xml:"results"`
}

// String returns a string representation of this object's fields and implements the Stringer interface
func (o SnapshotDeleteResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

// SnapshotDeleteResponseResult is a structure to represent a snapshot-delete ZAPI object's result
type SnapshotDeleteResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

// ToXML converts this object into an xml string representation
func (o *SnapshotDeleteResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

// NewSnapshotDeleteResponse is a factory method for creating new instances of SnapshotDeleteResponse objects
func NewSnapshotDeleteResponse() *SnapshotDeleteResponse { return &SnapshotDeleteResponse{} }

// String returns a string representation of this object's fields and implements the Stringer interface
func (o SnapshotDeleteResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
