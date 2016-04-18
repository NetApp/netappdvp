// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeSizeRequest struct {
	XMLName xml.Name `xml:"volume-size"`

	NewSizePtr *string `xml:"new-size"`
	VolumePtr  *string `xml:"volume"`
}

func (o *VolumeSizeRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeSizeRequest() *VolumeSizeRequest { return &VolumeSizeRequest{} }

func (r *VolumeSizeRequest) ExecuteUsing(zr *ZapiRunner) (VolumeSizeResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeSizeResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-size result:\n%s", n.Result)

	return n, err
}

func (o VolumeSizeRequest) String() string {
	var buffer bytes.Buffer
	if o.NewSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "new-size", *o.NewSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("new-size: nil\n"))
	}
	if o.VolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume", *o.VolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSizeRequest) NewSize() string {
	r := *o.NewSizePtr
	return r
}

func (o *VolumeSizeRequest) SetNewSize(newValue string) *VolumeSizeRequest {
	o.NewSizePtr = &newValue
	return o
}

func (o *VolumeSizeRequest) Volume() string {
	r := *o.VolumePtr
	return r
}

func (o *VolumeSizeRequest) SetVolume(newValue string) *VolumeSizeRequest {
	o.VolumePtr = &newValue
	return o
}

type VolumeSizeResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeSizeResponseResult `xml:"results"`
}

func (o VolumeSizeResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeSizeResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr         string  `xml:"status,attr"`
	ResultReasonAttr         string  `xml:"reason,attr"`
	ResultErrnoAttr          string  `xml:"errno,attr"`
	IsFixedSizeFlexVolumePtr *bool   `xml:"is-fixed-size-flex-volume"`
	IsReadonlyFlexVolumePtr  *bool   `xml:"is-readonly-flex-volume"`
	IsReplicaFlexVolumePtr   *bool   `xml:"is-replica-flex-volume"`
	VolumeSizePtr            *string `xml:"volume-size"`
}

func (o *VolumeSizeResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSizeResponse() *VolumeSizeResponse { return &VolumeSizeResponse{} }

func (o VolumeSizeResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.IsFixedSizeFlexVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-fixed-size-flex-volume", *o.IsFixedSizeFlexVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-fixed-size-flex-volume: nil\n"))
	}
	if o.IsReadonlyFlexVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-readonly-flex-volume", *o.IsReadonlyFlexVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-readonly-flex-volume: nil\n"))
	}
	if o.IsReplicaFlexVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-replica-flex-volume", *o.IsReplicaFlexVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-replica-flex-volume: nil\n"))
	}
	if o.VolumeSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-size", *o.VolumeSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-size: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSizeResponseResult) IsFixedSizeFlexVolume() bool {
	r := *o.IsFixedSizeFlexVolumePtr
	return r
}

func (o *VolumeSizeResponseResult) SetIsFixedSizeFlexVolume(newValue bool) *VolumeSizeResponseResult {
	o.IsFixedSizeFlexVolumePtr = &newValue
	return o
}

func (o *VolumeSizeResponseResult) IsReadonlyFlexVolume() bool {
	r := *o.IsReadonlyFlexVolumePtr
	return r
}

func (o *VolumeSizeResponseResult) SetIsReadonlyFlexVolume(newValue bool) *VolumeSizeResponseResult {
	o.IsReadonlyFlexVolumePtr = &newValue
	return o
}

func (o *VolumeSizeResponseResult) IsReplicaFlexVolume() bool {
	r := *o.IsReplicaFlexVolumePtr
	return r
}

func (o *VolumeSizeResponseResult) SetIsReplicaFlexVolume(newValue bool) *VolumeSizeResponseResult {
	o.IsReplicaFlexVolumePtr = &newValue
	return o
}

func (o *VolumeSizeResponseResult) VolumeSize() string {
	r := *o.VolumeSizePtr
	return r
}

func (o *VolumeSizeResponseResult) SetVolumeSize(newValue string) *VolumeSizeResponseResult {
	o.VolumeSizePtr = &newValue
	return o
}
