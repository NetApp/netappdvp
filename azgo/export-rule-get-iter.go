// Copyright 2017 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

// ExportRuleGetIterRequest is a structure to represent a export-rule-get-iter ZAPI request object
type ExportRuleGetIterRequest struct {
	XMLName xml.Name `xml:"export-rule-get-iter"`

	DesiredAttributesPtr *ExportRuleInfoType `xml:"desired-attributes>export-rule-info"`
	MaxRecordsPtr        *int                `xml:"max-records"`
	QueryPtr             *ExportRuleInfoType `xml:"query>export-rule-info"`
	TagPtr               *string             `xml:"tag"`
}

// ToXML converts this object into an xml string representation
func (o *ExportRuleGetIterRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	//if err != nil { log.Errorf("error: %v\n", err) }
	return string(output), err
}

// NewExportRuleGetIterRequest is a factory method for creating new instances of ExportRuleGetIterRequest objects
func NewExportRuleGetIterRequest() *ExportRuleGetIterRequest { return &ExportRuleGetIterRequest{} }

// ExecuteUsing converts this object to a ZAPI XML representation and uses the supplied ZapiRunner to send to a filer
func (o *ExportRuleGetIterRequest) ExecuteUsing(zr *ZapiRunner) (ExportRuleGetIterResponse, error) {

	if zr.DebugTraceFlags["method"] {
		fields := log.Fields{"Method": "ExecuteUsing", "Type": "ExportRuleGetIterRequest"}
		log.WithFields(fields).Debug(">>>> ExecuteUsing")
		defer log.WithFields(fields).Debug("<<<< ExecuteUsing")
	}

	combined := NewExportRuleGetIterResponse()
	var nextTagPtr *string
	done := false
	for done != true {

		resp, err := zr.SendZapi(o)
		if err != nil {
			log.Errorf("API invocation failed. %v", err.Error())
			return *combined, err
		}
		defer resp.Body.Close()
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Errorf("Error reading response body. %v", readErr.Error())
			return *combined, readErr
		}
		if zr.DebugTraceFlags["api"] {
			log.Debugf("response Body:\n%s", string(body))
		}

		var n ExportRuleGetIterResponse
		unmarshalErr := xml.Unmarshal(body, &n)
		if unmarshalErr != nil {
			log.WithField("body", string(body)).Warnf("Error unmarshaling response body. %v", unmarshalErr.Error())
			//return *combined, unmarshalErr
		}
		if zr.DebugTraceFlags["api"] {
			log.Debugf("export-rule-get-iter result:\n%s", n.Result)
		}

		if err == nil {
			nextTagPtr = n.Result.NextTagPtr
			if nextTagPtr == nil {
				done = true
			} else {
				o.SetTag(*nextTagPtr)
			}

			recordsRead := n.Result.NumRecords()
			if recordsRead == 0 {
				done = true
			}

			combined.Result.SetAttributesList(append(combined.Result.AttributesList(), n.Result.AttributesList()...))
			if done == true {
				combined.Result.ResultErrnoAttr = n.Result.ResultErrnoAttr
				combined.Result.ResultReasonAttr = n.Result.ResultReasonAttr
				combined.Result.ResultStatusAttr = n.Result.ResultStatusAttr
				combined.Result.SetNumRecords(len(combined.Result.AttributesList()))
			}
		}
	}

	return *combined, nil
}

// String returns a string representation of this object's fields and implements the Stringer interface
func (o ExportRuleGetIterRequest) String() string {
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

// DesiredAttributes is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterRequest) DesiredAttributes() ExportRuleInfoType {
	r := *o.DesiredAttributesPtr
	return r
}

// SetDesiredAttributes is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterRequest) SetDesiredAttributes(newValue ExportRuleInfoType) *ExportRuleGetIterRequest {
	o.DesiredAttributesPtr = &newValue
	return o
}

// MaxRecords is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterRequest) MaxRecords() int {
	r := *o.MaxRecordsPtr
	return r
}

// SetMaxRecords is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterRequest) SetMaxRecords(newValue int) *ExportRuleGetIterRequest {
	o.MaxRecordsPtr = &newValue
	return o
}

// Query is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterRequest) Query() ExportRuleInfoType {
	r := *o.QueryPtr
	return r
}

// SetQuery is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterRequest) SetQuery(newValue ExportRuleInfoType) *ExportRuleGetIterRequest {
	o.QueryPtr = &newValue
	return o
}

// Tag is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterRequest) Tag() string {
	r := *o.TagPtr
	return r
}

// SetTag is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterRequest) SetTag(newValue string) *ExportRuleGetIterRequest {
	o.TagPtr = &newValue
	return o
}

// ExportRuleGetIterResponse is a structure to represent a export-rule-get-iter ZAPI response object
type ExportRuleGetIterResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result ExportRuleGetIterResponseResult `xml:"results"`
}

// String returns a string representation of this object's fields and implements the Stringer interface
func (o ExportRuleGetIterResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

// ExportRuleGetIterResponseResult is a structure to represent a export-rule-get-iter ZAPI object's result
type ExportRuleGetIterResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr  string               `xml:"status,attr"`
	ResultReasonAttr  string               `xml:"reason,attr"`
	ResultErrnoAttr   string               `xml:"errno,attr"`
	AttributesListPtr []ExportRuleInfoType `xml:"attributes-list>export-rule-info"`
	NextTagPtr        *string              `xml:"next-tag"`
	NumRecordsPtr     *int                 `xml:"num-records"`
}

// ToXML converts this object into an xml string representation
func (o *ExportRuleGetIterResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	//if err != nil { log.Debugf("error: %v", err) }
	return string(output), err
}

// NewExportRuleGetIterResponse is a factory method for creating new instances of ExportRuleGetIterResponse objects
func NewExportRuleGetIterResponse() *ExportRuleGetIterResponse { return &ExportRuleGetIterResponse{} }

// String returns a string representation of this object's fields and implements the Stringer interface
func (o ExportRuleGetIterResponseResult) String() string {
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

// AttributesList is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterResponseResult) AttributesList() []ExportRuleInfoType {
	r := o.AttributesListPtr
	return r
}

// SetAttributesList is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterResponseResult) SetAttributesList(newValue []ExportRuleInfoType) *ExportRuleGetIterResponseResult {
	newSlice := make([]ExportRuleInfoType, len(newValue))
	copy(newSlice, newValue)
	o.AttributesListPtr = newSlice
	return o
}

// NextTag is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterResponseResult) NextTag() string {
	r := *o.NextTagPtr
	return r
}

// SetNextTag is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterResponseResult) SetNextTag(newValue string) *ExportRuleGetIterResponseResult {
	o.NextTagPtr = &newValue
	return o
}

// NumRecords is a fluent style 'getter' method that can be chained
func (o *ExportRuleGetIterResponseResult) NumRecords() int {
	r := *o.NumRecordsPtr
	return r
}

// SetNumRecords is a fluent style 'setter' method that can be chained
func (o *ExportRuleGetIterResponseResult) SetNumRecords(newValue int) *ExportRuleGetIterResponseResult {
	o.NumRecordsPtr = &newValue
	return o
}
