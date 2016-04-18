// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type IgroupAddRequest struct {
	XMLName xml.Name `xml:"igroup-add"`

	ForcePtr              *bool   `xml:"force"`
	InitiatorPtr          *string `xml:"initiator"`
	InitiatorGroupNamePtr *string `xml:"initiator-group-name"`
}

func (o *IgroupAddRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewIgroupAddRequest() *IgroupAddRequest { return &IgroupAddRequest{} }

func (r *IgroupAddRequest) ExecuteUsing(zr *ZapiRunner) (IgroupAddResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n IgroupAddResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("igroup-add result:\n%s", n.Result)

	return n, err
}

func (o IgroupAddRequest) String() string {
	var buffer bytes.Buffer
	if o.ForcePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force", *o.ForcePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force: nil\n"))
	}
	if o.InitiatorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator", *o.InitiatorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator: nil\n"))
	}
	if o.InitiatorGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-name", *o.InitiatorGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-name: nil\n"))
	}
	return buffer.String()
}

func (o *IgroupAddRequest) Force() bool {
	r := *o.ForcePtr
	return r
}

func (o *IgroupAddRequest) SetForce(newValue bool) *IgroupAddRequest {
	o.ForcePtr = &newValue
	return o
}

func (o *IgroupAddRequest) Initiator() string {
	r := *o.InitiatorPtr
	return r
}

func (o *IgroupAddRequest) SetInitiator(newValue string) *IgroupAddRequest {
	o.InitiatorPtr = &newValue
	return o
}

func (o *IgroupAddRequest) InitiatorGroupName() string {
	r := *o.InitiatorGroupNamePtr
	return r
}

func (o *IgroupAddRequest) SetInitiatorGroupName(newValue string) *IgroupAddRequest {
	o.InitiatorGroupNamePtr = &newValue
	return o
}

type IgroupAddResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result IgroupAddResponseResult `xml:"results"`
}

func (o IgroupAddResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type IgroupAddResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *IgroupAddResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewIgroupAddResponse() *IgroupAddResponse { return &IgroupAddResponse{} }

func (o IgroupAddResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
