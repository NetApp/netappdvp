// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeUnmountRequest struct {
	XMLName xml.Name `xml:"volume-unmount"`

	ForcePtr      *bool   `xml:"force"`
	VolumeNamePtr *string `xml:"volume-name"`
}

func (o *VolumeUnmountRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeUnmountRequest() *VolumeUnmountRequest { return &VolumeUnmountRequest{} }

func (r *VolumeUnmountRequest) ExecuteUsing(zr *ZapiRunner) (VolumeUnmountResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeUnmountResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-unmount result:\n%s", n.Result)

	return n, err
}

func (o VolumeUnmountRequest) String() string {
	var buffer bytes.Buffer
	if o.ForcePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force", *o.ForcePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force: nil\n"))
	}
	if o.VolumeNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-name", *o.VolumeNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-name: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeUnmountRequest) Force() bool {
	r := *o.ForcePtr
	return r
}

func (o *VolumeUnmountRequest) SetForce(newValue bool) *VolumeUnmountRequest {
	o.ForcePtr = &newValue
	return o
}

func (o *VolumeUnmountRequest) VolumeName() string {
	r := *o.VolumeNamePtr
	return r
}

func (o *VolumeUnmountRequest) SetVolumeName(newValue string) *VolumeUnmountRequest {
	o.VolumeNamePtr = &newValue
	return o
}

type VolumeUnmountResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeUnmountResponseResult `xml:"results"`
}

func (o VolumeUnmountResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeUnmountResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *VolumeUnmountResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeUnmountResponse() *VolumeUnmountResponse { return &VolumeUnmountResponse{} }

func (o VolumeUnmountResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
