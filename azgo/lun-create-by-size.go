// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type LunCreateBySizeRequest struct {
	XMLName xml.Name `xml:"lun-create-by-size"`

	CachingPolicyPtr           *string `xml:"caching-policy"`
	ClassPtr                   *string `xml:"class"`
	CommentPtr                 *string `xml:"comment"`
	ForeignDiskPtr             *string `xml:"foreign-disk"`
	OstypePtr                  *string `xml:"ostype"`
	PathPtr                    *string `xml:"path"`
	PrefixSizePtr              *int    `xml:"prefix-size"`
	QosPolicyGroupPtr          *string `xml:"qos-policy-group"`
	SizePtr                    *int    `xml:"size"`
	SpaceAllocationEnabledPtr  *bool   `xml:"space-allocation-enabled"`
	SpaceReservationEnabledPtr *bool   `xml:"space-reservation-enabled"`
	TypePtr                    *string `xml:"type"`
}

func (o *LunCreateBySizeRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewLunCreateBySizeRequest() *LunCreateBySizeRequest { return &LunCreateBySizeRequest{} }

func (r *LunCreateBySizeRequest) ExecuteUsing(zr *ZapiRunner) (LunCreateBySizeResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n LunCreateBySizeResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("lun-create-by-size result:\n%s", n.Result)

	return n, err
}

func (o LunCreateBySizeRequest) String() string {
	var buffer bytes.Buffer
	if o.CachingPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "caching-policy", *o.CachingPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("caching-policy: nil\n"))
	}
	if o.ClassPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "class", *o.ClassPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("class: nil\n"))
	}
	if o.CommentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "comment", *o.CommentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("comment: nil\n"))
	}
	if o.ForeignDiskPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "foreign-disk", *o.ForeignDiskPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("foreign-disk: nil\n"))
	}
	if o.OstypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "ostype", *o.OstypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("ostype: nil\n"))
	}
	if o.PathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "path", *o.PathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("path: nil\n"))
	}
	if o.PrefixSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "prefix-size", *o.PrefixSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("prefix-size: nil\n"))
	}
	if o.QosPolicyGroupPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "qos-policy-group", *o.QosPolicyGroupPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("qos-policy-group: nil\n"))
	}
	if o.SizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size", *o.SizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size: nil\n"))
	}
	if o.SpaceAllocationEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-allocation-enabled", *o.SpaceAllocationEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-allocation-enabled: nil\n"))
	}
	if o.SpaceReservationEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-reservation-enabled", *o.SpaceReservationEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-reservation-enabled: nil\n"))
	}
	if o.TypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "type", *o.TypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("type: nil\n"))
	}
	return buffer.String()
}

func (o *LunCreateBySizeRequest) CachingPolicy() string {
	r := *o.CachingPolicyPtr
	return r
}

func (o *LunCreateBySizeRequest) SetCachingPolicy(newValue string) *LunCreateBySizeRequest {
	o.CachingPolicyPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) Class() string {
	r := *o.ClassPtr
	return r
}

func (o *LunCreateBySizeRequest) SetClass(newValue string) *LunCreateBySizeRequest {
	o.ClassPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) Comment() string {
	r := *o.CommentPtr
	return r
}

func (o *LunCreateBySizeRequest) SetComment(newValue string) *LunCreateBySizeRequest {
	o.CommentPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) ForeignDisk() string {
	r := *o.ForeignDiskPtr
	return r
}

func (o *LunCreateBySizeRequest) SetForeignDisk(newValue string) *LunCreateBySizeRequest {
	o.ForeignDiskPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) Ostype() string {
	r := *o.OstypePtr
	return r
}

func (o *LunCreateBySizeRequest) SetOstype(newValue string) *LunCreateBySizeRequest {
	o.OstypePtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) Path() string {
	r := *o.PathPtr
	return r
}

func (o *LunCreateBySizeRequest) SetPath(newValue string) *LunCreateBySizeRequest {
	o.PathPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) PrefixSize() int {
	r := *o.PrefixSizePtr
	return r
}

func (o *LunCreateBySizeRequest) SetPrefixSize(newValue int) *LunCreateBySizeRequest {
	o.PrefixSizePtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) QosPolicyGroup() string {
	r := *o.QosPolicyGroupPtr
	return r
}

func (o *LunCreateBySizeRequest) SetQosPolicyGroup(newValue string) *LunCreateBySizeRequest {
	o.QosPolicyGroupPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) Size() int {
	r := *o.SizePtr
	return r
}

func (o *LunCreateBySizeRequest) SetSize(newValue int) *LunCreateBySizeRequest {
	o.SizePtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) SpaceAllocationEnabled() bool {
	r := *o.SpaceAllocationEnabledPtr
	return r
}

func (o *LunCreateBySizeRequest) SetSpaceAllocationEnabled(newValue bool) *LunCreateBySizeRequest {
	o.SpaceAllocationEnabledPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) SpaceReservationEnabled() bool {
	r := *o.SpaceReservationEnabledPtr
	return r
}

func (o *LunCreateBySizeRequest) SetSpaceReservationEnabled(newValue bool) *LunCreateBySizeRequest {
	o.SpaceReservationEnabledPtr = &newValue
	return o
}

func (o *LunCreateBySizeRequest) Type() string {
	r := *o.TypePtr
	return r
}

func (o *LunCreateBySizeRequest) SetType(newValue string) *LunCreateBySizeRequest {
	o.TypePtr = &newValue
	return o
}

type LunCreateBySizeResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result LunCreateBySizeResponseResult `xml:"results"`
}

func (o LunCreateBySizeResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type LunCreateBySizeResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
	ActualSizePtr    *int   `xml:"actual-size"`
}

func (o *LunCreateBySizeResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewLunCreateBySizeResponse() *LunCreateBySizeResponse { return &LunCreateBySizeResponse{} }

func (o LunCreateBySizeResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	if o.ActualSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "actual-size", *o.ActualSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("actual-size: nil\n"))
	}
	return buffer.String()
}

func (o *LunCreateBySizeResponseResult) ActualSize() int {
	r := *o.ActualSizePtr
	return r
}

func (o *LunCreateBySizeResponseResult) SetActualSize(newValue int) *LunCreateBySizeResponseResult {
	o.ActualSizePtr = &newValue
	return o
}
