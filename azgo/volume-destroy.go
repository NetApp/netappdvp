// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeDestroyRequest struct {
	XMLName xml.Name `xml:"volume-destroy"`

	NamePtr              *string `xml:"name"`
	UnmountAndOfflinePtr *bool   `xml:"unmount-and-offline"`
}

func (o *VolumeDestroyRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeDestroyRequest() *VolumeDestroyRequest { return &VolumeDestroyRequest{} }

func (r *VolumeDestroyRequest) ExecuteUsing(zr *ZapiRunner) (VolumeDestroyResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeDestroyResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-destroy result:\n%s", n.Result)

	return n, err
}

func (o VolumeDestroyRequest) String() string {
	var buffer bytes.Buffer
	if o.NamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name", *o.NamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name: nil\n"))
	}
	if o.UnmountAndOfflinePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "unmount-and-offline", *o.UnmountAndOfflinePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("unmount-and-offline: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeDestroyRequest) Name() string {
	r := *o.NamePtr
	return r
}

func (o *VolumeDestroyRequest) SetName(newValue string) *VolumeDestroyRequest {
	o.NamePtr = &newValue
	return o
}

func (o *VolumeDestroyRequest) UnmountAndOffline() bool {
	r := *o.UnmountAndOfflinePtr
	return r
}

func (o *VolumeDestroyRequest) SetUnmountAndOffline(newValue bool) *VolumeDestroyRequest {
	o.UnmountAndOfflinePtr = &newValue
	return o
}

type VolumeDestroyResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeDestroyResponseResult `xml:"results"`
}

func (o VolumeDestroyResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeDestroyResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *VolumeDestroyResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeDestroyResponse() *VolumeDestroyResponse { return &VolumeDestroyResponse{} }

func (o VolumeDestroyResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
