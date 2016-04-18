// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type IgroupCreateRequest struct {
	XMLName xml.Name `xml:"igroup-create"`

	BindPortsetPtr        *string `xml:"bind-portset"`
	InitiatorGroupNamePtr *string `xml:"initiator-group-name"`
	InitiatorGroupTypePtr *string `xml:"initiator-group-type"`
	OsTypePtr             *string `xml:"os-type"`
	OstypePtr             *string `xml:"ostype"`
}

func (o *IgroupCreateRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewIgroupCreateRequest() *IgroupCreateRequest { return &IgroupCreateRequest{} }

func (r *IgroupCreateRequest) ExecuteUsing(zr *ZapiRunner) (IgroupCreateResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n IgroupCreateResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("igroup-create result:\n%s", n.Result)

	return n, err
}

func (o IgroupCreateRequest) String() string {
	var buffer bytes.Buffer
	if o.BindPortsetPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "bind-portset", *o.BindPortsetPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("bind-portset: nil\n"))
	}
	if o.InitiatorGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-name", *o.InitiatorGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-name: nil\n"))
	}
	if o.InitiatorGroupTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-type", *o.InitiatorGroupTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-type: nil\n"))
	}
	if o.OsTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "os-type", *o.OsTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("os-type: nil\n"))
	}
	if o.OstypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "ostype", *o.OstypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("ostype: nil\n"))
	}
	return buffer.String()
}

func (o *IgroupCreateRequest) BindPortset() string {
	r := *o.BindPortsetPtr
	return r
}

func (o *IgroupCreateRequest) SetBindPortset(newValue string) *IgroupCreateRequest {
	o.BindPortsetPtr = &newValue
	return o
}

func (o *IgroupCreateRequest) InitiatorGroupName() string {
	r := *o.InitiatorGroupNamePtr
	return r
}

func (o *IgroupCreateRequest) SetInitiatorGroupName(newValue string) *IgroupCreateRequest {
	o.InitiatorGroupNamePtr = &newValue
	return o
}

func (o *IgroupCreateRequest) InitiatorGroupType() string {
	r := *o.InitiatorGroupTypePtr
	return r
}

func (o *IgroupCreateRequest) SetInitiatorGroupType(newValue string) *IgroupCreateRequest {
	o.InitiatorGroupTypePtr = &newValue
	return o
}

func (o *IgroupCreateRequest) OsType() string {
	r := *o.OsTypePtr
	return r
}

func (o *IgroupCreateRequest) SetOsType(newValue string) *IgroupCreateRequest {
	o.OsTypePtr = &newValue
	return o
}

func (o *IgroupCreateRequest) Ostype() string {
	r := *o.OstypePtr
	return r
}

func (o *IgroupCreateRequest) SetOstype(newValue string) *IgroupCreateRequest {
	o.OstypePtr = &newValue
	return o
}

type IgroupCreateResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result IgroupCreateResponseResult `xml:"results"`
}

func (o IgroupCreateResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type IgroupCreateResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *IgroupCreateResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewIgroupCreateResponse() *IgroupCreateResponse { return &IgroupCreateResponse{} }

func (o IgroupCreateResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
