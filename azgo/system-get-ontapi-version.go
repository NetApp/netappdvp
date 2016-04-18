// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type SystemGetOntapiVersionRequest struct {
	XMLName xml.Name `xml:"system-get-ontapi-version"`
}

func (o *SystemGetOntapiVersionRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewSystemGetOntapiVersionRequest() *SystemGetOntapiVersionRequest {
	return &SystemGetOntapiVersionRequest{}
}

func (r *SystemGetOntapiVersionRequest) ExecuteUsing(zr *ZapiRunner) (SystemGetOntapiVersionResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n SystemGetOntapiVersionResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("system-get-ontapi-version result:\n%s", n.Result)

	return n, err
}

func (o SystemGetOntapiVersionRequest) String() string {
	var buffer bytes.Buffer
	return buffer.String()
}

type SystemGetOntapiVersionResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result SystemGetOntapiVersionResponseResult `xml:"results"`
}

func (o SystemGetOntapiVersionResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type SystemGetOntapiVersionResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
	MajorVersionPtr  *int   `xml:"major-version"`
	MinorVersionPtr  *int   `xml:"minor-version"`
}

func (o *SystemGetOntapiVersionResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewSystemGetOntapiVersionResponse() *SystemGetOntapiVersionResponse {
	return &SystemGetOntapiVersionResponse{}
}

func (o SystemGetOntapiVersionResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.MajorVersionPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "major-version", *o.MajorVersionPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("major-version: nil\n"))
	}
	if o.MinorVersionPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "minor-version", *o.MinorVersionPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("minor-version: nil\n"))
	}
	return buffer.String()
}

func (o *SystemGetOntapiVersionResponseResult) MajorVersion() int {
	r := *o.MajorVersionPtr
	return r
}

func (o *SystemGetOntapiVersionResponseResult) SetMajorVersion(newValue int) *SystemGetOntapiVersionResponseResult {
	o.MajorVersionPtr = &newValue
	return o
}

func (o *SystemGetOntapiVersionResponseResult) MinorVersion() int {
	r := *o.MinorVersionPtr
	return r
}

func (o *SystemGetOntapiVersionResponseResult) SetMinorVersion(newValue int) *SystemGetOntapiVersionResponseResult {
	o.MinorVersionPtr = &newValue
	return o
}
