// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type VolumeCreateRequest struct {
	XMLName xml.Name `xml:"volume-create"`

	AntivirusOnAccessPolicyPtr      *string `xml:"antivirus-on-access-policy"`
	CachingPolicyPtr                *string `xml:"caching-policy"`
	ConstituentRolePtr              *string `xml:"constituent-role"`
	ContainingAggrNamePtr           *string `xml:"containing-aggr-name"`
	ExcludedFromAutobalancePtr      *bool   `xml:"excluded-from-autobalance"`
	ExportPolicyPtr                 *string `xml:"export-policy"`
	FlexcacheCachePolicyPtr         *string `xml:"flexcache-cache-policy"`
	FlexcacheFillPolicyPtr          *string `xml:"flexcache-fill-policy"`
	FlexcacheOriginVolumeNamePtr    *string `xml:"flexcache-origin-volume-name"`
	GroupIdPtr                      *int    `xml:"group-id"`
	IsJunctionActivePtr             *bool   `xml:"is-junction-active"`
	IsNvfailEnabledPtr              *string `xml:"is-nvfail-enabled"`
	IsVserverRootPtr                *bool   `xml:"is-vserver-root"`
	JunctionPathPtr                 *string `xml:"junction-path"`
	LanguageCodePtr                 *string `xml:"language-code"`
	MaxDirSizePtr                   *int    `xml:"max-dir-size"`
	MaxWriteAllocBlocksPtr          *int    `xml:"max-write-alloc-blocks"`
	PercentageSnapshotReservePtr    *int    `xml:"percentage-snapshot-reserve"`
	QosPolicyGroupNamePtr           *string `xml:"qos-policy-group-name"`
	SizePtr                         *string `xml:"size"`
	SnapshotPolicyPtr               *string `xml:"snapshot-policy"`
	SpaceReservePtr                 *string `xml:"space-reserve"`
	StorageServicePtr               *string `xml:"storage-service"`
	StripeAlgorithmPtr              *string `xml:"stripe-algorithm"`
	StripeConcurrencyPtr            *string `xml:"stripe-concurrency"`
	StripeConstituentVolumeCountPtr *int    `xml:"stripe-constituent-volume-count"`
	StripeOptimizePtr               *string `xml:"stripe-optimize"`
	StripeWidthPtr                  *int    `xml:"stripe-width"`
	UnixPermissionsPtr              *string `xml:"unix-permissions"`
	UserIdPtr                       *int    `xml:"user-id"`
	VmAlignSectorPtr                *int    `xml:"vm-align-sector"`
	VmAlignSuffixPtr                *string `xml:"vm-align-suffix"`
	VolumePtr                       *string `xml:"volume"`
	VolumeCommentPtr                *string `xml:"volume-comment"`
	VolumeSecurityStylePtr          *string `xml:"volume-security-style"`
	VolumeStatePtr                  *string `xml:"volume-state"`
	VolumeTypePtr                   *string `xml:"volume-type"`
}

func (o *VolumeCreateRequest) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v\n", err)
	}
	return string(output), err
}

func NewVolumeCreateRequest() *VolumeCreateRequest { return &VolumeCreateRequest{} }

func (r *VolumeCreateRequest) ExecuteUsing(zr *ZapiRunner) (VolumeCreateResponse, error) {
	resp, err := zr.SendZapi(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	var n VolumeCreateResponse
	xml.Unmarshal(body, &n)
	if err != nil {
		log.Errorf("err: %v", err.Error())
	}
	log.Debugf("volume-create result:\n%s", n.Result)

	return n, err
}

func (o VolumeCreateRequest) String() string {
	var buffer bytes.Buffer
	if o.AntivirusOnAccessPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "antivirus-on-access-policy", *o.AntivirusOnAccessPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("antivirus-on-access-policy: nil\n"))
	}
	if o.CachingPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "caching-policy", *o.CachingPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("caching-policy: nil\n"))
	}
	if o.ConstituentRolePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "constituent-role", *o.ConstituentRolePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("constituent-role: nil\n"))
	}
	if o.ContainingAggrNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "containing-aggr-name", *o.ContainingAggrNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("containing-aggr-name: nil\n"))
	}
	if o.ExcludedFromAutobalancePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "excluded-from-autobalance", *o.ExcludedFromAutobalancePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("excluded-from-autobalance: nil\n"))
	}
	if o.ExportPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "export-policy", *o.ExportPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("export-policy: nil\n"))
	}
	if o.FlexcacheCachePolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "flexcache-cache-policy", *o.FlexcacheCachePolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("flexcache-cache-policy: nil\n"))
	}
	if o.FlexcacheFillPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "flexcache-fill-policy", *o.FlexcacheFillPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("flexcache-fill-policy: nil\n"))
	}
	if o.FlexcacheOriginVolumeNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "flexcache-origin-volume-name", *o.FlexcacheOriginVolumeNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("flexcache-origin-volume-name: nil\n"))
	}
	if o.GroupIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "group-id", *o.GroupIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("group-id: nil\n"))
	}
	if o.IsJunctionActivePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-junction-active", *o.IsJunctionActivePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-junction-active: nil\n"))
	}
	if o.IsNvfailEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-nvfail-enabled", *o.IsNvfailEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-nvfail-enabled: nil\n"))
	}
	if o.IsVserverRootPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-vserver-root", *o.IsVserverRootPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-vserver-root: nil\n"))
	}
	if o.JunctionPathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "junction-path", *o.JunctionPathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("junction-path: nil\n"))
	}
	if o.LanguageCodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "language-code", *o.LanguageCodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("language-code: nil\n"))
	}
	if o.MaxDirSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-dir-size", *o.MaxDirSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-dir-size: nil\n"))
	}
	if o.MaxWriteAllocBlocksPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-write-alloc-blocks", *o.MaxWriteAllocBlocksPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-write-alloc-blocks: nil\n"))
	}
	if o.PercentageSnapshotReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-snapshot-reserve", *o.PercentageSnapshotReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-snapshot-reserve: nil\n"))
	}
	if o.QosPolicyGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "qos-policy-group-name", *o.QosPolicyGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("qos-policy-group-name: nil\n"))
	}
	if o.SizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size", *o.SizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size: nil\n"))
	}
	if o.SnapshotPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-policy", *o.SnapshotPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-policy: nil\n"))
	}
	if o.SpaceReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-reserve", *o.SpaceReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-reserve: nil\n"))
	}
	if o.StorageServicePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "storage-service", *o.StorageServicePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("storage-service: nil\n"))
	}
	if o.StripeAlgorithmPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "stripe-algorithm", *o.StripeAlgorithmPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("stripe-algorithm: nil\n"))
	}
	if o.StripeConcurrencyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "stripe-concurrency", *o.StripeConcurrencyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("stripe-concurrency: nil\n"))
	}
	if o.StripeConstituentVolumeCountPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "stripe-constituent-volume-count", *o.StripeConstituentVolumeCountPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("stripe-constituent-volume-count: nil\n"))
	}
	if o.StripeOptimizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "stripe-optimize", *o.StripeOptimizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("stripe-optimize: nil\n"))
	}
	if o.StripeWidthPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "stripe-width", *o.StripeWidthPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("stripe-width: nil\n"))
	}
	if o.UnixPermissionsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "unix-permissions", *o.UnixPermissionsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("unix-permissions: nil\n"))
	}
	if o.UserIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "user-id", *o.UserIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("user-id: nil\n"))
	}
	if o.VmAlignSectorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vm-align-sector", *o.VmAlignSectorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vm-align-sector: nil\n"))
	}
	if o.VmAlignSuffixPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vm-align-suffix", *o.VmAlignSuffixPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vm-align-suffix: nil\n"))
	}
	if o.VolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume", *o.VolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume: nil\n"))
	}
	if o.VolumeCommentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-comment", *o.VolumeCommentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-comment: nil\n"))
	}
	if o.VolumeSecurityStylePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-security-style", *o.VolumeSecurityStylePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-security-style: nil\n"))
	}
	if o.VolumeStatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-state", *o.VolumeStatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-state: nil\n"))
	}
	if o.VolumeTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-type", *o.VolumeTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-type: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeCreateRequest) AntivirusOnAccessPolicy() string {
	r := *o.AntivirusOnAccessPolicyPtr
	return r
}

func (o *VolumeCreateRequest) SetAntivirusOnAccessPolicy(newValue string) *VolumeCreateRequest {
	o.AntivirusOnAccessPolicyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) CachingPolicy() string {
	r := *o.CachingPolicyPtr
	return r
}

func (o *VolumeCreateRequest) SetCachingPolicy(newValue string) *VolumeCreateRequest {
	o.CachingPolicyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) ConstituentRole() string {
	r := *o.ConstituentRolePtr
	return r
}

func (o *VolumeCreateRequest) SetConstituentRole(newValue string) *VolumeCreateRequest {
	o.ConstituentRolePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) ContainingAggrName() string {
	r := *o.ContainingAggrNamePtr
	return r
}

func (o *VolumeCreateRequest) SetContainingAggrName(newValue string) *VolumeCreateRequest {
	o.ContainingAggrNamePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) ExcludedFromAutobalance() bool {
	r := *o.ExcludedFromAutobalancePtr
	return r
}

func (o *VolumeCreateRequest) SetExcludedFromAutobalance(newValue bool) *VolumeCreateRequest {
	o.ExcludedFromAutobalancePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) ExportPolicy() string {
	r := *o.ExportPolicyPtr
	return r
}

func (o *VolumeCreateRequest) SetExportPolicy(newValue string) *VolumeCreateRequest {
	o.ExportPolicyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) FlexcacheCachePolicy() string {
	r := *o.FlexcacheCachePolicyPtr
	return r
}

func (o *VolumeCreateRequest) SetFlexcacheCachePolicy(newValue string) *VolumeCreateRequest {
	o.FlexcacheCachePolicyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) FlexcacheFillPolicy() string {
	r := *o.FlexcacheFillPolicyPtr
	return r
}

func (o *VolumeCreateRequest) SetFlexcacheFillPolicy(newValue string) *VolumeCreateRequest {
	o.FlexcacheFillPolicyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) FlexcacheOriginVolumeName() string {
	r := *o.FlexcacheOriginVolumeNamePtr
	return r
}

func (o *VolumeCreateRequest) SetFlexcacheOriginVolumeName(newValue string) *VolumeCreateRequest {
	o.FlexcacheOriginVolumeNamePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) GroupId() int {
	r := *o.GroupIdPtr
	return r
}

func (o *VolumeCreateRequest) SetGroupId(newValue int) *VolumeCreateRequest {
	o.GroupIdPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) IsJunctionActive() bool {
	r := *o.IsJunctionActivePtr
	return r
}

func (o *VolumeCreateRequest) SetIsJunctionActive(newValue bool) *VolumeCreateRequest {
	o.IsJunctionActivePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) IsNvfailEnabled() string {
	r := *o.IsNvfailEnabledPtr
	return r
}

func (o *VolumeCreateRequest) SetIsNvfailEnabled(newValue string) *VolumeCreateRequest {
	o.IsNvfailEnabledPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) IsVserverRoot() bool {
	r := *o.IsVserverRootPtr
	return r
}

func (o *VolumeCreateRequest) SetIsVserverRoot(newValue bool) *VolumeCreateRequest {
	o.IsVserverRootPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) JunctionPath() string {
	r := *o.JunctionPathPtr
	return r
}

func (o *VolumeCreateRequest) SetJunctionPath(newValue string) *VolumeCreateRequest {
	o.JunctionPathPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) LanguageCode() string {
	r := *o.LanguageCodePtr
	return r
}

func (o *VolumeCreateRequest) SetLanguageCode(newValue string) *VolumeCreateRequest {
	o.LanguageCodePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) MaxDirSize() int {
	r := *o.MaxDirSizePtr
	return r
}

func (o *VolumeCreateRequest) SetMaxDirSize(newValue int) *VolumeCreateRequest {
	o.MaxDirSizePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) MaxWriteAllocBlocks() int {
	r := *o.MaxWriteAllocBlocksPtr
	return r
}

func (o *VolumeCreateRequest) SetMaxWriteAllocBlocks(newValue int) *VolumeCreateRequest {
	o.MaxWriteAllocBlocksPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) PercentageSnapshotReserve() int {
	r := *o.PercentageSnapshotReservePtr
	return r
}

func (o *VolumeCreateRequest) SetPercentageSnapshotReserve(newValue int) *VolumeCreateRequest {
	o.PercentageSnapshotReservePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) QosPolicyGroupName() string {
	r := *o.QosPolicyGroupNamePtr
	return r
}

func (o *VolumeCreateRequest) SetQosPolicyGroupName(newValue string) *VolumeCreateRequest {
	o.QosPolicyGroupNamePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) Size() string {
	r := *o.SizePtr
	return r
}

func (o *VolumeCreateRequest) SetSize(newValue string) *VolumeCreateRequest {
	o.SizePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) SnapshotPolicy() string {
	r := *o.SnapshotPolicyPtr
	return r
}

func (o *VolumeCreateRequest) SetSnapshotPolicy(newValue string) *VolumeCreateRequest {
	o.SnapshotPolicyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) SpaceReserve() string {
	r := *o.SpaceReservePtr
	return r
}

func (o *VolumeCreateRequest) SetSpaceReserve(newValue string) *VolumeCreateRequest {
	o.SpaceReservePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) StorageService() string {
	r := *o.StorageServicePtr
	return r
}

func (o *VolumeCreateRequest) SetStorageService(newValue string) *VolumeCreateRequest {
	o.StorageServicePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) StripeAlgorithm() string {
	r := *o.StripeAlgorithmPtr
	return r
}

func (o *VolumeCreateRequest) SetStripeAlgorithm(newValue string) *VolumeCreateRequest {
	o.StripeAlgorithmPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) StripeConcurrency() string {
	r := *o.StripeConcurrencyPtr
	return r
}

func (o *VolumeCreateRequest) SetStripeConcurrency(newValue string) *VolumeCreateRequest {
	o.StripeConcurrencyPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) StripeConstituentVolumeCount() int {
	r := *o.StripeConstituentVolumeCountPtr
	return r
}

func (o *VolumeCreateRequest) SetStripeConstituentVolumeCount(newValue int) *VolumeCreateRequest {
	o.StripeConstituentVolumeCountPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) StripeOptimize() string {
	r := *o.StripeOptimizePtr
	return r
}

func (o *VolumeCreateRequest) SetStripeOptimize(newValue string) *VolumeCreateRequest {
	o.StripeOptimizePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) StripeWidth() int {
	r := *o.StripeWidthPtr
	return r
}

func (o *VolumeCreateRequest) SetStripeWidth(newValue int) *VolumeCreateRequest {
	o.StripeWidthPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) UnixPermissions() string {
	r := *o.UnixPermissionsPtr
	return r
}

func (o *VolumeCreateRequest) SetUnixPermissions(newValue string) *VolumeCreateRequest {
	o.UnixPermissionsPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) UserId() int {
	r := *o.UserIdPtr
	return r
}

func (o *VolumeCreateRequest) SetUserId(newValue int) *VolumeCreateRequest {
	o.UserIdPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) VmAlignSector() int {
	r := *o.VmAlignSectorPtr
	return r
}

func (o *VolumeCreateRequest) SetVmAlignSector(newValue int) *VolumeCreateRequest {
	o.VmAlignSectorPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) VmAlignSuffix() string {
	r := *o.VmAlignSuffixPtr
	return r
}

func (o *VolumeCreateRequest) SetVmAlignSuffix(newValue string) *VolumeCreateRequest {
	o.VmAlignSuffixPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) Volume() string {
	r := *o.VolumePtr
	return r
}

func (o *VolumeCreateRequest) SetVolume(newValue string) *VolumeCreateRequest {
	o.VolumePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) VolumeComment() string {
	r := *o.VolumeCommentPtr
	return r
}

func (o *VolumeCreateRequest) SetVolumeComment(newValue string) *VolumeCreateRequest {
	o.VolumeCommentPtr = &newValue
	return o
}

func (o *VolumeCreateRequest) VolumeSecurityStyle() string {
	r := *o.VolumeSecurityStylePtr
	return r
}

func (o *VolumeCreateRequest) SetVolumeSecurityStyle(newValue string) *VolumeCreateRequest {
	o.VolumeSecurityStylePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) VolumeState() string {
	r := *o.VolumeStatePtr
	return r
}

func (o *VolumeCreateRequest) SetVolumeState(newValue string) *VolumeCreateRequest {
	o.VolumeStatePtr = &newValue
	return o
}

func (o *VolumeCreateRequest) VolumeType() string {
	r := *o.VolumeTypePtr
	return r
}

func (o *VolumeCreateRequest) SetVolumeType(newValue string) *VolumeCreateRequest {
	o.VolumeTypePtr = &newValue
	return o
}

type VolumeCreateResponse struct {
	XMLName xml.Name `xml:"netapp"`

	ResponseVersion string `xml:"version,attr"`
	ResponseXmlns   string `xml:"xmlns,attr"`

	Result VolumeCreateResponseResult `xml:"results"`
}

func (o VolumeCreateResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "version", o.ResponseVersion))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "xmlns", o.ResponseXmlns))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "results", o.Result))
	return buffer.String()
}

type VolumeCreateResponseResult struct {
	XMLName xml.Name `xml:"results"`

	ResultStatusAttr string `xml:"status,attr"`
	ResultReasonAttr string `xml:"reason,attr"`
	ResultErrnoAttr  string `xml:"errno,attr"`
}

func (o *VolumeCreateResponse) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Debugf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeCreateResponse() *VolumeCreateResponse { return &VolumeCreateResponse{} }

func (o VolumeCreateResponseResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultStatusAttr", o.ResultStatusAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultReasonAttr", o.ResultReasonAttr))
	buffer.WriteString(fmt.Sprintf("%s: %s\n", "resultErrnoAttr", o.ResultErrnoAttr))
	return buffer.String()
}
