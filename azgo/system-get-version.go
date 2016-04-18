// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type SystemGetVersionRequest struct {
	XMLName xml.Name `xml:"system-get-version"`
}

func (o *SystemGetVersionRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewSystemGetVersionRequest() *SystemGetVersionRequest { return &SystemGetVersionRequest{} }

func (r *SystemGetVersionRequest) ExecuteUsing(zr *ZapiRunner) (SystemGetVersionResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n SystemGetVersionResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("system-get-version result:\n%s", n.Result)

	return n, err
}

func (o SystemGetVersionRequest) String() string {
	var buffer bytes.Buffer
	return buffer.String()
}

type SystemGetVersionResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result SystemGetVersionResponseResult `xml:"results"`
}

func (o SystemGetVersionResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type SystemGetVersionResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr  string                  `xml:"status,attr"`
	ResultReasonAttr  string                  `xml:"reason,attr"`
	ResultErrnoAttr   string                  `xml:"errno,attr"`
	BuildTimestampPtr *int                    `xml:"build-timestamp"`
	IsClusteredPtr    *bool                   `xml:"is-clustered"`
	VersionPtr        *string                 `xml:"version"`
	VersionTuplePtr   *SystemVersionTupleType `xml:"version-tuple>system-version-tuple"`
}

func (o *SystemGetVersionResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewSystemGetVersionResponse() *SystemGetVersionResponse { return &SystemGetVersionResponse{} }

func (o SystemGetVersionResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.BuildTimestampPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "build-timestamp", *o.BuildTimestampPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("build-timestamp: nil\n"))
	}
	if o.IsClusteredPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-clustered", *o.IsClusteredPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-clustered: nil\n"))
	}
	if o.VersionPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "version", *o.VersionPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("version: nil\n"))
	}
	if o.VersionTuplePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "version-tuple", *o.VersionTuplePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("version-tuple: nil\n"))
	}
	return buffer.String()
}

func (o *SystemGetVersionResponseResult) BuildTimestamp() int {
	r := *o.BuildTimestampPtr
	return r
}

func (o *SystemGetVersionResponseResult) SetBuildTimestamp(newValue int) *SystemGetVersionResponseResult {
	o.BuildTimestampPtr = &newValue
	return o
}

func (o *SystemGetVersionResponseResult) IsClustered() bool {
	r := *o.IsClusteredPtr
	return r
}

func (o *SystemGetVersionResponseResult) SetIsClustered(newValue bool) *SystemGetVersionResponseResult {
	o.IsClusteredPtr = &newValue
	return o
}

func (o *SystemGetVersionResponseResult) Version() string {
	r := *o.VersionPtr
	return r
}

func (o *SystemGetVersionResponseResult) SetVersion(newValue string) *SystemGetVersionResponseResult {
	o.VersionPtr = &newValue
	return o
}

func (o *SystemGetVersionResponseResult) VersionTuple() SystemVersionTupleType {
	r := *o.VersionTuplePtr
	return r
}

func (o *SystemGetVersionResponseResult) SetVersionTuple(newValue SystemVersionTupleType) *SystemGetVersionResponseResult {
	o.VersionTuplePtr = &newValue
	return o
}
