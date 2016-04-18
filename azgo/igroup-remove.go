// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type IgroupRemoveRequest struct {
	XMLName xml.Name `xml:"igroup-remove"`

	ForcePtr              *bool   `xml:"force"`
	InitiatorPtr          *string `xml:"initiator"`
	InitiatorGroupNamePtr *string `xml:"initiator-group-name"`
}

func (o *IgroupRemoveRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewIgroupRemoveRequest() *IgroupRemoveRequest { return &IgroupRemoveRequest{} }

func (r *IgroupRemoveRequest) ExecuteUsing(zr *ZapiRunner) (IgroupRemoveResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n IgroupRemoveResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("igroup-remove result:\n%s", n.Result)

	return n, err
}

func (o IgroupRemoveRequest) String() string {
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

func (o *IgroupRemoveRequest) Force() bool {
	r := *o.ForcePtr
	return r
}

func (o *IgroupRemoveRequest) SetForce(newValue bool) *IgroupRemoveRequest {
	o.ForcePtr = &newValue
	return o
}

func (o *IgroupRemoveRequest) Initiator() string {
	r := *o.InitiatorPtr
	return r
}

func (o *IgroupRemoveRequest) SetInitiator(newValue string) *IgroupRemoveRequest {
	o.InitiatorPtr = &newValue
	return o
}

func (o *IgroupRemoveRequest) InitiatorGroupName() string {
	r := *o.InitiatorGroupNamePtr
	return r
}

func (o *IgroupRemoveRequest) SetInitiatorGroupName(newValue string) *IgroupRemoveRequest {
	o.InitiatorGroupNamePtr = &newValue
	return o
}

type IgroupRemoveResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result IgroupRemoveResponseResult `xml:"results"`
}

func (o IgroupRemoveResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type IgroupRemoveResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *IgroupRemoveResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewIgroupRemoveResponse() *IgroupRemoveResponse { return &IgroupRemoveResponse{} }

func (o IgroupRemoveResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
