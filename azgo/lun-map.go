// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type LunMapRequest struct {
	XMLName xml.Name `xml:"lun-map"`

	AdditionalReportingNodePtr *NodeNameType `xml:"additional-reporting-node>node-name"`
	ForcePtr                   *bool         `xml:"force"`
	InitiatorGroupPtr          *string       `xml:"initiator-group"`
	LunIdPtr                   *int          `xml:"lun-id"`
	PathPtr                    *string       `xml:"path"`
}

func (o *LunMapRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewLunMapRequest() *LunMapRequest { return &LunMapRequest{} }

func (r *LunMapRequest) ExecuteUsing(zr *ZapiRunner) (LunMapResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n LunMapResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("lun-map result:\n%s", n.Result)

	return n, err
}

func (o LunMapRequest) String() string {
	var buffer bytes.Buffer
	if o.AdditionalReportingNodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "additional-reporting-node", *o.AdditionalReportingNodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("additional-reporting-node: nil\n"))
	}
	if o.ForcePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force", *o.ForcePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force: nil\n"))
	}
	if o.InitiatorGroupPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group", *o.InitiatorGroupPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group: nil\n"))
	}
	if o.LunIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "lun-id", *o.LunIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("lun-id: nil\n"))
	}
	if o.PathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "path", *o.PathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("path: nil\n"))
	}
	return buffer.String()
}

func (o *LunMapRequest) AdditionalReportingNode() NodeNameType {
	r := *o.AdditionalReportingNodePtr
	return r
}

func (o *LunMapRequest) SetAdditionalReportingNode(newValue NodeNameType) *LunMapRequest {
	o.AdditionalReportingNodePtr = &newValue
	return o
}

func (o *LunMapRequest) Force() bool {
	r := *o.ForcePtr
	return r
}

func (o *LunMapRequest) SetForce(newValue bool) *LunMapRequest {
	o.ForcePtr = &newValue
	return o
}

func (o *LunMapRequest) InitiatorGroup() string {
	r := *o.InitiatorGroupPtr
	return r
}

func (o *LunMapRequest) SetInitiatorGroup(newValue string) *LunMapRequest {
	o.InitiatorGroupPtr = &newValue
	return o
}

func (o *LunMapRequest) LunId() int {
	r := *o.LunIdPtr
	return r
}

func (o *LunMapRequest) SetLunId(newValue int) *LunMapRequest {
	o.LunIdPtr = &newValue
	return o
}

func (o *LunMapRequest) Path() string {
	r := *o.PathPtr
	return r
}

func (o *LunMapRequest) SetPath(newValue string) *LunMapRequest {
	o.PathPtr = &newValue
	return o
}

type LunMapResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result LunMapResponseResult `xml:"results"`
}

func (o LunMapResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type LunMapResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
	LunIdAssignedPtr *int   `xml:"lun-id-assigned"`
}

func (o *LunMapResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewLunMapResponse() *LunMapResponse { return &LunMapResponse{} }

func (o LunMapResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.LunIdAssignedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "lun-id-assigned", *o.LunIdAssignedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("lun-id-assigned: nil\n"))
	}
	return buffer.String()
}

func (o *LunMapResponseResult) LunIdAssigned() int {
	r := *o.LunIdAssignedPtr
	return r
}

func (o *LunMapResponseResult) SetLunIdAssigned(newValue int) *LunMapResponseResult {
	o.LunIdAssignedPtr = &newValue
	return o
}
