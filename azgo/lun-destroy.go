// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type LunDestroyRequest struct {
	XMLName xml.Name `xml:"lun-destroy"`

	DestroyFencedLunPtr *bool   `xml:"destroy-fenced-lun"`
	ForcePtr            *bool   `xml:"force"`
	PathPtr             *string `xml:"path"`
}

func (o *LunDestroyRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewLunDestroyRequest() *LunDestroyRequest { return &LunDestroyRequest{} }

func (r *LunDestroyRequest) ExecuteUsing(zr *ZapiRunner) (LunDestroyResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n LunDestroyResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("lun-destroy result:\n%s", n.Result)

	return n, err
}

func (o LunDestroyRequest) String() string {
	var buffer bytes.Buffer
	if o.DestroyFencedLunPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "destroy-fenced-lun", *o.DestroyFencedLunPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("destroy-fenced-lun: nil\n"))
	}
	if o.ForcePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force", *o.ForcePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force: nil\n"))
	}
	if o.PathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "path", *o.PathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("path: nil\n"))
	}
	return buffer.String()
}

func (o *LunDestroyRequest) DestroyFencedLun() bool {
	r := *o.DestroyFencedLunPtr
	return r
}

func (o *LunDestroyRequest) SetDestroyFencedLun(newValue bool) *LunDestroyRequest {
	o.DestroyFencedLunPtr = &newValue
	return o
}

func (o *LunDestroyRequest) Force() bool {
	r := *o.ForcePtr
	return r
}

func (o *LunDestroyRequest) SetForce(newValue bool) *LunDestroyRequest {
	o.ForcePtr = &newValue
	return o
}

func (o *LunDestroyRequest) Path() string {
	r := *o.PathPtr
	return r
}

func (o *LunDestroyRequest) SetPath(newValue string) *LunDestroyRequest {
	o.PathPtr = &newValue
	return o
}

type LunDestroyResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result LunDestroyResponseResult `xml:"results"`
}

func (o LunDestroyResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type LunDestroyResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *LunDestroyResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewLunDestroyResponse() *LunDestroyResponse { return &LunDestroyResponse{} }

func (o LunDestroyResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
