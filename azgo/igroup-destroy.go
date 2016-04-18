// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type IgroupDestroyRequest struct {
	XMLName xml.Name `xml:"igroup-destroy"`

	ForcePtr              *bool   `xml:"force"`
	InitiatorGroupNamePtr *string `xml:"initiator-group-name"`
}

func (o *IgroupDestroyRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewIgroupDestroyRequest() *IgroupDestroyRequest { return &IgroupDestroyRequest{} }

func (r *IgroupDestroyRequest) ExecuteUsing(zr *ZapiRunner) (IgroupDestroyResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n IgroupDestroyResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("igroup-destroy result:\n%s", n.Result)

	return n, err
}

func (o IgroupDestroyRequest) String() string {
	var buffer bytes.Buffer
	if o.ForcePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force", *o.ForcePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force: nil\n"))
	}
	if o.InitiatorGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-name", *o.InitiatorGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-name: nil\n"))
	}
	return buffer.String()
}

func (o *IgroupDestroyRequest) Force() bool {
	r := *o.ForcePtr
	return r
}

func (o *IgroupDestroyRequest) SetForce(newValue bool) *IgroupDestroyRequest {
	o.ForcePtr = &newValue
	return o
}

func (o *IgroupDestroyRequest) InitiatorGroupName() string {
	r := *o.InitiatorGroupNamePtr
	return r
}

func (o *IgroupDestroyRequest) SetInitiatorGroupName(newValue string) *IgroupDestroyRequest {
	o.InitiatorGroupNamePtr = &newValue
	return o
}

type IgroupDestroyResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result IgroupDestroyResponseResult `xml:"results"`
}

func (o IgroupDestroyResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type IgroupDestroyResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *IgroupDestroyResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewIgroupDestroyResponse() *IgroupDestroyResponse { return &IgroupDestroyResponse{} }

func (o IgroupDestroyResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
