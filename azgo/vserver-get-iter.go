// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VserverGetIterRequest struct {
	XMLName xml.Name `xml:"vserver-get-iter"`

	DesiredAttributesPtr *VserverInfoType `xml:"desired-attributes>vserver-info"`
	MaxRecordsPtr        *int             `xml:"max-records"`
	QueryPtr             *VserverInfoType `xml:"query>vserver-info"`
	TagPtr               *string          `xml:"tag"`
}

func (o *VserverGetIterRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVserverGetIterRequest() *VserverGetIterRequest { return &VserverGetIterRequest{} }

func (r *VserverGetIterRequest) ExecuteUsing(zr *ZapiRunner) (VserverGetIterResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VserverGetIterResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("vserver-get-iter result:\n%s", n.Result)

	return n, err
}

func (o VserverGetIterRequest) String() string {
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

func (o *VserverGetIterRequest) DesiredAttributes() VserverInfoType {
	r := *o.DesiredAttributesPtr
	return r
}

func (o *VserverGetIterRequest) SetDesiredAttributes(newValue VserverInfoType) *VserverGetIterRequest {
	o.DesiredAttributesPtr = &newValue
	return o
}

func (o *VserverGetIterRequest) MaxRecords() int {
	r := *o.MaxRecordsPtr
	return r
}

func (o *VserverGetIterRequest) SetMaxRecords(newValue int) *VserverGetIterRequest {
	o.MaxRecordsPtr = &newValue
	return o
}

func (o *VserverGetIterRequest) Query() VserverInfoType {
	r := *o.QueryPtr
	return r
}

func (o *VserverGetIterRequest) SetQuery(newValue VserverInfoType) *VserverGetIterRequest {
	o.QueryPtr = &newValue
	return o
}

func (o *VserverGetIterRequest) Tag() string {
	r := *o.TagPtr
	return r
}

func (o *VserverGetIterRequest) SetTag(newValue string) *VserverGetIterRequest {
	o.TagPtr = &newValue
	return o
}

type VserverGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VserverGetIterResponseResult `xml:"results"`
}

func (o VserverGetIterResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VserverGetIterResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr  string            `xml:"status,attr"`
	ResultReasonAttr  string            `xml:"reason,attr"`
	ResultErrnoAttr   string            `xml:"errno,attr"`
	AttributesListPtr []VserverInfoType `xml:"attributes-list>vserver-info"`
	NextTagPtr        *string           `xml:"next-tag"`
	NumRecordsPtr     *int              `xml:"num-records"`
}

func (o *VserverGetIterResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	} // TODO: handle better
	return string(output), err
}

func NewVserverGetIterResponse() *VserverGetIterResponse { return &VserverGetIterResponse{} }

func (o VserverGetIterResponseResult) String() string {
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

func (o *VserverGetIterResponseResult) AttributesList() []VserverInfoType {
	r := o.AttributesListPtr
	return r
}

func (o *VserverGetIterResponseResult) SetAttributesList(newValue []VserverInfoType) *VserverGetIterResponseResult {
	newSlice := make([]VserverInfoType, len(newValue))
	copy(newSlice, newValue)
	o.AttributesListPtr = newSlice
	return o
}

func (o *VserverGetIterResponseResult) NextTag() string {
	r := *o.NextTagPtr
	return r
}

func (o *VserverGetIterResponseResult) SetNextTag(newValue string) *VserverGetIterResponseResult {
	o.NextTagPtr = &newValue
	return o
}

func (o *VserverGetIterResponseResult) NumRecords() int {
	r := *o.NumRecordsPtr
	return r
}

func (o *VserverGetIterResponseResult) SetNumRecords(newValue int) *VserverGetIterResponseResult {
	o.NumRecordsPtr = &newValue
	return o
}
