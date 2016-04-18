// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type NetInterfaceGetIterRequest struct {
	XMLName xml.Name `xml:"net-interface-get-iter"`

	DesiredAttributesPtr *NetInterfaceInfoType `xml:"desired-attributes>net-interface-info"`
	MaxRecordsPtr        *int                  `xml:"max-records"`
	QueryPtr             *NetInterfaceInfoType `xml:"query>net-interface-info"`
	TagPtr               *string               `xml:"tag"`
}

func (o *NetInterfaceGetIterRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewNetInterfaceGetIterRequest() *NetInterfaceGetIterRequest { return &NetInterfaceGetIterRequest{} }

func (r *NetInterfaceGetIterRequest) ExecuteUsing(zr *ZapiRunner) (NetInterfaceGetIterResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n NetInterfaceGetIterResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("net-interface-get-iter result:\n%s", n.Result)

	return n, err
}

func (o NetInterfaceGetIterRequest) String() string {
	var buffer bytes.Buffer
	if o.DesiredAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "desired-attributes", *o.DesiredAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("desired-attributes: nil\n"))
	}
	if o.MaxRecordsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-records", *o.MaxRecordsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-records: nil\n"))
	}
	if o.QueryPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "query", *o.QueryPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("query: nil\n"))
	}
	if o.TagPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "tag", *o.TagPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("tag: nil\n"))
	}
	return buffer.String()
}

func (o *NetInterfaceGetIterRequest) DesiredAttributes() NetInterfaceInfoType {
	r := *o.DesiredAttributesPtr
	return r
}

func (o *NetInterfaceGetIterRequest) SetDesiredAttributes(newValue NetInterfaceInfoType) *NetInterfaceGetIterRequest {
	o.DesiredAttributesPtr = &newValue
	return o
}

func (o *NetInterfaceGetIterRequest) MaxRecords() int {
	r := *o.MaxRecordsPtr
	return r
}

func (o *NetInterfaceGetIterRequest) SetMaxRecords(newValue int) *NetInterfaceGetIterRequest {
	o.MaxRecordsPtr = &newValue
	return o
}

func (o *NetInterfaceGetIterRequest) Query() NetInterfaceInfoType {
	r := *o.QueryPtr
	return r
}

func (o *NetInterfaceGetIterRequest) SetQuery(newValue NetInterfaceInfoType) *NetInterfaceGetIterRequest {
	o.QueryPtr = &newValue
	return o
}

func (o *NetInterfaceGetIterRequest) Tag() string {
	r := *o.TagPtr
	return r
}

func (o *NetInterfaceGetIterRequest) SetTag(newValue string) *NetInterfaceGetIterRequest {
	o.TagPtr = &newValue
	return o
}

type NetInterfaceGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result NetInterfaceGetIterResponseResult `xml:"results"`
}

func (o NetInterfaceGetIterResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type NetInterfaceGetIterResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr  string                 `xml:"status,attr"`
	ResultReasonAttr  string                 `xml:"reason,attr"`
	ResultErrnoAttr   string                 `xml:"errno,attr"`
	AttributesListPtr []NetInterfaceInfoType `xml:"attributes-list>net-interface-info"`
	NextTagPtr        *string                `xml:"next-tag"`
	NumRecordsPtr     *int                   `xml:"num-records"`
}

func (o *NetInterfaceGetIterResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewNetInterfaceGetIterResponse() *NetInterfaceGetIterResponse {
	return &NetInterfaceGetIterResponse{}
}

func (o NetInterfaceGetIterResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.AttributesListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "attributes-list", o.AttributesListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("attributes-list: nil\n"))
	}
	if o.NextTagPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "next-tag", *o.NextTagPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("next-tag: nil\n"))
	}
	if o.NumRecordsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "num-records", *o.NumRecordsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("num-records: nil\n"))
	}
	return buffer.String()
}

func (o *NetInterfaceGetIterResponseResult) AttributesList() []NetInterfaceInfoType {
	r := o.AttributesListPtr
	return r
}

func (o *NetInterfaceGetIterResponseResult) SetAttributesList(newValue []NetInterfaceInfoType) *NetInterfaceGetIterResponseResult {
	newSlice := make([]NetInterfaceInfoType, len(newValue))
	copy(newSlice, newValue)
	o.AttributesListPtr = newSlice
	return o
}

func (o *NetInterfaceGetIterResponseResult) NextTag() string {
	r := *o.NextTagPtr
	return r
}

func (o *NetInterfaceGetIterResponseResult) SetNextTag(newValue string) *NetInterfaceGetIterResponseResult {
	o.NextTagPtr = &newValue
	return o
}

func (o *NetInterfaceGetIterResponseResult) NumRecords() int {
	r := *o.NumRecordsPtr
	return r
}

func (o *NetInterfaceGetIterResponseResult) SetNumRecords(newValue int) *NetInterfaceGetIterResponseResult {
	o.NumRecordsPtr = &newValue
	return o
}
