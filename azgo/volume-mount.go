// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeMountRequest struct {
	XMLName xml.Name `xml:"volume-mount"`

	ActivateJunctionPtr     *bool   `xml:"activate-junction"`
	ExportPolicyOverridePtr *bool   `xml:"export-policy-override"`
	JunctionPathPtr         *string `xml:"junction-path"`
	VolumeNamePtr           *string `xml:"volume-name"`
}

func (o *VolumeMountRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeMountRequest() *VolumeMountRequest { return &VolumeMountRequest{} }

func (r *VolumeMountRequest) ExecuteUsing(zr *ZapiRunner) (VolumeMountResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeMountResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-mount result:\n%s", n.Result)

	return n, err
}

func (o VolumeMountRequest) String() string {
	var buffer bytes.Buffer
	if o.ActivateJunctionPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "activate-junction", *o.ActivateJunctionPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("activate-junction: nil\n"))
	}
	if o.ExportPolicyOverridePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "export-policy-override", *o.ExportPolicyOverridePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("export-policy-override: nil\n"))
	}
	if o.JunctionPathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "junction-path", *o.JunctionPathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("junction-path: nil\n"))
	}
	if o.VolumeNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-name", *o.VolumeNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-name: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeMountRequest) ActivateJunction() bool {
	r := *o.ActivateJunctionPtr
	return r
}

func (o *VolumeMountRequest) SetActivateJunction(newValue bool) *VolumeMountRequest {
	o.ActivateJunctionPtr = &newValue
	return o
}

func (o *VolumeMountRequest) ExportPolicyOverride() bool {
	r := *o.ExportPolicyOverridePtr
	return r
}

func (o *VolumeMountRequest) SetExportPolicyOverride(newValue bool) *VolumeMountRequest {
	o.ExportPolicyOverridePtr = &newValue
	return o
}

func (o *VolumeMountRequest) JunctionPath() string {
	r := *o.JunctionPathPtr
	return r
}

func (o *VolumeMountRequest) SetJunctionPath(newValue string) *VolumeMountRequest {
	o.JunctionPathPtr = &newValue
	return o
}

func (o *VolumeMountRequest) VolumeName() string {
	r := *o.VolumeNamePtr
	return r
}

func (o *VolumeMountRequest) SetVolumeName(newValue string) *VolumeMountRequest {
	o.VolumeNamePtr = &newValue
	return o
}

type VolumeMountResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeMountResponseResult `xml:"results"`
}

func (o VolumeMountResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeMountResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *VolumeMountResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeMountResponse() *VolumeMountResponse { return &VolumeMountResponse{} }

func (o VolumeMountResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
