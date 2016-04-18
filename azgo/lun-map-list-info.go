// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type LunMapListInfoRequest struct {
	XMLName xml.Name `xml:"lun-map-list-info"`

	PathPtr *string `xml:"path"`
}

func (o *LunMapListInfoRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewLunMapListInfoRequest() *LunMapListInfoRequest { return &LunMapListInfoRequest{} }

func (r *LunMapListInfoRequest) ExecuteUsing(zr *ZapiRunner) (LunMapListInfoResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n LunMapListInfoResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("lun-map-list-info result:\n%s", n.Result)

	return n, err
}

func (o LunMapListInfoRequest) String() string {
	var buffer bytes.Buffer
	if o.PathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "path", *o.PathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("path: nil\n"))
	}
	return buffer.String()
}

func (o *LunMapListInfoRequest) Path() string {
	r := *o.PathPtr
	return r
}

func (o *LunMapListInfoRequest) SetPath(newValue string) *LunMapListInfoRequest {
	o.PathPtr = &newValue
	return o
}

type LunMapListInfoResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result LunMapListInfoResponseResult `xml:"results"`
}

func (o LunMapListInfoResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type LunMapListInfoResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr   string                   `xml:"status,attr"`
	ResultReasonAttr   string                   `xml:"reason,attr"`
	ResultErrnoAttr    string                   `xml:"errno,attr"`
	InitiatorGroupsPtr []InitiatorGroupInfoType `xml:"initiator-groups>initiator-group-info"`
}

func (o *LunMapListInfoResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewLunMapListInfoResponse() *LunMapListInfoResponse { return &LunMapListInfoResponse{} }

func (o LunMapListInfoResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.InitiatorGroupsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-groups", o.InitiatorGroupsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-groups: nil\n"))
	}
	return buffer.String()
}

func (o *LunMapListInfoResponseResult) InitiatorGroups() []InitiatorGroupInfoType {
	r := o.InitiatorGroupsPtr
	return r
}

func (o *LunMapListInfoResponseResult) SetInitiatorGroups(newValue []InitiatorGroupInfoType) *LunMapListInfoResponseResult {
	newSlice := make([]InitiatorGroupInfoType, len(newValue))
	copy(newSlice, newValue)
	o.InitiatorGroupsPtr = newSlice
	return o
}
