// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

// VolumeCloneGetRequest is a structure to represent a volume-clone-get ZAPI request object
type VolumeCloneGetRequest struct {
	XMLName xml.Name `xml:"volume-clone-get"`

	DesiredAttributesPtr *VolumeCloneInfoType    `xml:"desired-attributes>volume-clone-info"`
	VolumePtr            *string                 `xml:"volume"`
}

// ToXML converts this object into an xml string representation
func (o *VolumeCloneGetRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

// NewVolumeCloneGetRequest is a factory method for creating new instances of VolumeCloneGetRequest objects
func NewVolumeCloneGetRequest() *VolumeCloneGetRequest { return &VolumeCloneGetRequest{} }

// ExecuteUsing converts this object to a ZAPI XML representation and uses the supplied ZapiRunner to send to a filer
func (o *VolumeCloneGetRequest) ExecuteUsing(zr *ZapiRunner) (VolumeCloneGetResponse, error) {
	resp, err := zr.SendZapi(o)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeCloneGetResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-clone-get result:\n%s", n.Result)

	return n, err
}

// String returns a string representation of this object's fields and implements the Stringer interface
func (o VolumeCloneGetRequest) String() string {
	var buffer bytes.Buffer
	if o.DesiredAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "desired-attributes", *o.DesiredAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("desired-attributes: nil\n"))
	}
	if o.VolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume", *o.VolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("tag: nil\n"))
	}
	return buffer.String()
}

// DesiredAttributes is a fluent style 'getter' method that can be chained
func (o *VolumeCloneGetRequest) DesiredAttributes() VolumeCloneInfoType {
	r := *o.DesiredAttributesPtr
	return r
}

// SetDesiredAttributes is a fluent style 'setter' method that can be chained
func (o *VolumeCloneGetRequest) SetDesiredAttributes(newValue VolumeCloneInfoType) *VolumeCloneGetRequest {
	o.DesiredAttributesPtr = &newValue
	return o
}

// Volume is a fluent style 'getter' method that can be chained
func (o *VolumeCloneGetRequest) Volume() string {
	r := *o.VolumePtr
	return r
}

// SetVolume is a fluent style 'setter' method that can be chained
func (o *VolumeCloneGetRequest) SetVolume(newValue string) *VolumeCloneGetRequest {
	o.VolumePtr = &newValue
	return o
}

// VolumeCloneGetResponse is a structure to represent a volume-clone-get ZAPI response object
type VolumeCloneGetResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeCloneGetResponseResult `xml:"results"`
}

// String returns a string representation of this object's fields and implements the Stringer interface
func (o VolumeCloneGetResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

// VolumeCloneGetResponseResult is a structure to represent a volume-clone-get ZAPI object's result
type VolumeCloneGetResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr  string                   `xml:"status,attr"`
	ResultReasonAttr  string                   `xml:"reason,attr"`
	ResultErrnoAttr   string                   `xml:"errno,attr"`
	AttributesPtr     *VolumeCloneInfoType     `xml:"attributes>volume-clone-info"`
}

// ToXML converts this object into an xml string representation
func (o *VolumeCloneGetResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

// NewVolumeCloneGetResponse is a factory method for creating new instances of VolumeCloneGetResponse objects
func NewVolumeCloneGetResponse() *VolumeCloneGetResponse { return &VolumeCloneGetResponse{} }

// String returns a string representation of this object's fields and implements the Stringer interface
func (o VolumeCloneGetResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.AttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "attributes", o.AttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("attributes: nil\n"))
	}
	return buffer.String()
}

// Attributes is a fluent style 'getter' method that can be chained
func (o *VolumeCloneGetResponseResult) Attributes() VolumeCloneInfoType {
	r := *o.AttributesPtr
	return r
}

// SetAttributes is a fluent style 'setter' method that can be chained
func (o *VolumeCloneGetResponseResult) SetAttributes(newValue VolumeCloneInfoType) *VolumeCloneGetResponseResult {
	o.AttributesPtr = &newValue
	return o
}
