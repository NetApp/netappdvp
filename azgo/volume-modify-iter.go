// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeModifyIterRequest struct {
	XMLName xml.Name `xml:"volume-modify-iter"`

	AttributesPtr        *VolumeAttributesType `xml:"attributes>volume-attributes"`
	ContinueOnFailurePtr *bool                 `xml:"continue-on-failure"`
	MaxFailureCountPtr   *int                  `xml:"max-failure-count"`
	MaxRecordsPtr        *int                  `xml:"max-records"`
	QueryPtr             *VolumeAttributesType `xml:"query>volume-attributes"`
	ReturnFailureListPtr *bool                 `xml:"return-failure-list"`
	ReturnSuccessListPtr *bool                 `xml:"return-success-list"`
	TagPtr               *string               `xml:"tag"`
}

func (o *VolumeModifyIterRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeModifyIterRequest() *VolumeModifyIterRequest { return &VolumeModifyIterRequest{} }

func (r *VolumeModifyIterRequest) ExecuteUsing(zr *ZapiRunner) (VolumeModifyIterResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeModifyIterResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-modify-iter result:\n%s", n.Result)

	return n, err
}

func (o VolumeModifyIterRequest) String() string {
	var buffer bytes.Buffer
	if o.AttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "attributes", *o.AttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("attributes: nil\n"))
	}
	if o.ContinueOnFailurePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "continue-on-failure", *o.ContinueOnFailurePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("continue-on-failure: nil\n"))
	}
	if o.MaxFailureCountPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-failure-count", *o.MaxFailureCountPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-failure-count: nil\n"))
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
	if o.ReturnFailureListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "return-failure-list", *o.ReturnFailureListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("return-failure-list: nil\n"))
	}
	if o.ReturnSuccessListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "return-success-list", *o.ReturnSuccessListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("return-success-list: nil\n"))
	}
	if o.TagPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "tag", *o.TagPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("tag: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeModifyIterRequest) Attributes() VolumeAttributesType {
	r := *o.AttributesPtr
	return r
}

func (o *VolumeModifyIterRequest) SetAttributes(newValue VolumeAttributesType) *VolumeModifyIterRequest {
	o.AttributesPtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) ContinueOnFailure() bool {
	r := *o.ContinueOnFailurePtr
	return r
}

func (o *VolumeModifyIterRequest) SetContinueOnFailure(newValue bool) *VolumeModifyIterRequest {
	o.ContinueOnFailurePtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) MaxFailureCount() int {
	r := *o.MaxFailureCountPtr
	return r
}

func (o *VolumeModifyIterRequest) SetMaxFailureCount(newValue int) *VolumeModifyIterRequest {
	o.MaxFailureCountPtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) MaxRecords() int {
	r := *o.MaxRecordsPtr
	return r
}

func (o *VolumeModifyIterRequest) SetMaxRecords(newValue int) *VolumeModifyIterRequest {
	o.MaxRecordsPtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) Query() VolumeAttributesType {
	r := *o.QueryPtr
	return r
}

func (o *VolumeModifyIterRequest) SetQuery(newValue VolumeAttributesType) *VolumeModifyIterRequest {
	o.QueryPtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) ReturnFailureList() bool {
	r := *o.ReturnFailureListPtr
	return r
}

func (o *VolumeModifyIterRequest) SetReturnFailureList(newValue bool) *VolumeModifyIterRequest {
	o.ReturnFailureListPtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) ReturnSuccessList() bool {
	r := *o.ReturnSuccessListPtr
	return r
}

func (o *VolumeModifyIterRequest) SetReturnSuccessList(newValue bool) *VolumeModifyIterRequest {
	o.ReturnSuccessListPtr = &newValue
	return o
}

func (o *VolumeModifyIterRequest) Tag() string {
	r := *o.TagPtr
	return r
}

func (o *VolumeModifyIterRequest) SetTag(newValue string) *VolumeModifyIterRequest {
	o.TagPtr = &newValue
	return o
}

type VolumeModifyIterResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeModifyIterResponseResult `xml:"results"`
}

func (o VolumeModifyIterResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeModifyIterResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string                     `xml:"status,attr"`
	ResultReasonAttr string                     `xml:"reason,attr"`
	ResultErrnoAttr  string                     `xml:"errno,attr"`
	FailureListPtr   []VolumeModifyIterInfoType `xml:"failure-list>volume-modify-iter-info"`
	NextTagPtr       *string                    `xml:"next-tag"`
	NumFailedPtr     *int                       `xml:"num-failed"`
	NumSucceededPtr  *int                       `xml:"num-succeeded"`
	SuccessListPtr   []VolumeModifyIterInfoType `xml:"success-list>volume-modify-iter-info"`
}

func (o *VolumeModifyIterResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeModifyIterResponse() *VolumeModifyIterResponse { return &VolumeModifyIterResponse{} }

func (o VolumeModifyIterResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.FailureListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "failure-list", o.FailureListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("failure-list: nil\n"))
	}
	if o.NextTagPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "next-tag", *o.NextTagPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("next-tag: nil\n"))
	}
	if o.NumFailedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "num-failed", *o.NumFailedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("num-failed: nil\n"))
	}
	if o.NumSucceededPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "num-succeeded", *o.NumSucceededPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("num-succeeded: nil\n"))
	}
	if o.SuccessListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "success-list", o.SuccessListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("success-list: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeModifyIterResponseResult) FailureList() []VolumeModifyIterInfoType {
	r := o.FailureListPtr
	return r
}

func (o *VolumeModifyIterResponseResult) SetFailureList(newValue []VolumeModifyIterInfoType) *VolumeModifyIterResponseResult {
	newSlice := make([]VolumeModifyIterInfoType, len(newValue))
	copy(newSlice, newValue)
	o.FailureListPtr = newSlice
	return o
}

func (o *VolumeModifyIterResponseResult) NextTag() string {
	r := *o.NextTagPtr
	return r
}

func (o *VolumeModifyIterResponseResult) SetNextTag(newValue string) *VolumeModifyIterResponseResult {
	o.NextTagPtr = &newValue
	return o
}

func (o *VolumeModifyIterResponseResult) NumFailed() int {
	r := *o.NumFailedPtr
	return r
}

func (o *VolumeModifyIterResponseResult) SetNumFailed(newValue int) *VolumeModifyIterResponseResult {
	o.NumFailedPtr = &newValue
	return o
}

func (o *VolumeModifyIterResponseResult) NumSucceeded() int {
	r := *o.NumSucceededPtr
	return r
}

func (o *VolumeModifyIterResponseResult) SetNumSucceeded(newValue int) *VolumeModifyIterResponseResult {
	o.NumSucceededPtr = &newValue
	return o
}

func (o *VolumeModifyIterResponseResult) SuccessList() []VolumeModifyIterInfoType {
	r := o.SuccessListPtr
	return r
}

func (o *VolumeModifyIterResponseResult) SetSuccessList(newValue []VolumeModifyIterInfoType) *VolumeModifyIterResponseResult {
	newSlice := make([]VolumeModifyIterInfoType, len(newValue))
	copy(newSlice, newValue)
	o.SuccessListPtr = newSlice
	return o
}
