// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeOfflineRequest struct {
	XMLName xml.Name `xml:"volume-offline"`

	NamePtr *string `xml:"name"`
}

func (o *VolumeOfflineRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeOfflineRequest() *VolumeOfflineRequest { return &VolumeOfflineRequest{} }

func (r *VolumeOfflineRequest) ExecuteUsing(zr *ZapiRunner) (VolumeOfflineResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeOfflineResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-offline result:\n%s", n.Result)

	return n, err
}

func (o VolumeOfflineRequest) String() string {
	var buffer bytes.Buffer
	if o.NamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name", *o.NamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeOfflineRequest) Name() string {
	r := *o.NamePtr
	return r
}

func (o *VolumeOfflineRequest) SetName(newValue string) *VolumeOfflineRequest {
	o.NamePtr = &newValue
	return o
}

type VolumeOfflineResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeOfflineResponseResult `xml:"results"`
}

func (o VolumeOfflineResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeOfflineResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *VolumeOfflineResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeOfflineResponse() *VolumeOfflineResponse { return &VolumeOfflineResponse{} }

func (o VolumeOfflineResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
