// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type LunGetSerialNumberRequest struct {
	XMLName xml.Name `xml:"lun-get-serial-number"`

	PathPtr *string `xml:"path"`
}

func (o *LunGetSerialNumberRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewLunGetSerialNumberRequest() *LunGetSerialNumberRequest { return &LunGetSerialNumberRequest{} }

func (r *LunGetSerialNumberRequest) ExecuteUsing(zr *ZapiRunner) (LunGetSerialNumberResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n LunGetSerialNumberResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("lun-get-serial-number result:\n%s", n.Result)

	return n, err
}

func (o LunGetSerialNumberRequest) String() string {
	var buffer bytes.Buffer
	if o.PathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "path", *o.PathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("path: nil\n"))
	}
	return buffer.String()
}

func (o *LunGetSerialNumberRequest) Path() string {
	r := *o.PathPtr
	return r
}

func (o *LunGetSerialNumberRequest) SetPath(newValue string) *LunGetSerialNumberRequest {
	o.PathPtr = &newValue
	return o
}

type LunGetSerialNumberResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result LunGetSerialNumberResponseResult `xml:"results"`
}

func (o LunGetSerialNumberResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type LunGetSerialNumberResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string  `xml:"status,attr"`
	ResultReasonAttr string  `xml:"reason,attr"`
	ResultErrnoAttr  string  `xml:"errno,attr"`
	SerialNumberPtr  *string `xml:"serial-number"`
}

func (o *LunGetSerialNumberResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewLunGetSerialNumberResponse() *LunGetSerialNumberResponse { return &LunGetSerialNumberResponse{} }

func (o LunGetSerialNumberResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.SerialNumberPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "serial-number", *o.SerialNumberPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("serial-number: nil\n"))
	}
	return buffer.String()
}

func (o *LunGetSerialNumberResponseResult) SerialNumber() string {
	r := *o.SerialNumberPtr
	return r
}

func (o *LunGetSerialNumberResponseResult) SetSerialNumber(newValue string) *LunGetSerialNumberResponseResult {
	o.SerialNumberPtr = &newValue
	return o
}
