// Copyright 2016 NetApp, Inc. All Rights Reserved.

package azgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

type VolumeCloneInfoType struct {
	XMLName xml.Name `xml:"volume-clone-info"`

	AggregatePtr               *string           `xml:"aggregate"`
	BlockPercentageCompletePtr *int              `xml:"block-percentage-complete"`
	BlocksScannedPtr           *int              `xml:"blocks-scanned"`
	BlocksUpdatedPtr           *int              `xml:"blocks-updated"`
	DsidPtr                    *int              `xml:"dsid"`
	FlexCloneUsedPercentPtr    *int              `xml:"flexclone-used-percent"`
	InodePercentageCompletePtr *int              `xml:"inode-percentage-complete"`
	InodesProcessedPtr         *int              `xml:"inodes-processed"`
	InodesTotalPtr             *int              `xml:"inodes-total"`
	JunctionActivePtr	   *bool             `xml:"junction-active"`
	JunctionPathPtr            *JunctionPathType `xml:"junction-path"`
	MsidPtr                    *int              `xml:"msid"`
	ParentSnapshotPtr          *string           `xml:"parent-snapshot"`
	ParentVolumePtr            *string           `xml:"parent-volume"`
        ParentVserverPtr           *string           `xml:"parent-vserver"`
	QosPolicyGroupNamePtr      *string           `xml:"qos-policy-group-name"`
	SizePtr                    *int              `xml:"size"`
	SpaceGuaranteeEnabledPtr   *bool             `xml:"space-guarantee-enabled"`
	SpaceReservePtr            *string           `xml:"space-reserve"`
	SplitEstimatePtr           *int              `xml:"split-estimate"`
	StatePtr                   *string           `xml:"state"`
	UsedPtr                    *int              `xml:"used"`
	VolumePtr                  *string           `xml:"volume"`	
	VolumeTypePtr              *string           `xml:"volume-type"`
	VserverPtr                 *string           `xml:"vserver"`	
	VserverDrProtectionPtr     *string           `xml:"vserver-dr-protection"`
}

func (o *VolumeCloneInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} // TODO: handle better
	//fmt.Println(string(output))
	return string(output), err
}

func NewVolumeCloneInfoType() *VolumeCloneInfoType { return &VolumeCloneInfoType{} }

func (o VolumeCloneInfoType) String() string {
	var buffer bytes.Buffer
	if o.AggregatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "aggregate", *o.AggregatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("aggregate: nil\n"))
	}
	if o.BlockPercentageCompletePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "block-percentage-complete", *o.BlockPercentageCompletePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("block-percentage-complete: nil\n"))
	}
	if o.BlocksScannedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "blocks-scanned", *o.BlocksScannedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("blocks-scanned: nil\n"))
	}
	if o.BlocksUpdatedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "blocks-updated", *o.BlocksUpdatedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("blocks-updated: nil\n"))
	}
	if o.DsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "dsid", *o.DsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("dsid: nil\n"))
	}
	if o.FlexCloneUsedPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "flexclone-used-percent", *o.FlexCloneUsedPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("flexclone-used-percent: nil\n"))
	}
	if o.InodePercentageCompletePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "inode-percentage-complete", *o.InodePercentageCompletePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("inode-percentage-complete: nil\n"))
	}
	if o.InodesProcessedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "inodes-processed", *o.InodesProcessedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("inodes-processed: nil\n"))
	}
	if o.InodesTotalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "inodes-total", *o.InodesTotalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("inodes-total: nil\n"))
	}
	if o.JunctionActivePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "junction-active", *o.JunctionActivePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("junction-active: nil\n"))
	}
	if o.JunctionPathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "junction-path", *o.JunctionPathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("junction-path: nil\n"))
	}
	if o.MsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "msid", *o.MsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("msid: nil\n"))
	}
	if o.ParentSnapshotPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "parent-snapshot", *o.ParentSnapshotPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("parent-snapshot: nil\n"))
	}
	if o.ParentVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "parent-volume", *o.ParentVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("parent-volume: nil\n"))
	}
	if o.QosPolicyGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "qos-policy-group-name", o.QosPolicyGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("qos-policy-group-name: nil\n"))
	}
	if o.SizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size", *o.SizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size: nil\n"))
	}
	if o.SpaceGuaranteeEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-guarantee-enabled", *o.SpaceGuaranteeEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-guarantee-enabled: nil\n"))
	}
	if o.SpaceReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-reserve", *o.SpaceReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-reserve: nil\n"))
	}
	if o.SplitEstimatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "split-estimate", *o.SplitEstimatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("split-estimate: nil\n"))
	}
	if o.StatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "state", *o.StatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("state: nil\n"))
	}
	if o.UsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "used", *o.UsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("used: nil\n"))
	}
	if o.VolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume", *o.VolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume: nil\n"))
	}
	if o.VolumeTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-type", *o.VolumeTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-type: nil\n"))
	}
	if o.VserverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver", *o.VserverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver: nil\n"))
	}
	if o.VserverDrProtectionPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver-dr-protection", *o.VserverDrProtectionPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver-dr-protection: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeCloneInfoType) Aggregate() string {
	r := *o.AggregatePtr
	return r
}

func (o *VolumeCloneInfoType) SetAggregate(newValue string) *VolumeCloneInfoType {
	o.AggregatePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) BlockPercentageComplete() int {
	r := *o.BlockPercentageCompletePtr
	return r
}

func (o *VolumeCloneInfoType) SetBlockPercentageComplete(newValue int) *VolumeCloneInfoType {
	o.BlockPercentageCompletePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) BlocksScanned() int {
	r := *o.BlocksScannedPtr
	return r
}

func (o *VolumeCloneInfoType) SetBlocksScanned(newValue int) *VolumeCloneInfoType {
	o.BlocksScannedPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) BlocksUpdated() int {
	r := *o.BlocksUpdatedPtr
	return r
}

func (o *VolumeCloneInfoType) SetBlocksUpdated(newValue int) *VolumeCloneInfoType {
	o.BlocksUpdatedPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) Dsid() int {
	r := *o.DsidPtr
	return r
}

func (o *VolumeCloneInfoType) SetDsid(newValue int) *VolumeCloneInfoType {
	o.DsidPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) FlexCloneUsedPercent() int {
	r := *o.FlexCloneUsedPercentPtr
	return r
}

func (o *VolumeCloneInfoType) SetFlexCloneUsedPercent(newValue int) *VolumeCloneInfoType {
	o.FlexCloneUsedPercentPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) InodePercentageComplete() int {
	r := *o.InodePercentageCompletePtr
	return r
}

func (o *VolumeCloneInfoType) SetInodePercentageComplete(newValue int) *VolumeCloneInfoType {
	o.InodePercentageCompletePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) InodesProcessed() int {
	r := *o.InodesProcessedPtr
	return r
}

func (o *VolumeCloneInfoType) SetInodesProcessed(newValue int) *VolumeCloneInfoType {
	o.InodesProcessedPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) InodesTotal() int {
	r := *o.InodesTotalPtr
	return r
}

func (o *VolumeCloneInfoType) SetInodesTotal(newValue int) *VolumeCloneInfoType {
	o.InodesTotalPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) JunctionActive() bool {
	r := *o.JunctionActivePtr
	return r
}

func (o *VolumeCloneInfoType) SetJunctionActive(newValue bool) *VolumeCloneInfoType {
	o.JunctionActivePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) JunctionPath() JunctionPathType {
	r := *o.JunctionPathPtr
	return r
}

func (o *VolumeCloneInfoType) SetJunctionPath(newValue JunctionPathType) *VolumeCloneInfoType {
	o.JunctionPathPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) Msid() int {
	r := *o.MsidPtr
	return r
}

func (o *VolumeCloneInfoType) SetMsid(newValue int) *VolumeCloneInfoType {
	o.MsidPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) ParentSnapshot() string {
	r := *o.ParentSnapshotPtr
	return r
}

func (o *VolumeCloneInfoType) SetParentSnapshot(newValue string) *VolumeCloneInfoType {
	o.ParentSnapshotPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) ParentVolume() string {
	r := *o.ParentVolumePtr
	return r
}

func (o *VolumeCloneInfoType) SetParentVolume(newValue string) *VolumeCloneInfoType {
	o.ParentVolumePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) ParentVserver() string {
	r := *o.ParentVserverPtr
	return r
}

func (o *VolumeCloneInfoType) SetParentVserver(newValue string) *VolumeCloneInfoType {
	o.ParentVserverPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) QosPolicyGroupName() string {
	r := *o.QosPolicyGroupNamePtr
	return r
}

func (o *VolumeCloneInfoType) SetQosPolicyGroupName(newValue string) *VolumeCloneInfoType {
	o.QosPolicyGroupNamePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) Size() int {
	r := *o.SizePtr
	return r
}

func (o *VolumeCloneInfoType) SetSize(newValue int) *VolumeCloneInfoType {
	o.SizePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) SpaceGuaranteeEnabled() bool {
	r := *o.SpaceGuaranteeEnabledPtr
	return r
}

func (o *VolumeCloneInfoType) SetSpaceGuaranteeEnabled(newValue bool) *VolumeCloneInfoType {
	o.SpaceGuaranteeEnabledPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) SpaceReserve() string {
	r := *o.SpaceReservePtr
	return r
}

func (o *VolumeCloneInfoType) SetSpaceReserve(newValue string) *VolumeCloneInfoType {
	o.SpaceReservePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) SplitEstimate() int {
	r := *o.SplitEstimatePtr
	return r
}

func (o *VolumeCloneInfoType) SetSplitEstimate(newValue int) *VolumeCloneInfoType {
	o.SplitEstimatePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) State() string {
	r := *o.StatePtr
	return r
}

func (o *VolumeCloneInfoType) SetState(newValue string) *VolumeCloneInfoType {
	o.StatePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) Used() int {
	r := *o.UsedPtr
	return r
}

func (o *VolumeCloneInfoType) SetUsed(newValue int) *VolumeCloneInfoType {
	o.UsedPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) Volume() string {
	r := *o.VolumePtr
	return r
}

func (o *VolumeCloneInfoType) SetVolume(newValue string) *VolumeCloneInfoType {
	o.VolumePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) VolumeType() string {
	r := *o.VolumeTypePtr
	return r
}

func (o *VolumeCloneInfoType) SetVolumeType(newValue string) *VolumeCloneInfoType {
	o.VolumeTypePtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) Vserver() string {
	r := *o.VserverPtr
	return r
}

func (o *VolumeCloneInfoType) SetVserver(newValue string) *VolumeCloneInfoType {
	o.VserverPtr = &newValue
	return o
}

func (o *VolumeCloneInfoType) VserverDrProtection() string {
	r := *o.VserverDrProtectionPtr
	return r
}

func (o *VolumeCloneInfoType) SetVserverDrProtection(newValue string) *VolumeCloneInfoType {
	o.VserverDrProtectionPtr = &newValue
	return o
}


type NullableSizeType string

type VserverAggrInfoType struct {
	XMLName xml.Name `xml:"vserver-aggr-info"`

	// We need to use this for compatibility with ONTAP 9.
	//AggrAvailsizePtr *SizeType     `xml:"aggr-availsize"`
	AggrAvailsizePtr *NullableSizeType `xml:"aggr-availsize"`
	AggrNamePtr      *AggrNameType     `xml:"aggr-name"`
}

func (o *VserverAggrInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	} // TODO: handle better
	return string(output), err
}

func NewVserverAggrInfoType() *VserverAggrInfoType { return &VserverAggrInfoType{} }

func (o VserverAggrInfoType) String() string {
	var buffer bytes.Buffer
	if o.AggrAvailsizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "aggr-availsize", *o.AggrAvailsizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("aggr-availsize: nil\n"))
	}
	if o.AggrNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "aggr-name", *o.AggrNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("aggr-name: nil\n"))
	}
	return buffer.String()
}

func (o *VserverAggrInfoType) AggrAvailsize() (SizeType, error) {
	r := *o.AggrAvailsizePtr
	if r == "" {
		return 0, nil
	}
	ret, err := strconv.Atoi(string(r))
	return SizeType(ret), err
}

func (o *VserverAggrInfoType) SetAggrAvailsize(newValue SizeType) *VserverAggrInfoType {
	n := NullableSizeType(strconv.Itoa(int(newValue)))
	o.AggrAvailsizePtr = &n
	return o
}

func (o *VserverAggrInfoType) AggrName() AggrNameType {
	r := *o.AggrNamePtr
	return r
}

func (o *VserverAggrInfoType) SetAggrName(newValue AggrNameType) *VserverAggrInfoType {
	o.AggrNamePtr = &newValue
	return o
}

type ProtocolType string

type AntivirusPolicyType string

type NmswitchType string

type NsswitchType string

type NisDomainType string

type VsoperstateType string

type VsopstopreasonType string

type SecurityStyleEnumType string

type SnapshotPolicyType string

type VsadminstateType string

type VserverInfoType struct {
	XMLName xml.Name `xml:"vserver-info"`

	AggrListPtr                      []AggrNameType         `xml:"aggr-list>aggr-name"`
	AllowedProtocolsPtr              []ProtocolType         `xml:"allowed-protocols>protocol"`
	AntivirusOnAccessPolicyPtr       *AntivirusPolicyType   `xml:"antivirus-on-access-policy"`
	CommentPtr                       *string                `xml:"comment"`
	DisallowedProtocolsPtr           []ProtocolType         `xml:"disallowed-protocols>protocol"`
	IpspacePtr                       *string                `xml:"ipspace"`
	IsConfigLockedForChangesPtr      *bool                  `xml:"is-config-locked-for-changes"`
	IsRepositoryVserverPtr           *bool                  `xml:"is-repository-vserver"`
	LanguagePtr                      *LanguageCodeType      `xml:"language"`
	LdapDomainPtr                    *string                `xml:"ldap-domain"`
	MaxVolumesPtr                    *string                `xml:"max-volumes"`
	NameMappingSwitchPtr             []NmswitchType         `xml:"name-mapping-switch>nmswitch"`
	NameServerSwitchPtr              []NsswitchType         `xml:"name-server-switch>nsswitch"`
	NisDomainPtr                     *NisDomainType         `xml:"nis-domain"`
	OperationalStatePtr              *VsoperstateType       `xml:"operational-state"`
	OperationalStateStoppedReasonPtr *VsopstopreasonType    `xml:"operational-state-stopped-reason"`
	QosPolicyGroupPtr                *string                `xml:"qos-policy-group"`
	QuotaPolicyPtr                   *string                `xml:"quota-policy"`
	RootVolumePtr                    *VolumeNameType        `xml:"root-volume"`
	RootVolumeAggregatePtr           *AggrNameType          `xml:"root-volume-aggregate"`
	RootVolumeSecurityStylePtr       *SecurityStyleEnumType `xml:"root-volume-security-style"`
	SnapshotPolicyPtr                *SnapshotPolicyType    `xml:"snapshot-policy"`
	StatePtr                         *VsadminstateType      `xml:"state"`
	UuidPtr                          *UuidType              `xml:"uuid"`
	VolumeDeleteRetentionHoursPtr    *int                   `xml:"volume-delete-retention-hours"`
	VserverAggrInfoListPtr           []VserverAggrInfoType  `xml:"vserver-aggr-info-list>vserver-aggr-info"`
	VserverNamePtr                   *string                `xml:"vserver-name"`
	VserverSubtypePtr                *string                `xml:"vserver-subtype"`
	VserverTypePtr                   *string                `xml:"vserver-type"`
}

func (o *VserverInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	} // TODO: handle better
	return string(output), err
}

func NewVserverInfoType() *VserverInfoType { return &VserverInfoType{} }

func (o VserverInfoType) String() string {
	var buffer bytes.Buffer
	if o.AggrListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "aggr-list", o.AggrListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("aggr-list: nil\n"))
	}
	if o.AllowedProtocolsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "allowed-protocols", o.AllowedProtocolsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("allowed-protocols: nil\n"))
	}
	if o.AntivirusOnAccessPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "antivirus-on-access-policy", *o.AntivirusOnAccessPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("antivirus-on-access-policy: nil\n"))
	}
	if o.CommentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "comment", *o.CommentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("comment: nil\n"))
	}
	if o.DisallowedProtocolsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "disallowed-protocols", o.DisallowedProtocolsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("disallowed-protocols: nil\n"))
	}
	if o.IpspacePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "ipspace", *o.IpspacePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("ipspace: nil\n"))
	}
	if o.IsConfigLockedForChangesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-config-locked-for-changes", *o.IsConfigLockedForChangesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-config-locked-for-changes: nil\n"))
	}
	if o.IsRepositoryVserverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-repository-vserver", *o.IsRepositoryVserverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-repository-vserver: nil\n"))
	}
	if o.LanguagePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "language", *o.LanguagePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("language: nil\n"))
	}
	if o.LdapDomainPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "ldap-domain", *o.LdapDomainPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("ldap-domain: nil\n"))
	}
	if o.MaxVolumesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-volumes", *o.MaxVolumesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-volumes: nil\n"))
	}
	if o.NameMappingSwitchPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name-mapping-switch", o.NameMappingSwitchPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name-mapping-switch: nil\n"))
	}
	if o.NameServerSwitchPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name-server-switch", o.NameServerSwitchPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name-server-switch: nil\n"))
	}
	if o.NisDomainPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "nis-domain", *o.NisDomainPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("nis-domain: nil\n"))
	}
	if o.OperationalStatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "operational-state", *o.OperationalStatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("operational-state: nil\n"))
	}
	if o.OperationalStateStoppedReasonPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "operational-state-stopped-reason", *o.OperationalStateStoppedReasonPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("operational-state-stopped-reason: nil\n"))
	}
	if o.QosPolicyGroupPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "qos-policy-group", *o.QosPolicyGroupPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("qos-policy-group: nil\n"))
	}
	if o.QuotaPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "quota-policy", *o.QuotaPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("quota-policy: nil\n"))
	}
	if o.RootVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "root-volume", *o.RootVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("root-volume: nil\n"))
	}
	if o.RootVolumeAggregatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "root-volume-aggregate", *o.RootVolumeAggregatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("root-volume-aggregate: nil\n"))
	}
	if o.RootVolumeSecurityStylePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "root-volume-security-style", *o.RootVolumeSecurityStylePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("root-volume-security-style: nil\n"))
	}
	if o.SnapshotPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-policy", *o.SnapshotPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-policy: nil\n"))
	}
	if o.StatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "state", *o.StatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("state: nil\n"))
	}
	if o.UuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "uuid", *o.UuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("uuid: nil\n"))
	}
	if o.VolumeDeleteRetentionHoursPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-delete-retention-hours", *o.VolumeDeleteRetentionHoursPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-delete-retention-hours: nil\n"))
	}
	if o.VserverAggrInfoListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver-aggr-info-list", o.VserverAggrInfoListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver-aggr-info-list: nil\n"))
	}
	if o.VserverNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver-name", *o.VserverNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver-name: nil\n"))
	}
	if o.VserverSubtypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver-subtype", *o.VserverSubtypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver-subtype: nil\n"))
	}
	if o.VserverTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver-type", *o.VserverTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver-type: nil\n"))
	}
	return buffer.String()
}

func (o *VserverInfoType) AggrList() []AggrNameType {
	r := o.AggrListPtr
	return r
}

func (o *VserverInfoType) SetAggrList(newValue []AggrNameType) *VserverInfoType {
	newSlice := make([]AggrNameType, len(newValue))
	copy(newSlice, newValue)
	o.AggrListPtr = newSlice
	return o
}

func (o *VserverInfoType) AllowedProtocols() []ProtocolType {
	r := o.AllowedProtocolsPtr
	return r
}

func (o *VserverInfoType) SetAllowedProtocols(newValue []ProtocolType) *VserverInfoType {
	newSlice := make([]ProtocolType, len(newValue))
	copy(newSlice, newValue)
	o.AllowedProtocolsPtr = newSlice
	return o
}

func (o *VserverInfoType) AntivirusOnAccessPolicy() AntivirusPolicyType {
	r := *o.AntivirusOnAccessPolicyPtr
	return r
}

func (o *VserverInfoType) SetAntivirusOnAccessPolicy(newValue AntivirusPolicyType) *VserverInfoType {
	o.AntivirusOnAccessPolicyPtr = &newValue
	return o
}

func (o *VserverInfoType) Comment() string {
	r := *o.CommentPtr
	return r
}

func (o *VserverInfoType) SetComment(newValue string) *VserverInfoType {
	o.CommentPtr = &newValue
	return o
}

func (o *VserverInfoType) DisallowedProtocols() []ProtocolType {
	r := o.DisallowedProtocolsPtr
	return r
}

func (o *VserverInfoType) SetDisallowedProtocols(newValue []ProtocolType) *VserverInfoType {
	newSlice := make([]ProtocolType, len(newValue))
	copy(newSlice, newValue)
	o.DisallowedProtocolsPtr = newSlice
	return o
}

func (o *VserverInfoType) Ipspace() string {
	r := *o.IpspacePtr
	return r
}

func (o *VserverInfoType) SetIpspace(newValue string) *VserverInfoType {
	o.IpspacePtr = &newValue
	return o
}

func (o *VserverInfoType) IsConfigLockedForChanges() bool {
	r := *o.IsConfigLockedForChangesPtr
	return r
}

func (o *VserverInfoType) SetIsConfigLockedForChanges(newValue bool) *VserverInfoType {
	o.IsConfigLockedForChangesPtr = &newValue
	return o
}

func (o *VserverInfoType) IsRepositoryVserver() bool {
	r := *o.IsRepositoryVserverPtr
	return r
}

func (o *VserverInfoType) SetIsRepositoryVserver(newValue bool) *VserverInfoType {
	o.IsRepositoryVserverPtr = &newValue
	return o
}

func (o *VserverInfoType) Language() LanguageCodeType {
	r := *o.LanguagePtr
	return r
}

func (o *VserverInfoType) SetLanguage(newValue LanguageCodeType) *VserverInfoType {
	o.LanguagePtr = &newValue
	return o
}

func (o *VserverInfoType) LdapDomain() string {
	r := *o.LdapDomainPtr
	return r
}

func (o *VserverInfoType) SetLdapDomain(newValue string) *VserverInfoType {
	o.LdapDomainPtr = &newValue
	return o
}

func (o *VserverInfoType) MaxVolumes() string {
	r := *o.MaxVolumesPtr
	return r
}

func (o *VserverInfoType) SetMaxVolumes(newValue string) *VserverInfoType {
	o.MaxVolumesPtr = &newValue
	return o
}

func (o *VserverInfoType) NameMappingSwitch() []NmswitchType {
	r := o.NameMappingSwitchPtr
	return r
}

func (o *VserverInfoType) SetNameMappingSwitch(newValue []NmswitchType) *VserverInfoType {
	newSlice := make([]NmswitchType, len(newValue))
	copy(newSlice, newValue)
	o.NameMappingSwitchPtr = newSlice
	return o
}

func (o *VserverInfoType) NameServerSwitch() []NsswitchType {
	r := o.NameServerSwitchPtr
	return r
}

func (o *VserverInfoType) SetNameServerSwitch(newValue []NsswitchType) *VserverInfoType {
	newSlice := make([]NsswitchType, len(newValue))
	copy(newSlice, newValue)
	o.NameServerSwitchPtr = newSlice
	return o
}

func (o *VserverInfoType) NisDomain() NisDomainType {
	r := *o.NisDomainPtr
	return r
}

func (o *VserverInfoType) SetNisDomain(newValue NisDomainType) *VserverInfoType {
	o.NisDomainPtr = &newValue
	return o
}

func (o *VserverInfoType) OperationalState() VsoperstateType {
	r := *o.OperationalStatePtr
	return r
}

func (o *VserverInfoType) SetOperationalState(newValue VsoperstateType) *VserverInfoType {
	o.OperationalStatePtr = &newValue
	return o
}

func (o *VserverInfoType) OperationalStateStoppedReason() VsopstopreasonType {
	r := *o.OperationalStateStoppedReasonPtr
	return r
}

func (o *VserverInfoType) SetOperationalStateStoppedReason(newValue VsopstopreasonType) *VserverInfoType {
	o.OperationalStateStoppedReasonPtr = &newValue
	return o
}

func (o *VserverInfoType) QosPolicyGroup() string {
	r := *o.QosPolicyGroupPtr
	return r
}

func (o *VserverInfoType) SetQosPolicyGroup(newValue string) *VserverInfoType {
	o.QosPolicyGroupPtr = &newValue
	return o
}

func (o *VserverInfoType) QuotaPolicy() string {
	r := *o.QuotaPolicyPtr
	return r
}

func (o *VserverInfoType) SetQuotaPolicy(newValue string) *VserverInfoType {
	o.QuotaPolicyPtr = &newValue
	return o
}

func (o *VserverInfoType) RootVolume() VolumeNameType {
	r := *o.RootVolumePtr
	return r
}

func (o *VserverInfoType) SetRootVolume(newValue VolumeNameType) *VserverInfoType {
	o.RootVolumePtr = &newValue
	return o
}

func (o *VserverInfoType) RootVolumeAggregate() AggrNameType {
	r := *o.RootVolumeAggregatePtr
	return r
}

func (o *VserverInfoType) SetRootVolumeAggregate(newValue AggrNameType) *VserverInfoType {
	o.RootVolumeAggregatePtr = &newValue
	return o
}

func (o *VserverInfoType) RootVolumeSecurityStyle() SecurityStyleEnumType {
	r := *o.RootVolumeSecurityStylePtr
	return r
}

func (o *VserverInfoType) SetRootVolumeSecurityStyle(newValue SecurityStyleEnumType) *VserverInfoType {
	o.RootVolumeSecurityStylePtr = &newValue
	return o
}

func (o *VserverInfoType) SnapshotPolicy() SnapshotPolicyType {
	r := *o.SnapshotPolicyPtr
	return r
}

func (o *VserverInfoType) SetSnapshotPolicy(newValue SnapshotPolicyType) *VserverInfoType {
	o.SnapshotPolicyPtr = &newValue
	return o
}

func (o *VserverInfoType) State() VsadminstateType {
	r := *o.StatePtr
	return r
}

func (o *VserverInfoType) SetState(newValue VsadminstateType) *VserverInfoType {
	o.StatePtr = &newValue
	return o
}

func (o *VserverInfoType) Uuid() UuidType {
	r := *o.UuidPtr
	return r
}

func (o *VserverInfoType) SetUuid(newValue UuidType) *VserverInfoType {
	o.UuidPtr = &newValue
	return o
}

func (o *VserverInfoType) VolumeDeleteRetentionHours() int {
	r := *o.VolumeDeleteRetentionHoursPtr
	return r
}

func (o *VserverInfoType) SetVolumeDeleteRetentionHours(newValue int) *VserverInfoType {
	o.VolumeDeleteRetentionHoursPtr = &newValue
	return o
}

func (o *VserverInfoType) VserverAggrInfoList() []VserverAggrInfoType {
	r := o.VserverAggrInfoListPtr
	return r
}

func (o *VserverInfoType) SetVserverAggrInfoList(newValue []VserverAggrInfoType) *VserverInfoType {
	newSlice := make([]VserverAggrInfoType, len(newValue))
	copy(newSlice, newValue)
	o.VserverAggrInfoListPtr = newSlice
	return o
}

func (o *VserverInfoType) VserverName() string {
	r := *o.VserverNamePtr
	return r
}

func (o *VserverInfoType) SetVserverName(newValue string) *VserverInfoType {
	o.VserverNamePtr = &newValue
	return o
}

func (o *VserverInfoType) VserverSubtype() string {
	r := *o.VserverSubtypePtr
	return r
}

func (o *VserverInfoType) SetVserverSubtype(newValue string) *VserverInfoType {
	o.VserverSubtypePtr = &newValue
	return o
}

func (o *VserverInfoType) VserverType() string {
	r := *o.VserverTypePtr
	return r
}

func (o *VserverInfoType) SetVserverType(newValue string) *VserverInfoType {
	o.VserverTypePtr = &newValue
	return o
}

// SnaplocktypeType is a structure to represent a snaplocktype ZAPI object
type SnaplocktypeType string

// AggregatetypeType is a structure to represent a aggregatetype ZAPI object
type AggregatetypeType string

// ShowAggregatesType is a structure to represent a show-aggregates ZAPI object
type ShowAggregatesType struct {
	XMLName xml.Name `xml:"show-aggregates"`

	AggregateNamePtr *AggrNameType      `xml:"aggregate-name"`
	AggregateTypePtr *AggregatetypeType `xml:"aggregate-type"`
	AvailableSizePtr *SizeType          `xml:"available-size"`
	SnaplockTypePtr  *SnaplocktypeType  `xml:"snaplock-type"`
	VserverNamePtr   *string            `xml:"vserver-name"`
}

// ToXML converts this object into an xml string representation
func (o *ShowAggregatesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

// NewShowAggregatesType is a factory method for creating new instances of ShowAggregatesType objects
func NewShowAggregatesType() *ShowAggregatesType { return &ShowAggregatesType{} }

// String returns a string representation of this object's fields and implements the Stringer interface
func (o ShowAggregatesType) String() string {
	var buffer bytes.Buffer
	if o.AggregateNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "aggregate-name", *o.AggregateNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("aggregate-name: nil\n"))
	}
	if o.AggregateTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "aggregate-type", *o.AggregateTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("aggregate-type: nil\n"))
	}
	if o.AvailableSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "available-size", *o.AvailableSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("available-size: nil\n"))
	}
	if o.SnaplockTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snaplock-type", *o.SnaplockTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snaplock-type: nil\n"))
	}
	if o.VserverNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver-name", *o.VserverNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver-name: nil\n"))
	}
	return buffer.String()
}

// AggregateName is a fluent style 'getter' method that can be chained
func (o *ShowAggregatesType) AggregateName() AggrNameType {
	r := *o.AggregateNamePtr
	return r
}

// SetAggregateName is a fluent style 'setter' method that can be chained
func (o *ShowAggregatesType) SetAggregateName(newValue AggrNameType) *ShowAggregatesType {
	o.AggregateNamePtr = &newValue
	return o
}

// AggregateType is a fluent style 'getter' method that can be chained
func (o *ShowAggregatesType) AggregateType() AggregatetypeType {
	r := *o.AggregateTypePtr
	return r
}

// SetAggregateType is a fluent style 'setter' method that can be chained
func (o *ShowAggregatesType) SetAggregateType(newValue AggregatetypeType) *ShowAggregatesType {
	o.AggregateTypePtr = &newValue
	return o
}

// AvailableSize is a fluent style 'getter' method that can be chained
func (o *ShowAggregatesType) AvailableSize() SizeType {
	r := *o.AvailableSizePtr
	return r
}

// SetAvailableSize is a fluent style 'setter' method that can be chained
func (o *ShowAggregatesType) SetAvailableSize(newValue SizeType) *ShowAggregatesType {
	o.AvailableSizePtr = &newValue
	return o
}

// SnaplockType is a fluent style 'getter' method that can be chained
func (o *ShowAggregatesType) SnaplockType() SnaplocktypeType {
	r := *o.SnaplockTypePtr
	return r
}

// SetSnaplockType is a fluent style 'setter' method that can be chained
func (o *ShowAggregatesType) SetSnaplockType(newValue SnaplocktypeType) *ShowAggregatesType {
	o.SnaplockTypePtr = &newValue
	return o
}

// VserverName is a fluent style 'getter' method that can be chained
func (o *ShowAggregatesType) VserverName() string {
	r := *o.VserverNamePtr
	return r
}

// SetVserverName is a fluent style 'setter' method that can be chained
func (o *ShowAggregatesType) SetVserverName(newValue string) *ShowAggregatesType {
	o.VserverNamePtr = &newValue
	return o
}

type AggrAttributesType struct {
	XMLName               xml.Name                `xml:"aggr-attributes"`
	AggrRaidAttributesPtr *AggrRaidAttributesType `xml:"aggr-raid-attributes"`
	AggregateNamePtr      *string                 `xml:"aggregate-name"`
}

func (o *AggrAttributesType) AggrRaidAttributes() AggrRaidAttributesType {
	r := *o.AggrRaidAttributesPtr
	return r
}

func (o *AggrAttributesType) AggregateName() string {
	r := *o.AggregateNamePtr
	return r
}

type AggrRaidAttributesType struct {
	AggregateTypePtr *string `xml:"aggregate-type"`
	RaidTypePtr      *string `xml:"raid-type"`
}

func (o *AggrRaidAttributesType) AggregateType() string {
	r := *o.AggregateTypePtr
	return r
}

func (o *AggrRaidAttributesType) RaidType() string {
	r := *o.RaidTypePtr
	return r
}

type VolumeModifyIterInfoType struct {
	XMLName xml.Name `xml:"volume-modify-iter-info"`

	ErrorCodePtr    *int                  `xml:"error-code"`
	ErrorMessagePtr *string               `xml:"error-message"`
	VolumeKeyPtr    *VolumeAttributesType `xml:"volume-key"`
}

func (o *VolumeModifyIterInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeModifyIterInfoType() *VolumeModifyIterInfoType { return &VolumeModifyIterInfoType{} }

func (o VolumeModifyIterInfoType) String() string {
	var buffer bytes.Buffer
	if o.ErrorCodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "error-code", *o.ErrorCodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("error-code: nil\n"))
	}
	if o.ErrorMessagePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "error-message", *o.ErrorMessagePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("error-message: nil\n"))
	}
	if o.VolumeKeyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-key", *o.VolumeKeyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-key: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeModifyIterInfoType) ErrorCode() int {
	r := *o.ErrorCodePtr
	return r
}

func (o *VolumeModifyIterInfoType) SetErrorCode(newValue int) *VolumeModifyIterInfoType {
	o.ErrorCodePtr = &newValue
	return o
}

func (o *VolumeModifyIterInfoType) ErrorMessage() string {
	r := *o.ErrorMessagePtr
	return r
}

func (o *VolumeModifyIterInfoType) SetErrorMessage(newValue string) *VolumeModifyIterInfoType {
	o.ErrorMessagePtr = &newValue
	return o
}

func (o *VolumeModifyIterInfoType) VolumeKey() VolumeAttributesType {
	r := *o.VolumeKeyPtr
	return r
}

func (o *VolumeModifyIterInfoType) SetVolumeKey(newValue VolumeAttributesType) *VolumeModifyIterInfoType {
	o.VolumeKeyPtr = &newValue
	return o
}

type VolumeVmAlignAttributesType struct {
	XMLName xml.Name `xml:"volume-vm-align-attributes"`

	VmAlignSectorPtr *int    `xml:"vm-align-sector"`
	VmAlignSuffixPtr *string `xml:"vm-align-suffix"`
}

func (o *VolumeVmAlignAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeVmAlignAttributesType() *VolumeVmAlignAttributesType {
	return &VolumeVmAlignAttributesType{}
}

func (o VolumeVmAlignAttributesType) String() string {
	var buffer bytes.Buffer
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
	return buffer.String()
}

func (o *VolumeVmAlignAttributesType) VmAlignSector() int {
	r := *o.VmAlignSectorPtr
	return r
}

func (o *VolumeVmAlignAttributesType) SetVmAlignSector(newValue int) *VolumeVmAlignAttributesType {
	o.VmAlignSectorPtr = &newValue
	return o
}

func (o *VolumeVmAlignAttributesType) VmAlignSuffix() string {
	r := *o.VmAlignSuffixPtr
	return r
}

func (o *VolumeVmAlignAttributesType) SetVmAlignSuffix(newValue string) *VolumeVmAlignAttributesType {
	o.VmAlignSuffixPtr = &newValue
	return o
}

type VolumeTransitionAttributesType struct {
	XMLName xml.Name `xml:"volume-transition-attributes"`

	IsCopiedForTransitionPtr *bool   `xml:"is-copied-for-transition"`
	IsTransitionedPtr        *bool   `xml:"is-transitioned"`
	TransitionBehaviorPtr    *string `xml:"transition-behavior"`
}

func (o *VolumeTransitionAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeTransitionAttributesType() *VolumeTransitionAttributesType {
	return &VolumeTransitionAttributesType{}
}

func (o VolumeTransitionAttributesType) String() string {
	var buffer bytes.Buffer
	if o.IsCopiedForTransitionPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-copied-for-transition", *o.IsCopiedForTransitionPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-copied-for-transition: nil\n"))
	}
	if o.IsTransitionedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-transitioned", *o.IsTransitionedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-transitioned: nil\n"))
	}
	if o.TransitionBehaviorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "transition-behavior", *o.TransitionBehaviorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("transition-behavior: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeTransitionAttributesType) IsCopiedForTransition() bool {
	r := *o.IsCopiedForTransitionPtr
	return r
}

func (o *VolumeTransitionAttributesType) SetIsCopiedForTransition(newValue bool) *VolumeTransitionAttributesType {
	o.IsCopiedForTransitionPtr = &newValue
	return o
}

func (o *VolumeTransitionAttributesType) IsTransitioned() bool {
	r := *o.IsTransitionedPtr
	return r
}

func (o *VolumeTransitionAttributesType) SetIsTransitioned(newValue bool) *VolumeTransitionAttributesType {
	o.IsTransitionedPtr = &newValue
	return o
}

func (o *VolumeTransitionAttributesType) TransitionBehavior() string {
	r := *o.TransitionBehaviorPtr
	return r
}

func (o *VolumeTransitionAttributesType) SetTransitionBehavior(newValue string) *VolumeTransitionAttributesType {
	o.TransitionBehaviorPtr = &newValue
	return o
}

type VolumeStateAttributesType struct {
	XMLName xml.Name `xml:"volume-state-attributes"`

	BecomeNodeRootAfterRebootPtr *bool   `xml:"become-node-root-after-reboot"`
	ForceNvfailOnDrPtr           *bool   `xml:"force-nvfail-on-dr"`
	IgnoreInconsistentPtr        *bool   `xml:"ignore-inconsistent"`
	InNvfailedStatePtr           *bool   `xml:"in-nvfailed-state"`
	IsClusterVolumePtr           *bool   `xml:"is-cluster-volume"`
	IsConstituentPtr             *bool   `xml:"is-constituent"`
	IsInconsistentPtr            *bool   `xml:"is-inconsistent"`
	IsInvalidPtr                 *bool   `xml:"is-invalid"`
	IsJunctionActivePtr          *bool   `xml:"is-junction-active"`
	IsMovingPtr                  *bool   `xml:"is-moving"`
	IsNodeRootPtr                *bool   `xml:"is-node-root"`
	IsNvfailEnabledPtr           *bool   `xml:"is-nvfail-enabled"`
	IsQuiescedInMemoryPtr        *bool   `xml:"is-quiesced-in-memory"`
	IsQuiescedOnDiskPtr          *bool   `xml:"is-quiesced-on-disk"`
	IsUnrecoverablePtr           *bool   `xml:"is-unrecoverable"`
	IsVolumeInCutoverPtr         *bool   `xml:"is-volume-in-cutover"`
	IsVserverRootPtr             *bool   `xml:"is-vserver-root"`
	StatePtr                     *string `xml:"state"`
}

func (o *VolumeStateAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeStateAttributesType() *VolumeStateAttributesType { return &VolumeStateAttributesType{} }

func (o VolumeStateAttributesType) String() string {
	var buffer bytes.Buffer
	if o.BecomeNodeRootAfterRebootPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "become-node-root-after-reboot", *o.BecomeNodeRootAfterRebootPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("become-node-root-after-reboot: nil\n"))
	}
	if o.ForceNvfailOnDrPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force-nvfail-on-dr", *o.ForceNvfailOnDrPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force-nvfail-on-dr: nil\n"))
	}
	if o.IgnoreInconsistentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "ignore-inconsistent", *o.IgnoreInconsistentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("ignore-inconsistent: nil\n"))
	}
	if o.InNvfailedStatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "in-nvfailed-state", *o.InNvfailedStatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("in-nvfailed-state: nil\n"))
	}
	if o.IsClusterVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-cluster-volume", *o.IsClusterVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-cluster-volume: nil\n"))
	}
	if o.IsConstituentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-constituent", *o.IsConstituentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-constituent: nil\n"))
	}
	if o.IsInconsistentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-inconsistent", *o.IsInconsistentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-inconsistent: nil\n"))
	}
	if o.IsInvalidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-invalid", *o.IsInvalidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-invalid: nil\n"))
	}
	if o.IsJunctionActivePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-junction-active", *o.IsJunctionActivePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-junction-active: nil\n"))
	}
	if o.IsMovingPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-moving", *o.IsMovingPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-moving: nil\n"))
	}
	if o.IsNodeRootPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-node-root", *o.IsNodeRootPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-node-root: nil\n"))
	}
	if o.IsNvfailEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-nvfail-enabled", *o.IsNvfailEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-nvfail-enabled: nil\n"))
	}
	if o.IsQuiescedInMemoryPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-quiesced-in-memory", *o.IsQuiescedInMemoryPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-quiesced-in-memory: nil\n"))
	}
	if o.IsQuiescedOnDiskPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-quiesced-on-disk", *o.IsQuiescedOnDiskPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-quiesced-on-disk: nil\n"))
	}
	if o.IsUnrecoverablePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-unrecoverable", *o.IsUnrecoverablePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-unrecoverable: nil\n"))
	}
	if o.IsVolumeInCutoverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-volume-in-cutover", *o.IsVolumeInCutoverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-volume-in-cutover: nil\n"))
	}
	if o.IsVserverRootPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-vserver-root", *o.IsVserverRootPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-vserver-root: nil\n"))
	}
	if o.StatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "state", *o.StatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("state: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeStateAttributesType) BecomeNodeRootAfterReboot() bool {
	r := *o.BecomeNodeRootAfterRebootPtr
	return r
}

func (o *VolumeStateAttributesType) SetBecomeNodeRootAfterReboot(newValue bool) *VolumeStateAttributesType {
	o.BecomeNodeRootAfterRebootPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) ForceNvfailOnDr() bool {
	r := *o.ForceNvfailOnDrPtr
	return r
}

func (o *VolumeStateAttributesType) SetForceNvfailOnDr(newValue bool) *VolumeStateAttributesType {
	o.ForceNvfailOnDrPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IgnoreInconsistent() bool {
	r := *o.IgnoreInconsistentPtr
	return r
}

func (o *VolumeStateAttributesType) SetIgnoreInconsistent(newValue bool) *VolumeStateAttributesType {
	o.IgnoreInconsistentPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) InNvfailedState() bool {
	r := *o.InNvfailedStatePtr
	return r
}

func (o *VolumeStateAttributesType) SetInNvfailedState(newValue bool) *VolumeStateAttributesType {
	o.InNvfailedStatePtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsClusterVolume() bool {
	r := *o.IsClusterVolumePtr
	return r
}

func (o *VolumeStateAttributesType) SetIsClusterVolume(newValue bool) *VolumeStateAttributesType {
	o.IsClusterVolumePtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsConstituent() bool {
	r := *o.IsConstituentPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsConstituent(newValue bool) *VolumeStateAttributesType {
	o.IsConstituentPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsInconsistent() bool {
	r := *o.IsInconsistentPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsInconsistent(newValue bool) *VolumeStateAttributesType {
	o.IsInconsistentPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsInvalid() bool {
	r := *o.IsInvalidPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsInvalid(newValue bool) *VolumeStateAttributesType {
	o.IsInvalidPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsJunctionActive() bool {
	r := *o.IsJunctionActivePtr
	return r
}

func (o *VolumeStateAttributesType) SetIsJunctionActive(newValue bool) *VolumeStateAttributesType {
	o.IsJunctionActivePtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsMoving() bool {
	r := *o.IsMovingPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsMoving(newValue bool) *VolumeStateAttributesType {
	o.IsMovingPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsNodeRoot() bool {
	r := *o.IsNodeRootPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsNodeRoot(newValue bool) *VolumeStateAttributesType {
	o.IsNodeRootPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsNvfailEnabled() bool {
	r := *o.IsNvfailEnabledPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsNvfailEnabled(newValue bool) *VolumeStateAttributesType {
	o.IsNvfailEnabledPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsQuiescedInMemory() bool {
	r := *o.IsQuiescedInMemoryPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsQuiescedInMemory(newValue bool) *VolumeStateAttributesType {
	o.IsQuiescedInMemoryPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsQuiescedOnDisk() bool {
	r := *o.IsQuiescedOnDiskPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsQuiescedOnDisk(newValue bool) *VolumeStateAttributesType {
	o.IsQuiescedOnDiskPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsUnrecoverable() bool {
	r := *o.IsUnrecoverablePtr
	return r
}

func (o *VolumeStateAttributesType) SetIsUnrecoverable(newValue bool) *VolumeStateAttributesType {
	o.IsUnrecoverablePtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsVolumeInCutover() bool {
	r := *o.IsVolumeInCutoverPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsVolumeInCutover(newValue bool) *VolumeStateAttributesType {
	o.IsVolumeInCutoverPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) IsVserverRoot() bool {
	r := *o.IsVserverRootPtr
	return r
}

func (o *VolumeStateAttributesType) SetIsVserverRoot(newValue bool) *VolumeStateAttributesType {
	o.IsVserverRootPtr = &newValue
	return o
}

func (o *VolumeStateAttributesType) State() string {
	r := *o.StatePtr
	return r
}

func (o *VolumeStateAttributesType) SetState(newValue string) *VolumeStateAttributesType {
	o.StatePtr = &newValue
	return o
}

type VolumeSpaceAttributesType struct {
	XMLName xml.Name `xml:"volume-space-attributes"`

	FilesystemSizePtr                  *int    `xml:"filesystem-size"`
	IsFilesysSizeFixedPtr              *bool   `xml:"is-filesys-size-fixed"`
	IsSpaceGuaranteeEnabledPtr         *bool   `xml:"is-space-guarantee-enabled"`
	OverwriteReservePtr                *int    `xml:"overwrite-reserve"`
	OverwriteReserveRequiredPtr        *int    `xml:"overwrite-reserve-required"`
	OverwriteReserveUsedPtr            *int    `xml:"overwrite-reserve-used"`
	OverwriteReserveUsedActualPtr      *int    `xml:"overwrite-reserve-used-actual"`
	PercentageFractionalReservePtr     *int    `xml:"percentage-fractional-reserve"`
	PercentageSizeUsedPtr              *int    `xml:"percentage-size-used"`
	PercentageSnapshotReservePtr       *int    `xml:"percentage-snapshot-reserve"`
	PercentageSnapshotReserveUsedPtr   *int    `xml:"percentage-snapshot-reserve-used"`
	PhysicalUsedPtr                    *int    `xml:"physical-used"`
	PhysicalUsedPercentPtr             *int    `xml:"physical-used-percent"`
	SizePtr                            *int    `xml:"size"`
	SizeAvailablePtr                   *int    `xml:"size-available"`
	SizeAvailableForSnapshotsPtr       *int    `xml:"size-available-for-snapshots"`
	SizeTotalPtr                       *int    `xml:"size-total"`
	SizeUsedPtr                        *int    `xml:"size-used"`
	SizeUsedBySnapshotsPtr             *int    `xml:"size-used-by-snapshots"`
	SnapshotReserveSizePtr             *int    `xml:"snapshot-reserve-size"`
	SpaceFullThresholdPercentPtr       *int    `xml:"space-full-threshold-percent"`
	SpaceGuaranteePtr                  *string `xml:"space-guarantee"`
	SpaceMgmtOptionTryFirstPtr         *string `xml:"space-mgmt-option-try-first"`
	SpaceNearlyFullThresholdPercentPtr *int    `xml:"space-nearly-full-threshold-percent"`
}

func (o *VolumeSpaceAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSpaceAttributesType() *VolumeSpaceAttributesType { return &VolumeSpaceAttributesType{} }

func (o VolumeSpaceAttributesType) String() string {
	var buffer bytes.Buffer
	if o.FilesystemSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "filesystem-size", *o.FilesystemSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("filesystem-size: nil\n"))
	}
	if o.IsFilesysSizeFixedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-filesys-size-fixed", *o.IsFilesysSizeFixedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-filesys-size-fixed: nil\n"))
	}
	if o.IsSpaceGuaranteeEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-space-guarantee-enabled", *o.IsSpaceGuaranteeEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-space-guarantee-enabled: nil\n"))
	}
	if o.OverwriteReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "overwrite-reserve", *o.OverwriteReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("overwrite-reserve: nil\n"))
	}
	if o.OverwriteReserveRequiredPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "overwrite-reserve-required", *o.OverwriteReserveRequiredPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("overwrite-reserve-required: nil\n"))
	}
	if o.OverwriteReserveUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "overwrite-reserve-used", *o.OverwriteReserveUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("overwrite-reserve-used: nil\n"))
	}
	if o.OverwriteReserveUsedActualPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "overwrite-reserve-used-actual", *o.OverwriteReserveUsedActualPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("overwrite-reserve-used-actual: nil\n"))
	}
	if o.PercentageFractionalReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-fractional-reserve", *o.PercentageFractionalReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-fractional-reserve: nil\n"))
	}
	if o.PercentageSizeUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-size-used", *o.PercentageSizeUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-size-used: nil\n"))
	}
	if o.PercentageSnapshotReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-snapshot-reserve", *o.PercentageSnapshotReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-snapshot-reserve: nil\n"))
	}
	if o.PercentageSnapshotReserveUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-snapshot-reserve-used", *o.PercentageSnapshotReserveUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-snapshot-reserve-used: nil\n"))
	}
	if o.PhysicalUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "physical-used", *o.PhysicalUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("physical-used: nil\n"))
	}
	if o.PhysicalUsedPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "physical-used-percent", *o.PhysicalUsedPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("physical-used-percent: nil\n"))
	}
	if o.SizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size", *o.SizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size: nil\n"))
	}
	if o.SizeAvailablePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size-available", *o.SizeAvailablePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size-available: nil\n"))
	}
	if o.SizeAvailableForSnapshotsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size-available-for-snapshots", *o.SizeAvailableForSnapshotsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size-available-for-snapshots: nil\n"))
	}
	if o.SizeTotalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size-total", *o.SizeTotalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size-total: nil\n"))
	}
	if o.SizeUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size-used", *o.SizeUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size-used: nil\n"))
	}
	if o.SizeUsedBySnapshotsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "size-used-by-snapshots", *o.SizeUsedBySnapshotsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("size-used-by-snapshots: nil\n"))
	}
	if o.SnapshotReserveSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-reserve-size", *o.SnapshotReserveSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-reserve-size: nil\n"))
	}
	if o.SpaceFullThresholdPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-full-threshold-percent", *o.SpaceFullThresholdPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-full-threshold-percent: nil\n"))
	}
	if o.SpaceGuaranteePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-guarantee", *o.SpaceGuaranteePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-guarantee: nil\n"))
	}
	if o.SpaceMgmtOptionTryFirstPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-mgmt-option-try-first", *o.SpaceMgmtOptionTryFirstPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-mgmt-option-try-first: nil\n"))
	}
	if o.SpaceNearlyFullThresholdPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "space-nearly-full-threshold-percent", *o.SpaceNearlyFullThresholdPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("space-nearly-full-threshold-percent: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSpaceAttributesType) FilesystemSize() int {
	r := *o.FilesystemSizePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetFilesystemSize(newValue int) *VolumeSpaceAttributesType {
	o.FilesystemSizePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) IsFilesysSizeFixed() bool {
	r := *o.IsFilesysSizeFixedPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetIsFilesysSizeFixed(newValue bool) *VolumeSpaceAttributesType {
	o.IsFilesysSizeFixedPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) IsSpaceGuaranteeEnabled() bool {
	r := *o.IsSpaceGuaranteeEnabledPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetIsSpaceGuaranteeEnabled(newValue bool) *VolumeSpaceAttributesType {
	o.IsSpaceGuaranteeEnabledPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) OverwriteReserve() int {
	r := *o.OverwriteReservePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetOverwriteReserve(newValue int) *VolumeSpaceAttributesType {
	o.OverwriteReservePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) OverwriteReserveRequired() int {
	r := *o.OverwriteReserveRequiredPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetOverwriteReserveRequired(newValue int) *VolumeSpaceAttributesType {
	o.OverwriteReserveRequiredPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) OverwriteReserveUsed() int {
	r := *o.OverwriteReserveUsedPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetOverwriteReserveUsed(newValue int) *VolumeSpaceAttributesType {
	o.OverwriteReserveUsedPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) OverwriteReserveUsedActual() int {
	r := *o.OverwriteReserveUsedActualPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetOverwriteReserveUsedActual(newValue int) *VolumeSpaceAttributesType {
	o.OverwriteReserveUsedActualPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) PercentageFractionalReserve() int {
	r := *o.PercentageFractionalReservePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetPercentageFractionalReserve(newValue int) *VolumeSpaceAttributesType {
	o.PercentageFractionalReservePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) PercentageSizeUsed() int {
	r := *o.PercentageSizeUsedPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetPercentageSizeUsed(newValue int) *VolumeSpaceAttributesType {
	o.PercentageSizeUsedPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) PercentageSnapshotReserve() int {
	r := *o.PercentageSnapshotReservePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetPercentageSnapshotReserve(newValue int) *VolumeSpaceAttributesType {
	o.PercentageSnapshotReservePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) PercentageSnapshotReserveUsed() int {
	r := *o.PercentageSnapshotReserveUsedPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetPercentageSnapshotReserveUsed(newValue int) *VolumeSpaceAttributesType {
	o.PercentageSnapshotReserveUsedPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) PhysicalUsed() int {
	r := *o.PhysicalUsedPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetPhysicalUsed(newValue int) *VolumeSpaceAttributesType {
	o.PhysicalUsedPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) PhysicalUsedPercent() int {
	r := *o.PhysicalUsedPercentPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetPhysicalUsedPercent(newValue int) *VolumeSpaceAttributesType {
	o.PhysicalUsedPercentPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) Size() int {
	r := *o.SizePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSize(newValue int) *VolumeSpaceAttributesType {
	o.SizePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SizeAvailable() int {
	r := *o.SizeAvailablePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSizeAvailable(newValue int) *VolumeSpaceAttributesType {
	o.SizeAvailablePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SizeAvailableForSnapshots() int {
	r := *o.SizeAvailableForSnapshotsPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSizeAvailableForSnapshots(newValue int) *VolumeSpaceAttributesType {
	o.SizeAvailableForSnapshotsPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SizeTotal() int {
	r := *o.SizeTotalPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSizeTotal(newValue int) *VolumeSpaceAttributesType {
	o.SizeTotalPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SizeUsed() int {
	r := *o.SizeUsedPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSizeUsed(newValue int) *VolumeSpaceAttributesType {
	o.SizeUsedPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SizeUsedBySnapshots() int {
	r := *o.SizeUsedBySnapshotsPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSizeUsedBySnapshots(newValue int) *VolumeSpaceAttributesType {
	o.SizeUsedBySnapshotsPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SnapshotReserveSize() int {
	r := *o.SnapshotReserveSizePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSnapshotReserveSize(newValue int) *VolumeSpaceAttributesType {
	o.SnapshotReserveSizePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SpaceFullThresholdPercent() int {
	r := *o.SpaceFullThresholdPercentPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSpaceFullThresholdPercent(newValue int) *VolumeSpaceAttributesType {
	o.SpaceFullThresholdPercentPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SpaceGuarantee() string {
	r := *o.SpaceGuaranteePtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSpaceGuarantee(newValue string) *VolumeSpaceAttributesType {
	o.SpaceGuaranteePtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SpaceMgmtOptionTryFirst() string {
	r := *o.SpaceMgmtOptionTryFirstPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSpaceMgmtOptionTryFirst(newValue string) *VolumeSpaceAttributesType {
	o.SpaceMgmtOptionTryFirstPtr = &newValue
	return o
}

func (o *VolumeSpaceAttributesType) SpaceNearlyFullThresholdPercent() int {
	r := *o.SpaceNearlyFullThresholdPercentPtr
	return r
}

func (o *VolumeSpaceAttributesType) SetSpaceNearlyFullThresholdPercent(newValue int) *VolumeSpaceAttributesType {
	o.SpaceNearlyFullThresholdPercentPtr = &newValue
	return o
}

type VolumeSnapshotAutodeleteAttributesType struct {
	XMLName xml.Name `xml:"volume-snapshot-autodelete-attributes"`

	CommitmentPtr          *string `xml:"commitment"`
	DeferDeletePtr         *string `xml:"defer-delete"`
	DeleteOrderPtr         *string `xml:"delete-order"`
	DestroyListPtr         *string `xml:"destroy-list"`
	IsAutodeleteEnabledPtr *bool   `xml:"is-autodelete-enabled"`
	PrefixPtr              *string `xml:"prefix"`
	TargetFreeSpacePtr     *int    `xml:"target-free-space"`
	TriggerPtr             *string `xml:"trigger"`
}

func (o *VolumeSnapshotAutodeleteAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSnapshotAutodeleteAttributesType() *VolumeSnapshotAutodeleteAttributesType {
	return &VolumeSnapshotAutodeleteAttributesType{}
}

func (o VolumeSnapshotAutodeleteAttributesType) String() string {
	var buffer bytes.Buffer
	if o.CommitmentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "commitment", *o.CommitmentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("commitment: nil\n"))
	}
	if o.DeferDeletePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "defer-delete", *o.DeferDeletePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("defer-delete: nil\n"))
	}
	if o.DeleteOrderPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "delete-order", *o.DeleteOrderPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("delete-order: nil\n"))
	}
	if o.DestroyListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "destroy-list", *o.DestroyListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("destroy-list: nil\n"))
	}
	if o.IsAutodeleteEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-autodelete-enabled", *o.IsAutodeleteEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-autodelete-enabled: nil\n"))
	}
	if o.PrefixPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "prefix", *o.PrefixPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("prefix: nil\n"))
	}
	if o.TargetFreeSpacePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "target-free-space", *o.TargetFreeSpacePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("target-free-space: nil\n"))
	}
	if o.TriggerPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "trigger", *o.TriggerPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("trigger: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSnapshotAutodeleteAttributesType) Commitment() string {
	r := *o.CommitmentPtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetCommitment(newValue string) *VolumeSnapshotAutodeleteAttributesType {
	o.CommitmentPtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) DeferDelete() string {
	r := *o.DeferDeletePtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetDeferDelete(newValue string) *VolumeSnapshotAutodeleteAttributesType {
	o.DeferDeletePtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) DeleteOrder() string {
	r := *o.DeleteOrderPtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetDeleteOrder(newValue string) *VolumeSnapshotAutodeleteAttributesType {
	o.DeleteOrderPtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) DestroyList() string {
	r := *o.DestroyListPtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetDestroyList(newValue string) *VolumeSnapshotAutodeleteAttributesType {
	o.DestroyListPtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) IsAutodeleteEnabled() bool {
	r := *o.IsAutodeleteEnabledPtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetIsAutodeleteEnabled(newValue bool) *VolumeSnapshotAutodeleteAttributesType {
	o.IsAutodeleteEnabledPtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) Prefix() string {
	r := *o.PrefixPtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetPrefix(newValue string) *VolumeSnapshotAutodeleteAttributesType {
	o.PrefixPtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) TargetFreeSpace() int {
	r := *o.TargetFreeSpacePtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetTargetFreeSpace(newValue int) *VolumeSnapshotAutodeleteAttributesType {
	o.TargetFreeSpacePtr = &newValue
	return o
}

func (o *VolumeSnapshotAutodeleteAttributesType) Trigger() string {
	r := *o.TriggerPtr
	return r
}

func (o *VolumeSnapshotAutodeleteAttributesType) SetTrigger(newValue string) *VolumeSnapshotAutodeleteAttributesType {
	o.TriggerPtr = &newValue
	return o
}

type VolumeSnapshotAttributesType struct {
	XMLName xml.Name `xml:"volume-snapshot-attributes"`

	AutoSnapshotsEnabledPtr           *bool   `xml:"auto-snapshots-enabled"`
	SnapdirAccessEnabledPtr           *bool   `xml:"snapdir-access-enabled"`
	SnapshotCloneDependencyEnabledPtr *bool   `xml:"snapshot-clone-dependency-enabled"`
	SnapshotCountPtr                  *int    `xml:"snapshot-count"`
	SnapshotPolicyPtr                 *string `xml:"snapshot-policy"`
}

func (o *VolumeSnapshotAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSnapshotAttributesType() *VolumeSnapshotAttributesType {
	return &VolumeSnapshotAttributesType{}
}

func (o VolumeSnapshotAttributesType) String() string {
	var buffer bytes.Buffer
	if o.AutoSnapshotsEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "auto-snapshots-enabled", *o.AutoSnapshotsEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("auto-snapshots-enabled: nil\n"))
	}
	if o.SnapdirAccessEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapdir-access-enabled", *o.SnapdirAccessEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapdir-access-enabled: nil\n"))
	}
	if o.SnapshotCloneDependencyEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-clone-dependency-enabled", *o.SnapshotCloneDependencyEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-clone-dependency-enabled: nil\n"))
	}
	if o.SnapshotCountPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-count", *o.SnapshotCountPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-count: nil\n"))
	}
	if o.SnapshotPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-policy", *o.SnapshotPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-policy: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSnapshotAttributesType) AutoSnapshotsEnabled() bool {
	r := *o.AutoSnapshotsEnabledPtr
	return r
}

func (o *VolumeSnapshotAttributesType) SetAutoSnapshotsEnabled(newValue bool) *VolumeSnapshotAttributesType {
	o.AutoSnapshotsEnabledPtr = &newValue
	return o
}

func (o *VolumeSnapshotAttributesType) SnapdirAccessEnabled() bool {
	r := *o.SnapdirAccessEnabledPtr
	return r
}

func (o *VolumeSnapshotAttributesType) SetSnapdirAccessEnabled(newValue bool) *VolumeSnapshotAttributesType {
	o.SnapdirAccessEnabledPtr = &newValue
	return o
}

func (o *VolumeSnapshotAttributesType) SnapshotCloneDependencyEnabled() bool {
	r := *o.SnapshotCloneDependencyEnabledPtr
	return r
}

func (o *VolumeSnapshotAttributesType) SetSnapshotCloneDependencyEnabled(newValue bool) *VolumeSnapshotAttributesType {
	o.SnapshotCloneDependencyEnabledPtr = &newValue
	return o
}

func (o *VolumeSnapshotAttributesType) SnapshotCount() int {
	r := *o.SnapshotCountPtr
	return r
}

func (o *VolumeSnapshotAttributesType) SetSnapshotCount(newValue int) *VolumeSnapshotAttributesType {
	o.SnapshotCountPtr = &newValue
	return o
}

func (o *VolumeSnapshotAttributesType) SnapshotPolicy() string {
	r := *o.SnapshotPolicyPtr
	return r
}

func (o *VolumeSnapshotAttributesType) SetSnapshotPolicy(newValue string) *VolumeSnapshotAttributesType {
	o.SnapshotPolicyPtr = &newValue
	return o
}

type VolumeSisAttributesType struct {
	XMLName xml.Name `xml:"volume-sis-attributes"`

	CompressionSpaceSavedPtr             *int      `xml:"compression-space-saved"`
	DeduplicationSpaceSavedPtr           *int      `xml:"deduplication-space-saved"`
	DeduplicationSpaceSharedPtr          *SizeType `xml:"deduplication-space-shared"`
	IsSisLoggingEnabledPtr               *bool     `xml:"is-sis-logging-enabled"`
	IsSisVolumePtr                       *bool     `xml:"is-sis-volume"`
	PercentageCompressionSpaceSavedPtr   *int      `xml:"percentage-compression-space-saved"`
	PercentageDeduplicationSpaceSavedPtr *int      `xml:"percentage-deduplication-space-saved"`
	PercentageTotalSpaceSavedPtr         *int      `xml:"percentage-total-space-saved"`
	TotalSpaceSavedPtr                   *int      `xml:"total-space-saved"`
}

func (o *VolumeSisAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSisAttributesType() *VolumeSisAttributesType { return &VolumeSisAttributesType{} }

func (o VolumeSisAttributesType) String() string {
	var buffer bytes.Buffer
	if o.CompressionSpaceSavedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "compression-space-saved", *o.CompressionSpaceSavedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("compression-space-saved: nil\n"))
	}
	if o.DeduplicationSpaceSavedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "deduplication-space-saved", *o.DeduplicationSpaceSavedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("deduplication-space-saved: nil\n"))
	}
	if o.DeduplicationSpaceSharedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "deduplication-space-shared", *o.DeduplicationSpaceSharedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("deduplication-space-shared: nil\n"))
	}
	if o.IsSisLoggingEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-sis-logging-enabled", *o.IsSisLoggingEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-sis-logging-enabled: nil\n"))
	}
	if o.IsSisVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-sis-volume", *o.IsSisVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-sis-volume: nil\n"))
	}
	if o.PercentageCompressionSpaceSavedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-compression-space-saved", *o.PercentageCompressionSpaceSavedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-compression-space-saved: nil\n"))
	}
	if o.PercentageDeduplicationSpaceSavedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-deduplication-space-saved", *o.PercentageDeduplicationSpaceSavedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-deduplication-space-saved: nil\n"))
	}
	if o.PercentageTotalSpaceSavedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-total-space-saved", *o.PercentageTotalSpaceSavedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-total-space-saved: nil\n"))
	}
	if o.TotalSpaceSavedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "total-space-saved", *o.TotalSpaceSavedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("total-space-saved: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSisAttributesType) CompressionSpaceSaved() int {
	r := *o.CompressionSpaceSavedPtr
	return r
}

func (o *VolumeSisAttributesType) SetCompressionSpaceSaved(newValue int) *VolumeSisAttributesType {
	o.CompressionSpaceSavedPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) DeduplicationSpaceSaved() int {
	r := *o.DeduplicationSpaceSavedPtr
	return r
}

func (o *VolumeSisAttributesType) SetDeduplicationSpaceSaved(newValue int) *VolumeSisAttributesType {
	o.DeduplicationSpaceSavedPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) DeduplicationSpaceShared() SizeType {
	r := *o.DeduplicationSpaceSharedPtr
	return r
}

func (o *VolumeSisAttributesType) SetDeduplicationSpaceShared(newValue SizeType) *VolumeSisAttributesType {
	o.DeduplicationSpaceSharedPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) IsSisLoggingEnabled() bool {
	r := *o.IsSisLoggingEnabledPtr
	return r
}

func (o *VolumeSisAttributesType) SetIsSisLoggingEnabled(newValue bool) *VolumeSisAttributesType {
	o.IsSisLoggingEnabledPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) IsSisVolume() bool {
	r := *o.IsSisVolumePtr
	return r
}

func (o *VolumeSisAttributesType) SetIsSisVolume(newValue bool) *VolumeSisAttributesType {
	o.IsSisVolumePtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) PercentageCompressionSpaceSaved() int {
	r := *o.PercentageCompressionSpaceSavedPtr
	return r
}

func (o *VolumeSisAttributesType) SetPercentageCompressionSpaceSaved(newValue int) *VolumeSisAttributesType {
	o.PercentageCompressionSpaceSavedPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) PercentageDeduplicationSpaceSaved() int {
	r := *o.PercentageDeduplicationSpaceSavedPtr
	return r
}

func (o *VolumeSisAttributesType) SetPercentageDeduplicationSpaceSaved(newValue int) *VolumeSisAttributesType {
	o.PercentageDeduplicationSpaceSavedPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) PercentageTotalSpaceSaved() int {
	r := *o.PercentageTotalSpaceSavedPtr
	return r
}

func (o *VolumeSisAttributesType) SetPercentageTotalSpaceSaved(newValue int) *VolumeSisAttributesType {
	o.PercentageTotalSpaceSavedPtr = &newValue
	return o
}

func (o *VolumeSisAttributesType) TotalSpaceSaved() int {
	r := *o.TotalSpaceSavedPtr
	return r
}

func (o *VolumeSisAttributesType) SetTotalSpaceSaved(newValue int) *VolumeSisAttributesType {
	o.TotalSpaceSavedPtr = &newValue
	return o
}

type VolumeSecurityUnixAttributesType struct {
	XMLName xml.Name `xml:"volume-security-unix-attributes"`

	GroupIdPtr     *int    `xml:"group-id"`
	PermissionsPtr *string `xml:"permissions"`
	UserIdPtr      *int    `xml:"user-id"`
}

func (o *VolumeSecurityUnixAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSecurityUnixAttributesType() *VolumeSecurityUnixAttributesType {
	return &VolumeSecurityUnixAttributesType{}
}

func (o VolumeSecurityUnixAttributesType) String() string {
	var buffer bytes.Buffer
	if o.GroupIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "group-id", *o.GroupIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("group-id: nil\n"))
	}
	if o.PermissionsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "permissions", *o.PermissionsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("permissions: nil\n"))
	}
	if o.UserIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "user-id", *o.UserIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("user-id: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSecurityUnixAttributesType) GroupId() int {
	r := *o.GroupIdPtr
	return r
}

func (o *VolumeSecurityUnixAttributesType) SetGroupId(newValue int) *VolumeSecurityUnixAttributesType {
	o.GroupIdPtr = &newValue
	return o
}

func (o *VolumeSecurityUnixAttributesType) Permissions() string {
	r := *o.PermissionsPtr
	return r
}

func (o *VolumeSecurityUnixAttributesType) SetPermissions(newValue string) *VolumeSecurityUnixAttributesType {
	o.PermissionsPtr = &newValue
	return o
}

func (o *VolumeSecurityUnixAttributesType) UserId() int {
	r := *o.UserIdPtr
	return r
}

func (o *VolumeSecurityUnixAttributesType) SetUserId(newValue int) *VolumeSecurityUnixAttributesType {
	o.UserIdPtr = &newValue
	return o
}

type VolumeSecurityAttributesType struct {
	XMLName xml.Name `xml:"volume-security-attributes"`

	StylePtr                        *string                           `xml:"style"`
	VolumeSecurityUnixAttributesPtr *VolumeSecurityUnixAttributesType `xml:"volume-security-unix-attributes"`
}

func (o *VolumeSecurityAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeSecurityAttributesType() *VolumeSecurityAttributesType {
	return &VolumeSecurityAttributesType{}
}

func (o VolumeSecurityAttributesType) String() string {
	var buffer bytes.Buffer
	if o.StylePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "style", *o.StylePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("style: nil\n"))
	}
	if o.VolumeSecurityUnixAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-security-unix-attributes", *o.VolumeSecurityUnixAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-security-unix-attributes: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeSecurityAttributesType) Style() string {
	r := *o.StylePtr
	return r
}

func (o *VolumeSecurityAttributesType) SetStyle(newValue string) *VolumeSecurityAttributesType {
	o.StylePtr = &newValue
	return o
}

func (o *VolumeSecurityAttributesType) VolumeSecurityUnixAttributes() VolumeSecurityUnixAttributesType {
	r := *o.VolumeSecurityUnixAttributesPtr
	return r
}

func (o *VolumeSecurityAttributesType) SetVolumeSecurityUnixAttributes(newValue VolumeSecurityUnixAttributesType) *VolumeSecurityAttributesType {
	o.VolumeSecurityUnixAttributesPtr = &newValue
	return o
}

type VolumeQosAttributesType struct {
	XMLName xml.Name `xml:"volume-qos-attributes"`

	PolicyGroupNamePtr *string `xml:"policy-group-name"`
}

func (o *VolumeQosAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeQosAttributesType() *VolumeQosAttributesType { return &VolumeQosAttributesType{} }

func (o VolumeQosAttributesType) String() string {
	var buffer bytes.Buffer
	if o.PolicyGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "policy-group-name", *o.PolicyGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("policy-group-name: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeQosAttributesType) PolicyGroupName() string {
	r := *o.PolicyGroupNamePtr
	return r
}

func (o *VolumeQosAttributesType) SetPolicyGroupName(newValue string) *VolumeQosAttributesType {
	o.PolicyGroupNamePtr = &newValue
	return o
}

type VolumePerformanceAttributesType struct {
	XMLName xml.Name `xml:"volume-performance-attributes"`

	ExtentEnabledPtr        *string `xml:"extent-enabled"`
	FcDelegsEnabledPtr      *bool   `xml:"fc-delegs-enabled"`
	IsAtimeUpdateEnabledPtr *bool   `xml:"is-atime-update-enabled"`
	MaxWriteAllocBlocksPtr  *int    `xml:"max-write-alloc-blocks"`
	MinimalReadAheadPtr     *bool   `xml:"minimal-read-ahead"`
	ReadReallocPtr          *string `xml:"read-realloc"`
}

func (o *VolumePerformanceAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumePerformanceAttributesType() *VolumePerformanceAttributesType {
	return &VolumePerformanceAttributesType{}
}

func (o VolumePerformanceAttributesType) String() string {
	var buffer bytes.Buffer
	if o.ExtentEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "extent-enabled", *o.ExtentEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("extent-enabled: nil\n"))
	}
	if o.FcDelegsEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "fc-delegs-enabled", *o.FcDelegsEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("fc-delegs-enabled: nil\n"))
	}
	if o.IsAtimeUpdateEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-atime-update-enabled", *o.IsAtimeUpdateEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-atime-update-enabled: nil\n"))
	}
	if o.MaxWriteAllocBlocksPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-write-alloc-blocks", *o.MaxWriteAllocBlocksPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-write-alloc-blocks: nil\n"))
	}
	if o.MinimalReadAheadPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "minimal-read-ahead", *o.MinimalReadAheadPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("minimal-read-ahead: nil\n"))
	}
	if o.ReadReallocPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "read-realloc", *o.ReadReallocPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("read-realloc: nil\n"))
	}
	return buffer.String()
}

func (o *VolumePerformanceAttributesType) ExtentEnabled() string {
	r := *o.ExtentEnabledPtr
	return r
}

func (o *VolumePerformanceAttributesType) SetExtentEnabled(newValue string) *VolumePerformanceAttributesType {
	o.ExtentEnabledPtr = &newValue
	return o
}

func (o *VolumePerformanceAttributesType) FcDelegsEnabled() bool {
	r := *o.FcDelegsEnabledPtr
	return r
}

func (o *VolumePerformanceAttributesType) SetFcDelegsEnabled(newValue bool) *VolumePerformanceAttributesType {
	o.FcDelegsEnabledPtr = &newValue
	return o
}

func (o *VolumePerformanceAttributesType) IsAtimeUpdateEnabled() bool {
	r := *o.IsAtimeUpdateEnabledPtr
	return r
}

func (o *VolumePerformanceAttributesType) SetIsAtimeUpdateEnabled(newValue bool) *VolumePerformanceAttributesType {
	o.IsAtimeUpdateEnabledPtr = &newValue
	return o
}

func (o *VolumePerformanceAttributesType) MaxWriteAllocBlocks() int {
	r := *o.MaxWriteAllocBlocksPtr
	return r
}

func (o *VolumePerformanceAttributesType) SetMaxWriteAllocBlocks(newValue int) *VolumePerformanceAttributesType {
	o.MaxWriteAllocBlocksPtr = &newValue
	return o
}

func (o *VolumePerformanceAttributesType) MinimalReadAhead() bool {
	r := *o.MinimalReadAheadPtr
	return r
}

func (o *VolumePerformanceAttributesType) SetMinimalReadAhead(newValue bool) *VolumePerformanceAttributesType {
	o.MinimalReadAheadPtr = &newValue
	return o
}

func (o *VolumePerformanceAttributesType) ReadRealloc() string {
	r := *o.ReadReallocPtr
	return r
}

func (o *VolumePerformanceAttributesType) SetReadRealloc(newValue string) *VolumePerformanceAttributesType {
	o.ReadReallocPtr = &newValue
	return o
}

type VolumeMirrorAttributesType struct {
	XMLName xml.Name `xml:"volume-mirror-attributes"`

	IsDataProtectionMirrorPtr   *bool `xml:"is-data-protection-mirror"`
	IsLoadSharingMirrorPtr      *bool `xml:"is-load-sharing-mirror"`
	IsMoveMirrorPtr             *bool `xml:"is-move-mirror"`
	IsReplicaVolumePtr          *bool `xml:"is-replica-volume"`
	MirrorTransferInProgressPtr *bool `xml:"mirror-transfer-in-progress"`
	RedirectSnapshotIdPtr       *int  `xml:"redirect-snapshot-id"`
}

func (o *VolumeMirrorAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeMirrorAttributesType() *VolumeMirrorAttributesType { return &VolumeMirrorAttributesType{} }

func (o VolumeMirrorAttributesType) String() string {
	var buffer bytes.Buffer
	if o.IsDataProtectionMirrorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-data-protection-mirror", *o.IsDataProtectionMirrorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-data-protection-mirror: nil\n"))
	}
	if o.IsLoadSharingMirrorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-load-sharing-mirror", *o.IsLoadSharingMirrorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-load-sharing-mirror: nil\n"))
	}
	if o.IsMoveMirrorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-move-mirror", *o.IsMoveMirrorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-move-mirror: nil\n"))
	}
	if o.IsReplicaVolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-replica-volume", *o.IsReplicaVolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-replica-volume: nil\n"))
	}
	if o.MirrorTransferInProgressPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "mirror-transfer-in-progress", *o.MirrorTransferInProgressPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("mirror-transfer-in-progress: nil\n"))
	}
	if o.RedirectSnapshotIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "redirect-snapshot-id", *o.RedirectSnapshotIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("redirect-snapshot-id: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeMirrorAttributesType) IsDataProtectionMirror() bool {
	r := *o.IsDataProtectionMirrorPtr
	return r
}

func (o *VolumeMirrorAttributesType) SetIsDataProtectionMirror(newValue bool) *VolumeMirrorAttributesType {
	o.IsDataProtectionMirrorPtr = &newValue
	return o
}

func (o *VolumeMirrorAttributesType) IsLoadSharingMirror() bool {
	r := *o.IsLoadSharingMirrorPtr
	return r
}

func (o *VolumeMirrorAttributesType) SetIsLoadSharingMirror(newValue bool) *VolumeMirrorAttributesType {
	o.IsLoadSharingMirrorPtr = &newValue
	return o
}

func (o *VolumeMirrorAttributesType) IsMoveMirror() bool {
	r := *o.IsMoveMirrorPtr
	return r
}

func (o *VolumeMirrorAttributesType) SetIsMoveMirror(newValue bool) *VolumeMirrorAttributesType {
	o.IsMoveMirrorPtr = &newValue
	return o
}

func (o *VolumeMirrorAttributesType) IsReplicaVolume() bool {
	r := *o.IsReplicaVolumePtr
	return r
}

func (o *VolumeMirrorAttributesType) SetIsReplicaVolume(newValue bool) *VolumeMirrorAttributesType {
	o.IsReplicaVolumePtr = &newValue
	return o
}

func (o *VolumeMirrorAttributesType) MirrorTransferInProgress() bool {
	r := *o.MirrorTransferInProgressPtr
	return r
}

func (o *VolumeMirrorAttributesType) SetMirrorTransferInProgress(newValue bool) *VolumeMirrorAttributesType {
	o.MirrorTransferInProgressPtr = &newValue
	return o
}

func (o *VolumeMirrorAttributesType) RedirectSnapshotId() int {
	r := *o.RedirectSnapshotIdPtr
	return r
}

func (o *VolumeMirrorAttributesType) SetRedirectSnapshotId(newValue int) *VolumeMirrorAttributesType {
	o.RedirectSnapshotIdPtr = &newValue
	return o
}

type LanguageCodeType string

type VolumeLanguageAttributesType struct {
	XMLName xml.Name `xml:"volume-language-attributes"`

	IsConvertUcodeEnabledPtr *bool             `xml:"is-convert-ucode-enabled"`
	IsCreateUcodeEnabledPtr  *bool             `xml:"is-create-ucode-enabled"`
	LanguagePtr              *string           `xml:"language"`
	LanguageCodePtr          *LanguageCodeType `xml:"language-code"`
	NfsCharacterSetPtr       *string           `xml:"nfs-character-set"`
	OemCharacterSetPtr       *string           `xml:"oem-character-set"`
}

func (o *VolumeLanguageAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeLanguageAttributesType() *VolumeLanguageAttributesType {
	return &VolumeLanguageAttributesType{}
}

func (o VolumeLanguageAttributesType) String() string {
	var buffer bytes.Buffer
	if o.IsConvertUcodeEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-convert-ucode-enabled", *o.IsConvertUcodeEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-convert-ucode-enabled: nil\n"))
	}
	if o.IsCreateUcodeEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-create-ucode-enabled", *o.IsCreateUcodeEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-create-ucode-enabled: nil\n"))
	}
	if o.LanguagePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "language", *o.LanguagePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("language: nil\n"))
	}
	if o.LanguageCodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "language-code", *o.LanguageCodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("language-code: nil\n"))
	}
	if o.NfsCharacterSetPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "nfs-character-set", *o.NfsCharacterSetPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("nfs-character-set: nil\n"))
	}
	if o.OemCharacterSetPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "oem-character-set", *o.OemCharacterSetPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("oem-character-set: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeLanguageAttributesType) IsConvertUcodeEnabled() bool {
	r := *o.IsConvertUcodeEnabledPtr
	return r
}

func (o *VolumeLanguageAttributesType) SetIsConvertUcodeEnabled(newValue bool) *VolumeLanguageAttributesType {
	o.IsConvertUcodeEnabledPtr = &newValue
	return o
}

func (o *VolumeLanguageAttributesType) IsCreateUcodeEnabled() bool {
	r := *o.IsCreateUcodeEnabledPtr
	return r
}

func (o *VolumeLanguageAttributesType) SetIsCreateUcodeEnabled(newValue bool) *VolumeLanguageAttributesType {
	o.IsCreateUcodeEnabledPtr = &newValue
	return o
}

func (o *VolumeLanguageAttributesType) Language() string {
	r := *o.LanguagePtr
	return r
}

func (o *VolumeLanguageAttributesType) SetLanguage(newValue string) *VolumeLanguageAttributesType {
	o.LanguagePtr = &newValue
	return o
}

func (o *VolumeLanguageAttributesType) LanguageCode() LanguageCodeType {
	r := *o.LanguageCodePtr
	return r
}

func (o *VolumeLanguageAttributesType) SetLanguageCode(newValue LanguageCodeType) *VolumeLanguageAttributesType {
	o.LanguageCodePtr = &newValue
	return o
}

func (o *VolumeLanguageAttributesType) NfsCharacterSet() string {
	r := *o.NfsCharacterSetPtr
	return r
}

func (o *VolumeLanguageAttributesType) SetNfsCharacterSet(newValue string) *VolumeLanguageAttributesType {
	o.NfsCharacterSetPtr = &newValue
	return o
}

func (o *VolumeLanguageAttributesType) OemCharacterSet() string {
	r := *o.OemCharacterSetPtr
	return r
}

func (o *VolumeLanguageAttributesType) SetOemCharacterSet(newValue string) *VolumeLanguageAttributesType {
	o.OemCharacterSetPtr = &newValue
	return o
}

type VolumeInodeAttributesType struct {
	XMLName xml.Name `xml:"volume-inode-attributes"`

	BlockTypePtr                *string `xml:"block-type"`
	FilesPrivateUsedPtr         *int    `xml:"files-private-used"`
	FilesTotalPtr               *int    `xml:"files-total"`
	FilesUsedPtr                *int    `xml:"files-used"`
	InodefilePrivateCapacityPtr *int    `xml:"inodefile-private-capacity"`
	InodefilePublicCapacityPtr  *int    `xml:"inodefile-public-capacity"`
}

func (o *VolumeInodeAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeInodeAttributesType() *VolumeInodeAttributesType { return &VolumeInodeAttributesType{} }

func (o VolumeInodeAttributesType) String() string {
	var buffer bytes.Buffer
	if o.BlockTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "block-type", *o.BlockTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("block-type: nil\n"))
	}
	if o.FilesPrivateUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "files-private-used", *o.FilesPrivateUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("files-private-used: nil\n"))
	}
	if o.FilesTotalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "files-total", *o.FilesTotalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("files-total: nil\n"))
	}
	if o.FilesUsedPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "files-used", *o.FilesUsedPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("files-used: nil\n"))
	}
	if o.InodefilePrivateCapacityPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "inodefile-private-capacity", *o.InodefilePrivateCapacityPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("inodefile-private-capacity: nil\n"))
	}
	if o.InodefilePublicCapacityPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "inodefile-public-capacity", *o.InodefilePublicCapacityPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("inodefile-public-capacity: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeInodeAttributesType) BlockType() string {
	r := *o.BlockTypePtr
	return r
}

func (o *VolumeInodeAttributesType) SetBlockType(newValue string) *VolumeInodeAttributesType {
	o.BlockTypePtr = &newValue
	return o
}

func (o *VolumeInodeAttributesType) FilesPrivateUsed() int {
	r := *o.FilesPrivateUsedPtr
	return r
}

func (o *VolumeInodeAttributesType) SetFilesPrivateUsed(newValue int) *VolumeInodeAttributesType {
	o.FilesPrivateUsedPtr = &newValue
	return o
}

func (o *VolumeInodeAttributesType) FilesTotal() int {
	r := *o.FilesTotalPtr
	return r
}

func (o *VolumeInodeAttributesType) SetFilesTotal(newValue int) *VolumeInodeAttributesType {
	o.FilesTotalPtr = &newValue
	return o
}

func (o *VolumeInodeAttributesType) FilesUsed() int {
	r := *o.FilesUsedPtr
	return r
}

func (o *VolumeInodeAttributesType) SetFilesUsed(newValue int) *VolumeInodeAttributesType {
	o.FilesUsedPtr = &newValue
	return o
}

func (o *VolumeInodeAttributesType) InodefilePrivateCapacity() int {
	r := *o.InodefilePrivateCapacityPtr
	return r
}

func (o *VolumeInodeAttributesType) SetInodefilePrivateCapacity(newValue int) *VolumeInodeAttributesType {
	o.InodefilePrivateCapacityPtr = &newValue
	return o
}

func (o *VolumeInodeAttributesType) InodefilePublicCapacity() int {
	r := *o.InodefilePublicCapacityPtr
	return r
}

func (o *VolumeInodeAttributesType) SetInodefilePublicCapacity(newValue int) *VolumeInodeAttributesType {
	o.InodefilePublicCapacityPtr = &newValue
	return o
}

type AggrNameType string

type ReposConstituentRoleType string

type VolumeInfinitevolAttributesType struct {
	XMLName xml.Name `xml:"volume-infinitevol-attributes"`

	ConstituentRolePtr             *ReposConstituentRoleType `xml:"constituent-role"`
	EnableSnapdiffPtr              *bool                     `xml:"enable-snapdiff"`
	IsManagedByServicePtr          *bool                     `xml:"is-managed-by-service"`
	MaxDataConstituentSizePtr      *SizeType                 `xml:"max-data-constituent-size"`
	MaxNamespaceConstituentSizePtr *SizeType                 `xml:"max-namespace-constituent-size"`
	NamespaceMirrorAggrListPtr     []AggrNameType            `xml:"namespace-mirror-aggr-list>aggr-name"`
	StorageServicePtr              *string                   `xml:"storage-service"`
}

func (o *VolumeInfinitevolAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeInfinitevolAttributesType() *VolumeInfinitevolAttributesType {
	return &VolumeInfinitevolAttributesType{}
}

func (o VolumeInfinitevolAttributesType) String() string {
	var buffer bytes.Buffer
	if o.ConstituentRolePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "constituent-role", *o.ConstituentRolePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("constituent-role: nil\n"))
	}
	if o.EnableSnapdiffPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "enable-snapdiff", *o.EnableSnapdiffPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("enable-snapdiff: nil\n"))
	}
	if o.IsManagedByServicePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-managed-by-service", *o.IsManagedByServicePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-managed-by-service: nil\n"))
	}
	if o.MaxDataConstituentSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-data-constituent-size", *o.MaxDataConstituentSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-data-constituent-size: nil\n"))
	}
	if o.MaxNamespaceConstituentSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-namespace-constituent-size", *o.MaxNamespaceConstituentSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-namespace-constituent-size: nil\n"))
	}
	if o.NamespaceMirrorAggrListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "namespace-mirror-aggr-list", o.NamespaceMirrorAggrListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("namespace-mirror-aggr-list: nil\n"))
	}
	if o.StorageServicePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "storage-service", *o.StorageServicePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("storage-service: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeInfinitevolAttributesType) ConstituentRole() ReposConstituentRoleType {
	r := *o.ConstituentRolePtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetConstituentRole(newValue ReposConstituentRoleType) *VolumeInfinitevolAttributesType {
	o.ConstituentRolePtr = &newValue
	return o
}

func (o *VolumeInfinitevolAttributesType) EnableSnapdiff() bool {
	r := *o.EnableSnapdiffPtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetEnableSnapdiff(newValue bool) *VolumeInfinitevolAttributesType {
	o.EnableSnapdiffPtr = &newValue
	return o
}

func (o *VolumeInfinitevolAttributesType) IsManagedByService() bool {
	r := *o.IsManagedByServicePtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetIsManagedByService(newValue bool) *VolumeInfinitevolAttributesType {
	o.IsManagedByServicePtr = &newValue
	return o
}

func (o *VolumeInfinitevolAttributesType) MaxDataConstituentSize() SizeType {
	r := *o.MaxDataConstituentSizePtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetMaxDataConstituentSize(newValue SizeType) *VolumeInfinitevolAttributesType {
	o.MaxDataConstituentSizePtr = &newValue
	return o
}

func (o *VolumeInfinitevolAttributesType) MaxNamespaceConstituentSize() SizeType {
	r := *o.MaxNamespaceConstituentSizePtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetMaxNamespaceConstituentSize(newValue SizeType) *VolumeInfinitevolAttributesType {
	o.MaxNamespaceConstituentSizePtr = &newValue
	return o
}

func (o *VolumeInfinitevolAttributesType) NamespaceMirrorAggrList() []AggrNameType {
	r := o.NamespaceMirrorAggrListPtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetNamespaceMirrorAggrList(newValue []AggrNameType) *VolumeInfinitevolAttributesType {
	newSlice := make([]AggrNameType, len(newValue))
	copy(newSlice, newValue)
	o.NamespaceMirrorAggrListPtr = newSlice
	return o
}

func (o *VolumeInfinitevolAttributesType) StorageService() string {
	r := *o.StorageServicePtr
	return r
}

func (o *VolumeInfinitevolAttributesType) SetStorageService(newValue string) *VolumeInfinitevolAttributesType {
	o.StorageServicePtr = &newValue
	return o
}

type JunctionPathType string

type VolumeIdAttributesType struct {
	XMLName xml.Name `xml:"volume-id-attributes"`

	CommentPtr                 *string           `xml:"comment"`
	ContainingAggregateNamePtr *string           `xml:"containing-aggregate-name"`
	ContainingAggregateUuidPtr *UuidType         `xml:"containing-aggregate-uuid"`
	CreationTimePtr            *int              `xml:"creation-time"`
	DsidPtr                    *int              `xml:"dsid"`
	FsidPtr                    *string           `xml:"fsid"`
	InstanceUuidPtr            *UuidType         `xml:"instance-uuid"`
	JunctionParentNamePtr      *VolumeNameType   `xml:"junction-parent-name"`
	JunctionPathPtr            *JunctionPathType `xml:"junction-path"`
	MsidPtr                    *int              `xml:"msid"`
	NamePtr                    *VolumeNameType   `xml:"name"`
	NameOrdinalPtr             *string           `xml:"name-ordinal"`
	NodePtr                    *NodeNameType     `xml:"node"`
	OwningVserverNamePtr       *string           `xml:"owning-vserver-name"`
	OwningVserverUuidPtr       *UuidType         `xml:"owning-vserver-uuid"`
	ProvenanceUuidPtr          *UuidType         `xml:"provenance-uuid"`
	StylePtr                   *string           `xml:"style"`
	TypePtr                    *string           `xml:"type"`
	UuidPtr                    *UuidType         `xml:"uuid"`
}

func (o *VolumeIdAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeIdAttributesType() *VolumeIdAttributesType { return &VolumeIdAttributesType{} }

func (o VolumeIdAttributesType) String() string {
	var buffer bytes.Buffer
	if o.CommentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "comment", *o.CommentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("comment: nil\n"))
	}
	if o.ContainingAggregateNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "containing-aggregate-name", *o.ContainingAggregateNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("containing-aggregate-name: nil\n"))
	}
	if o.ContainingAggregateUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "containing-aggregate-uuid", *o.ContainingAggregateUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("containing-aggregate-uuid: nil\n"))
	}
	if o.CreationTimePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "creation-time", *o.CreationTimePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("creation-time: nil\n"))
	}
	if o.DsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "dsid", *o.DsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("dsid: nil\n"))
	}
	if o.FsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "fsid", *o.FsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("fsid: nil\n"))
	}
	if o.InstanceUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "instance-uuid", *o.InstanceUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("instance-uuid: nil\n"))
	}
	if o.JunctionParentNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "junction-parent-name", *o.JunctionParentNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("junction-parent-name: nil\n"))
	}
	if o.JunctionPathPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "junction-path", *o.JunctionPathPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("junction-path: nil\n"))
	}
	if o.MsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "msid", *o.MsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("msid: nil\n"))
	}
	if o.NamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name", *o.NamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name: nil\n"))
	}
	if o.NameOrdinalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name-ordinal", *o.NameOrdinalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name-ordinal: nil\n"))
	}
	if o.NodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "node", *o.NodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("node: nil\n"))
	}
	if o.OwningVserverNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "owning-vserver-name", *o.OwningVserverNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("owning-vserver-name: nil\n"))
	}
	if o.OwningVserverUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "owning-vserver-uuid", *o.OwningVserverUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("owning-vserver-uuid: nil\n"))
	}
	if o.ProvenanceUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "provenance-uuid", *o.ProvenanceUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("provenance-uuid: nil\n"))
	}
	if o.StylePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "style", *o.StylePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("style: nil\n"))
	}
	if o.TypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "type", *o.TypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("type: nil\n"))
	}
	if o.UuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "uuid", *o.UuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("uuid: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeIdAttributesType) Comment() string {
	r := *o.CommentPtr
	return r
}

func (o *VolumeIdAttributesType) SetComment(newValue string) *VolumeIdAttributesType {
	o.CommentPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) ContainingAggregateName() string {
	r := *o.ContainingAggregateNamePtr
	return r
}

func (o *VolumeIdAttributesType) SetContainingAggregateName(newValue string) *VolumeIdAttributesType {
	o.ContainingAggregateNamePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) ContainingAggregateUuid() UuidType {
	r := *o.ContainingAggregateUuidPtr
	return r
}

func (o *VolumeIdAttributesType) SetContainingAggregateUuid(newValue UuidType) *VolumeIdAttributesType {
	o.ContainingAggregateUuidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) CreationTime() int {
	r := *o.CreationTimePtr
	return r
}

func (o *VolumeIdAttributesType) SetCreationTime(newValue int) *VolumeIdAttributesType {
	o.CreationTimePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Dsid() int {
	r := *o.DsidPtr
	return r
}

func (o *VolumeIdAttributesType) SetDsid(newValue int) *VolumeIdAttributesType {
	o.DsidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Fsid() string {
	r := *o.FsidPtr
	return r
}

func (o *VolumeIdAttributesType) SetFsid(newValue string) *VolumeIdAttributesType {
	o.FsidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) InstanceUuid() UuidType {
	r := *o.InstanceUuidPtr
	return r
}

func (o *VolumeIdAttributesType) SetInstanceUuid(newValue UuidType) *VolumeIdAttributesType {
	o.InstanceUuidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) JunctionParentName() VolumeNameType {
	r := *o.JunctionParentNamePtr
	return r
}

func (o *VolumeIdAttributesType) SetJunctionParentName(newValue VolumeNameType) *VolumeIdAttributesType {
	o.JunctionParentNamePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) JunctionPath() JunctionPathType {
	r := *o.JunctionPathPtr
	return r
}

func (o *VolumeIdAttributesType) SetJunctionPath(newValue JunctionPathType) *VolumeIdAttributesType {
	o.JunctionPathPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Msid() int {
	r := *o.MsidPtr
	return r
}

func (o *VolumeIdAttributesType) SetMsid(newValue int) *VolumeIdAttributesType {
	o.MsidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Name() VolumeNameType {
	r := *o.NamePtr
	return r
}

func (o *VolumeIdAttributesType) SetName(newValue VolumeNameType) *VolumeIdAttributesType {
	o.NamePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) NameOrdinal() string {
	r := *o.NameOrdinalPtr
	return r
}

func (o *VolumeIdAttributesType) SetNameOrdinal(newValue string) *VolumeIdAttributesType {
	o.NameOrdinalPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Node() NodeNameType {
	r := *o.NodePtr
	return r
}

func (o *VolumeIdAttributesType) SetNode(newValue NodeNameType) *VolumeIdAttributesType {
	o.NodePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) OwningVserverName() string {
	r := *o.OwningVserverNamePtr
	return r
}

func (o *VolumeIdAttributesType) SetOwningVserverName(newValue string) *VolumeIdAttributesType {
	o.OwningVserverNamePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) OwningVserverUuid() UuidType {
	r := *o.OwningVserverUuidPtr
	return r
}

func (o *VolumeIdAttributesType) SetOwningVserverUuid(newValue UuidType) *VolumeIdAttributesType {
	o.OwningVserverUuidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) ProvenanceUuid() UuidType {
	r := *o.ProvenanceUuidPtr
	return r
}

func (o *VolumeIdAttributesType) SetProvenanceUuid(newValue UuidType) *VolumeIdAttributesType {
	o.ProvenanceUuidPtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Style() string {
	r := *o.StylePtr
	return r
}

func (o *VolumeIdAttributesType) SetStyle(newValue string) *VolumeIdAttributesType {
	o.StylePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Type() string {
	r := *o.TypePtr
	return r
}

func (o *VolumeIdAttributesType) SetType(newValue string) *VolumeIdAttributesType {
	o.TypePtr = &newValue
	return o
}

func (o *VolumeIdAttributesType) Uuid() UuidType {
	r := *o.UuidPtr
	return r
}

func (o *VolumeIdAttributesType) SetUuid(newValue UuidType) *VolumeIdAttributesType {
	o.UuidPtr = &newValue
	return o
}

type VolumeHybridCacheAttributesType struct {
	XMLName xml.Name `xml:"volume-hybrid-cache-attributes"`

	CachingPolicyPtr                 *string `xml:"caching-policy"`
	EligibilityPtr                   *string `xml:"eligibility"`
	WriteCacheIneligibilityReasonPtr *string `xml:"write-cache-ineligibility-reason"`
}

func (o *VolumeHybridCacheAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeHybridCacheAttributesType() *VolumeHybridCacheAttributesType {
	return &VolumeHybridCacheAttributesType{}
}

func (o VolumeHybridCacheAttributesType) String() string {
	var buffer bytes.Buffer
	if o.CachingPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "caching-policy", *o.CachingPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("caching-policy: nil\n"))
	}
	if o.EligibilityPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "eligibility", *o.EligibilityPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("eligibility: nil\n"))
	}
	if o.WriteCacheIneligibilityReasonPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "write-cache-ineligibility-reason", *o.WriteCacheIneligibilityReasonPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("write-cache-ineligibility-reason: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeHybridCacheAttributesType) CachingPolicy() string {
	r := *o.CachingPolicyPtr
	return r
}

func (o *VolumeHybridCacheAttributesType) SetCachingPolicy(newValue string) *VolumeHybridCacheAttributesType {
	o.CachingPolicyPtr = &newValue
	return o
}

func (o *VolumeHybridCacheAttributesType) Eligibility() string {
	r := *o.EligibilityPtr
	return r
}

func (o *VolumeHybridCacheAttributesType) SetEligibility(newValue string) *VolumeHybridCacheAttributesType {
	o.EligibilityPtr = &newValue
	return o
}

func (o *VolumeHybridCacheAttributesType) WriteCacheIneligibilityReason() string {
	r := *o.WriteCacheIneligibilityReasonPtr
	return r
}

func (o *VolumeHybridCacheAttributesType) SetWriteCacheIneligibilityReason(newValue string) *VolumeHybridCacheAttributesType {
	o.WriteCacheIneligibilityReasonPtr = &newValue
	return o
}

type SizeType int

type CachePolicyType string

type VolumeFlexcacheAttributesType struct {
	XMLName xml.Name `xml:"volume-flexcache-attributes"`

	CachePolicyPtr *CachePolicyType `xml:"cache-policy"`
	FillPolicyPtr  *CachePolicyType `xml:"fill-policy"`
	MinReservePtr  *SizeType        `xml:"min-reserve"`
	OriginPtr      *VolumeNameType  `xml:"origin"`
}

func (o *VolumeFlexcacheAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeFlexcacheAttributesType() *VolumeFlexcacheAttributesType {
	return &VolumeFlexcacheAttributesType{}
}

func (o VolumeFlexcacheAttributesType) String() string {
	var buffer bytes.Buffer
	if o.CachePolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "cache-policy", *o.CachePolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("cache-policy: nil\n"))
	}
	if o.FillPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "fill-policy", *o.FillPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("fill-policy: nil\n"))
	}
	if o.MinReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "min-reserve", *o.MinReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("min-reserve: nil\n"))
	}
	if o.OriginPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "origin", *o.OriginPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("origin: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeFlexcacheAttributesType) CachePolicy() CachePolicyType {
	r := *o.CachePolicyPtr
	return r
}

func (o *VolumeFlexcacheAttributesType) SetCachePolicy(newValue CachePolicyType) *VolumeFlexcacheAttributesType {
	o.CachePolicyPtr = &newValue
	return o
}

func (o *VolumeFlexcacheAttributesType) FillPolicy() CachePolicyType {
	r := *o.FillPolicyPtr
	return r
}

func (o *VolumeFlexcacheAttributesType) SetFillPolicy(newValue CachePolicyType) *VolumeFlexcacheAttributesType {
	o.FillPolicyPtr = &newValue
	return o
}

func (o *VolumeFlexcacheAttributesType) MinReserve() SizeType {
	r := *o.MinReservePtr
	return r
}

func (o *VolumeFlexcacheAttributesType) SetMinReserve(newValue SizeType) *VolumeFlexcacheAttributesType {
	o.MinReservePtr = &newValue
	return o
}

func (o *VolumeFlexcacheAttributesType) Origin() VolumeNameType {
	r := *o.OriginPtr
	return r
}

func (o *VolumeFlexcacheAttributesType) SetOrigin(newValue VolumeNameType) *VolumeFlexcacheAttributesType {
	o.OriginPtr = &newValue
	return o
}

type VolumeExportAttributesType struct {
	XMLName xml.Name `xml:"volume-export-attributes"`

	PolicyPtr *string `xml:"policy"`
}

func (o *VolumeExportAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeExportAttributesType() *VolumeExportAttributesType { return &VolumeExportAttributesType{} }

func (o VolumeExportAttributesType) String() string {
	var buffer bytes.Buffer
	if o.PolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "policy", *o.PolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("policy: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeExportAttributesType) Policy() string {
	r := *o.PolicyPtr
	return r
}

func (o *VolumeExportAttributesType) SetPolicy(newValue string) *VolumeExportAttributesType {
	o.PolicyPtr = &newValue
	return o
}

type VolumeDirectoryAttributesType struct {
	XMLName xml.Name `xml:"volume-directory-attributes"`

	I2pEnabledPtr *bool   `xml:"i2p-enabled"`
	MaxDirSizePtr *int    `xml:"max-dir-size"`
	RootDirGenPtr *string `xml:"root-dir-gen"`
}

func (o *VolumeDirectoryAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeDirectoryAttributesType() *VolumeDirectoryAttributesType {
	return &VolumeDirectoryAttributesType{}
}

func (o VolumeDirectoryAttributesType) String() string {
	var buffer bytes.Buffer
	if o.I2pEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "i2p-enabled", *o.I2pEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("i2p-enabled: nil\n"))
	}
	if o.MaxDirSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "max-dir-size", *o.MaxDirSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("max-dir-size: nil\n"))
	}
	if o.RootDirGenPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "root-dir-gen", *o.RootDirGenPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("root-dir-gen: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeDirectoryAttributesType) I2pEnabled() bool {
	r := *o.I2pEnabledPtr
	return r
}

func (o *VolumeDirectoryAttributesType) SetI2pEnabled(newValue bool) *VolumeDirectoryAttributesType {
	o.I2pEnabledPtr = &newValue
	return o
}

func (o *VolumeDirectoryAttributesType) MaxDirSize() int {
	r := *o.MaxDirSizePtr
	return r
}

func (o *VolumeDirectoryAttributesType) SetMaxDirSize(newValue int) *VolumeDirectoryAttributesType {
	o.MaxDirSizePtr = &newValue
	return o
}

func (o *VolumeDirectoryAttributesType) RootDirGen() string {
	r := *o.RootDirGenPtr
	return r
}

func (o *VolumeDirectoryAttributesType) SetRootDirGen(newValue string) *VolumeDirectoryAttributesType {
	o.RootDirGenPtr = &newValue
	return o
}

type VolumeNameType string

type VolumeCloneParentAttributesType struct {
	XMLName xml.Name `xml:"volume-clone-parent-attributes"`

	DsidPtr         *int            `xml:"dsid"`
	MsidPtr         *int            `xml:"msid"`
	NamePtr         *VolumeNameType `xml:"name"`
	SnapshotIdPtr   *int            `xml:"snapshot-id"`
	SnapshotNamePtr *string         `xml:"snapshot-name"`
	UuidPtr         *UuidType       `xml:"uuid"`
}

func (o *VolumeCloneParentAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeCloneParentAttributesType() *VolumeCloneParentAttributesType {
	return &VolumeCloneParentAttributesType{}
}

func (o VolumeCloneParentAttributesType) String() string {
	var buffer bytes.Buffer
	if o.DsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "dsid", *o.DsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("dsid: nil\n"))
	}
	if o.MsidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "msid", *o.MsidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("msid: nil\n"))
	}
	if o.NamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name", *o.NamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name: nil\n"))
	}
	if o.SnapshotIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-id", *o.SnapshotIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-id: nil\n"))
	}
	if o.SnapshotNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-name", *o.SnapshotNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-name: nil\n"))
	}
	if o.UuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "uuid", *o.UuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("uuid: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeCloneParentAttributesType) Dsid() int {
	r := *o.DsidPtr
	return r
}

func (o *VolumeCloneParentAttributesType) SetDsid(newValue int) *VolumeCloneParentAttributesType {
	o.DsidPtr = &newValue
	return o
}

func (o *VolumeCloneParentAttributesType) Msid() int {
	r := *o.MsidPtr
	return r
}

func (o *VolumeCloneParentAttributesType) SetMsid(newValue int) *VolumeCloneParentAttributesType {
	o.MsidPtr = &newValue
	return o
}

func (o *VolumeCloneParentAttributesType) Name() VolumeNameType {
	r := *o.NamePtr
	return r
}

func (o *VolumeCloneParentAttributesType) SetName(newValue VolumeNameType) *VolumeCloneParentAttributesType {
	o.NamePtr = &newValue
	return o
}

func (o *VolumeCloneParentAttributesType) SnapshotId() int {
	r := *o.SnapshotIdPtr
	return r
}

func (o *VolumeCloneParentAttributesType) SetSnapshotId(newValue int) *VolumeCloneParentAttributesType {
	o.SnapshotIdPtr = &newValue
	return o
}

func (o *VolumeCloneParentAttributesType) SnapshotName() string {
	r := *o.SnapshotNamePtr
	return r
}

func (o *VolumeCloneParentAttributesType) SetSnapshotName(newValue string) *VolumeCloneParentAttributesType {
	o.SnapshotNamePtr = &newValue
	return o
}

func (o *VolumeCloneParentAttributesType) Uuid() UuidType {
	r := *o.UuidPtr
	return r
}

func (o *VolumeCloneParentAttributesType) SetUuid(newValue UuidType) *VolumeCloneParentAttributesType {
	o.UuidPtr = &newValue
	return o
}

type VolumeCloneAttributesType struct {
	XMLName xml.Name `xml:"volume-clone-attributes"`

	CloneChildCountPtr             *int                             `xml:"clone-child-count"`
	VolumeCloneParentAttributesPtr *VolumeCloneParentAttributesType `xml:"volume-clone-parent-attributes"`
}

func (o *VolumeCloneAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeCloneAttributesType() *VolumeCloneAttributesType { return &VolumeCloneAttributesType{} }

func (o VolumeCloneAttributesType) String() string {
	var buffer bytes.Buffer
	if o.CloneChildCountPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "clone-child-count", *o.CloneChildCountPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("clone-child-count: nil\n"))
	}
	if o.VolumeCloneParentAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-clone-parent-attributes", *o.VolumeCloneParentAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-clone-parent-attributes: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeCloneAttributesType) CloneChildCount() int {
	r := *o.CloneChildCountPtr
	return r
}

func (o *VolumeCloneAttributesType) SetCloneChildCount(newValue int) *VolumeCloneAttributesType {
	o.CloneChildCountPtr = &newValue
	return o
}

func (o *VolumeCloneAttributesType) VolumeCloneParentAttributes() VolumeCloneParentAttributesType {
	r := *o.VolumeCloneParentAttributesPtr
	return r
}

func (o *VolumeCloneAttributesType) SetVolumeCloneParentAttributes(newValue VolumeCloneParentAttributesType) *VolumeCloneAttributesType {
	o.VolumeCloneParentAttributesPtr = &newValue
	return o
}

type VolumeAutosizeAttributesType struct {
	XMLName xml.Name `xml:"volume-autosize-attributes"`

	GrowThresholdPercentPtr   *int    `xml:"grow-threshold-percent"`
	IncrementPercentPtr       *int    `xml:"increment-percent"`
	IncrementSizePtr          *int    `xml:"increment-size"`
	IsEnabledPtr              *bool   `xml:"is-enabled"`
	MaximumSizePtr            *int    `xml:"maximum-size"`
	MinimumSizePtr            *int    `xml:"minimum-size"`
	ModePtr                   *string `xml:"mode"`
	ResetPtr                  *bool   `xml:"reset"`
	ShrinkThresholdPercentPtr *int    `xml:"shrink-threshold-percent"`
}

func (o *VolumeAutosizeAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeAutosizeAttributesType() *VolumeAutosizeAttributesType {
	return &VolumeAutosizeAttributesType{}
}

func (o VolumeAutosizeAttributesType) String() string {
	var buffer bytes.Buffer
	if o.GrowThresholdPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "grow-threshold-percent", *o.GrowThresholdPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("grow-threshold-percent: nil\n"))
	}
	if o.IncrementPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "increment-percent", *o.IncrementPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("increment-percent: nil\n"))
	}
	if o.IncrementSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "increment-size", *o.IncrementSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("increment-size: nil\n"))
	}
	if o.IsEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-enabled", *o.IsEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-enabled: nil\n"))
	}
	if o.MaximumSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "maximum-size", *o.MaximumSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("maximum-size: nil\n"))
	}
	if o.MinimumSizePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "minimum-size", *o.MinimumSizePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("minimum-size: nil\n"))
	}
	if o.ModePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "mode", *o.ModePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("mode: nil\n"))
	}
	if o.ResetPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "reset", *o.ResetPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("reset: nil\n"))
	}
	if o.ShrinkThresholdPercentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "shrink-threshold-percent", *o.ShrinkThresholdPercentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("shrink-threshold-percent: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeAutosizeAttributesType) GrowThresholdPercent() int {
	r := *o.GrowThresholdPercentPtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetGrowThresholdPercent(newValue int) *VolumeAutosizeAttributesType {
	o.GrowThresholdPercentPtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) IncrementPercent() int {
	r := *o.IncrementPercentPtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetIncrementPercent(newValue int) *VolumeAutosizeAttributesType {
	o.IncrementPercentPtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) IncrementSize() int {
	r := *o.IncrementSizePtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetIncrementSize(newValue int) *VolumeAutosizeAttributesType {
	o.IncrementSizePtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) IsEnabled() bool {
	r := *o.IsEnabledPtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetIsEnabled(newValue bool) *VolumeAutosizeAttributesType {
	o.IsEnabledPtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) MaximumSize() int {
	r := *o.MaximumSizePtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetMaximumSize(newValue int) *VolumeAutosizeAttributesType {
	o.MaximumSizePtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) MinimumSize() int {
	r := *o.MinimumSizePtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetMinimumSize(newValue int) *VolumeAutosizeAttributesType {
	o.MinimumSizePtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) Mode() string {
	r := *o.ModePtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetMode(newValue string) *VolumeAutosizeAttributesType {
	o.ModePtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) Reset() bool {
	r := *o.ResetPtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetReset(newValue bool) *VolumeAutosizeAttributesType {
	o.ResetPtr = &newValue
	return o
}

func (o *VolumeAutosizeAttributesType) ShrinkThresholdPercent() int {
	r := *o.ShrinkThresholdPercentPtr
	return r
}

func (o *VolumeAutosizeAttributesType) SetShrinkThresholdPercent(newValue int) *VolumeAutosizeAttributesType {
	o.ShrinkThresholdPercentPtr = &newValue
	return o
}

type VolumeAutobalanceAttributesType struct {
	XMLName xml.Name `xml:"volume-autobalance-attributes"`

	IsAutobalanceEligiblePtr *bool `xml:"is-autobalance-eligible"`
}

func (o *VolumeAutobalanceAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeAutobalanceAttributesType() *VolumeAutobalanceAttributesType {
	return &VolumeAutobalanceAttributesType{}
}

func (o VolumeAutobalanceAttributesType) String() string {
	var buffer bytes.Buffer
	if o.IsAutobalanceEligiblePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-autobalance-eligible", *o.IsAutobalanceEligiblePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-autobalance-eligible: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeAutobalanceAttributesType) IsAutobalanceEligible() bool {
	r := *o.IsAutobalanceEligiblePtr
	return r
}

func (o *VolumeAutobalanceAttributesType) SetIsAutobalanceEligible(newValue bool) *VolumeAutobalanceAttributesType {
	o.IsAutobalanceEligiblePtr = &newValue
	return o
}

type VolumeAntivirusAttributesType struct {
	XMLName xml.Name `xml:"volume-antivirus-attributes"`

	OnAccessPolicyPtr *string `xml:"on-access-policy"`
}

func (o *VolumeAntivirusAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeAntivirusAttributesType() *VolumeAntivirusAttributesType {
	return &VolumeAntivirusAttributesType{}
}

func (o VolumeAntivirusAttributesType) String() string {
	var buffer bytes.Buffer
	if o.OnAccessPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "on-access-policy", *o.OnAccessPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("on-access-policy: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeAntivirusAttributesType) OnAccessPolicy() string {
	r := *o.OnAccessPolicyPtr
	return r
}

func (o *VolumeAntivirusAttributesType) SetOnAccessPolicy(newValue string) *VolumeAntivirusAttributesType {
	o.OnAccessPolicyPtr = &newValue
	return o
}

type VolumeAttributesType struct {
	XMLName xml.Name `xml:"volume-attributes"`

	VolumeAntivirusAttributesPtr          *VolumeAntivirusAttributesType          `xml:"volume-antivirus-attributes"`
	VolumeAutobalanceAttributesPtr        *VolumeAutobalanceAttributesType        `xml:"volume-autobalance-attributes"`
	VolumeAutosizeAttributesPtr           *VolumeAutosizeAttributesType           `xml:"volume-autosize-attributes"`
	VolumeCloneAttributesPtr              *VolumeCloneAttributesType              `xml:"volume-clone-attributes"`
	VolumeDirectoryAttributesPtr          *VolumeDirectoryAttributesType          `xml:"volume-directory-attributes"`
	VolumeExportAttributesPtr             *VolumeExportAttributesType             `xml:"volume-export-attributes"`
	VolumeFlexcacheAttributesPtr          *VolumeFlexcacheAttributesType          `xml:"volume-flexcache-attributes"`
	VolumeHybridCacheAttributesPtr        *VolumeHybridCacheAttributesType        `xml:"volume-hybrid-cache-attributes"`
	VolumeIdAttributesPtr                 *VolumeIdAttributesType                 `xml:"volume-id-attributes"`
	VolumeInfinitevolAttributesPtr        *VolumeInfinitevolAttributesType        `xml:"volume-infinitevol-attributes"`
	VolumeInodeAttributesPtr              *VolumeInodeAttributesType              `xml:"volume-inode-attributes"`
	VolumeLanguageAttributesPtr           *VolumeLanguageAttributesType           `xml:"volume-language-attributes"`
	VolumeMirrorAttributesPtr             *VolumeMirrorAttributesType             `xml:"volume-mirror-attributes"`
	VolumePerformanceAttributesPtr        *VolumePerformanceAttributesType        `xml:"volume-performance-attributes"`
	VolumeQosAttributesPtr                *VolumeQosAttributesType                `xml:"volume-qos-attributes"`
	VolumeSecurityAttributesPtr           *VolumeSecurityAttributesType           `xml:"volume-security-attributes"`
	VolumeSisAttributesPtr                *VolumeSisAttributesType                `xml:"volume-sis-attributes"`
	VolumeSnapshotAttributesPtr           *VolumeSnapshotAttributesType           `xml:"volume-snapshot-attributes"`
	VolumeSnapshotAutodeleteAttributesPtr *VolumeSnapshotAutodeleteAttributesType `xml:"volume-snapshot-autodelete-attributes"`
	VolumeSpaceAttributesPtr              *VolumeSpaceAttributesType              `xml:"volume-space-attributes"`
	VolumeStateAttributesPtr              *VolumeStateAttributesType              `xml:"volume-state-attributes"`
	VolumeTransitionAttributesPtr         *VolumeTransitionAttributesType         `xml:"volume-transition-attributes"`
	VolumeVmAlignAttributesPtr            *VolumeVmAlignAttributesType            `xml:"volume-vm-align-attributes"`
}

func (o *VolumeAttributesType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewVolumeAttributesType() *VolumeAttributesType { return &VolumeAttributesType{} }

func (o VolumeAttributesType) String() string {
	var buffer bytes.Buffer
	if o.VolumeAntivirusAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-antivirus-attributes", *o.VolumeAntivirusAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-antivirus-attributes: nil\n"))
	}
	if o.VolumeAutobalanceAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-autobalance-attributes", *o.VolumeAutobalanceAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-autobalance-attributes: nil\n"))
	}
	if o.VolumeAutosizeAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-autosize-attributes", *o.VolumeAutosizeAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-autosize-attributes: nil\n"))
	}
	if o.VolumeCloneAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-clone-attributes", *o.VolumeCloneAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-clone-attributes: nil\n"))
	}
	if o.VolumeDirectoryAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-directory-attributes", *o.VolumeDirectoryAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-directory-attributes: nil\n"))
	}
	if o.VolumeExportAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-export-attributes", *o.VolumeExportAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-export-attributes: nil\n"))
	}
	if o.VolumeFlexcacheAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-flexcache-attributes", *o.VolumeFlexcacheAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-flexcache-attributes: nil\n"))
	}
	if o.VolumeHybridCacheAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-hybrid-cache-attributes", *o.VolumeHybridCacheAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-hybrid-cache-attributes: nil\n"))
	}
	if o.VolumeIdAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-id-attributes", *o.VolumeIdAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-id-attributes: nil\n"))
	}
	if o.VolumeInfinitevolAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-infinitevol-attributes", *o.VolumeInfinitevolAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-infinitevol-attributes: nil\n"))
	}
	if o.VolumeInodeAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-inode-attributes", *o.VolumeInodeAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-inode-attributes: nil\n"))
	}
	if o.VolumeLanguageAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-language-attributes", *o.VolumeLanguageAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-language-attributes: nil\n"))
	}
	if o.VolumeMirrorAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-mirror-attributes", *o.VolumeMirrorAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-mirror-attributes: nil\n"))
	}
	if o.VolumePerformanceAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-performance-attributes", *o.VolumePerformanceAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-performance-attributes: nil\n"))
	}
	if o.VolumeQosAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-qos-attributes", *o.VolumeQosAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-qos-attributes: nil\n"))
	}
	if o.VolumeSecurityAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-security-attributes", *o.VolumeSecurityAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-security-attributes: nil\n"))
	}
	if o.VolumeSisAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-sis-attributes", *o.VolumeSisAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-sis-attributes: nil\n"))
	}
	if o.VolumeSnapshotAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-snapshot-attributes", *o.VolumeSnapshotAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-snapshot-attributes: nil\n"))
	}
	if o.VolumeSnapshotAutodeleteAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-snapshot-autodelete-attributes", *o.VolumeSnapshotAutodeleteAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-snapshot-autodelete-attributes: nil\n"))
	}
	if o.VolumeSpaceAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-space-attributes", *o.VolumeSpaceAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-space-attributes: nil\n"))
	}
	if o.VolumeStateAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-state-attributes", *o.VolumeStateAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-state-attributes: nil\n"))
	}
	if o.VolumeTransitionAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-transition-attributes", *o.VolumeTransitionAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-transition-attributes: nil\n"))
	}
	if o.VolumeVmAlignAttributesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-vm-align-attributes", *o.VolumeVmAlignAttributesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-vm-align-attributes: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeAttributesType) VolumeAntivirusAttributes() VolumeAntivirusAttributesType {
	r := *o.VolumeAntivirusAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeAntivirusAttributes(newValue VolumeAntivirusAttributesType) *VolumeAttributesType {
	o.VolumeAntivirusAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeAutobalanceAttributes() VolumeAutobalanceAttributesType {
	r := *o.VolumeAutobalanceAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeAutobalanceAttributes(newValue VolumeAutobalanceAttributesType) *VolumeAttributesType {
	o.VolumeAutobalanceAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeAutosizeAttributes() VolumeAutosizeAttributesType {
	r := *o.VolumeAutosizeAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeAutosizeAttributes(newValue VolumeAutosizeAttributesType) *VolumeAttributesType {
	o.VolumeAutosizeAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeCloneAttributes() VolumeCloneAttributesType {
	r := *o.VolumeCloneAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeCloneAttributes(newValue VolumeCloneAttributesType) *VolumeAttributesType {
	o.VolumeCloneAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeDirectoryAttributes() VolumeDirectoryAttributesType {
	r := *o.VolumeDirectoryAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeDirectoryAttributes(newValue VolumeDirectoryAttributesType) *VolumeAttributesType {
	o.VolumeDirectoryAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeExportAttributes() VolumeExportAttributesType {
	r := *o.VolumeExportAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeExportAttributes(newValue VolumeExportAttributesType) *VolumeAttributesType {
	o.VolumeExportAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeFlexcacheAttributes() VolumeFlexcacheAttributesType {
	r := *o.VolumeFlexcacheAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeFlexcacheAttributes(newValue VolumeFlexcacheAttributesType) *VolumeAttributesType {
	o.VolumeFlexcacheAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeHybridCacheAttributes() VolumeHybridCacheAttributesType {
	r := *o.VolumeHybridCacheAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeHybridCacheAttributes(newValue VolumeHybridCacheAttributesType) *VolumeAttributesType {
	o.VolumeHybridCacheAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeIdAttributes() VolumeIdAttributesType {
	r := *o.VolumeIdAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeIdAttributes(newValue VolumeIdAttributesType) *VolumeAttributesType {
	o.VolumeIdAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeInfinitevolAttributes() VolumeInfinitevolAttributesType {
	r := *o.VolumeInfinitevolAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeInfinitevolAttributes(newValue VolumeInfinitevolAttributesType) *VolumeAttributesType {
	o.VolumeInfinitevolAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeInodeAttributes() VolumeInodeAttributesType {
	r := *o.VolumeInodeAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeInodeAttributes(newValue VolumeInodeAttributesType) *VolumeAttributesType {
	o.VolumeInodeAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeLanguageAttributes() VolumeLanguageAttributesType {
	r := *o.VolumeLanguageAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeLanguageAttributes(newValue VolumeLanguageAttributesType) *VolumeAttributesType {
	o.VolumeLanguageAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeMirrorAttributes() VolumeMirrorAttributesType {
	r := *o.VolumeMirrorAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeMirrorAttributes(newValue VolumeMirrorAttributesType) *VolumeAttributesType {
	o.VolumeMirrorAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumePerformanceAttributes() VolumePerformanceAttributesType {
	r := *o.VolumePerformanceAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumePerformanceAttributes(newValue VolumePerformanceAttributesType) *VolumeAttributesType {
	o.VolumePerformanceAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeQosAttributes() VolumeQosAttributesType {
	r := *o.VolumeQosAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeQosAttributes(newValue VolumeQosAttributesType) *VolumeAttributesType {
	o.VolumeQosAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeSecurityAttributes() VolumeSecurityAttributesType {
	r := *o.VolumeSecurityAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeSecurityAttributes(newValue VolumeSecurityAttributesType) *VolumeAttributesType {
	o.VolumeSecurityAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeSisAttributes() VolumeSisAttributesType {
	r := *o.VolumeSisAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeSisAttributes(newValue VolumeSisAttributesType) *VolumeAttributesType {
	o.VolumeSisAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeSnapshotAttributes() VolumeSnapshotAttributesType {
	r := *o.VolumeSnapshotAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeSnapshotAttributes(newValue VolumeSnapshotAttributesType) *VolumeAttributesType {
	o.VolumeSnapshotAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeSnapshotAutodeleteAttributes() VolumeSnapshotAutodeleteAttributesType {
	r := *o.VolumeSnapshotAutodeleteAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeSnapshotAutodeleteAttributes(newValue VolumeSnapshotAutodeleteAttributesType) *VolumeAttributesType {
	o.VolumeSnapshotAutodeleteAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeSpaceAttributes() VolumeSpaceAttributesType {
	r := *o.VolumeSpaceAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeSpaceAttributes(newValue VolumeSpaceAttributesType) *VolumeAttributesType {
	o.VolumeSpaceAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeStateAttributes() VolumeStateAttributesType {
	r := *o.VolumeStateAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeStateAttributes(newValue VolumeStateAttributesType) *VolumeAttributesType {
	o.VolumeStateAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeTransitionAttributes() VolumeTransitionAttributesType {
	r := *o.VolumeTransitionAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeTransitionAttributes(newValue VolumeTransitionAttributesType) *VolumeAttributesType {
	o.VolumeTransitionAttributesPtr = &newValue
	return o
}

func (o *VolumeAttributesType) VolumeVmAlignAttributes() VolumeVmAlignAttributesType {
	r := *o.VolumeVmAlignAttributesPtr
	return r
}

func (o *VolumeAttributesType) SetVolumeVmAlignAttributes(newValue VolumeVmAlignAttributesType) *VolumeAttributesType {
	o.VolumeVmAlignAttributesPtr = &newValue
	return o
}

type SystemVersionTupleType struct {
	XMLName xml.Name `xml:"system-version-tuple"`

	GenerationPtr *int `xml:"generation"`
	MajorPtr      *int `xml:"major"`
	MinorPtr      *int `xml:"minor"`
}

func (o *SystemVersionTupleType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewSystemVersionTupleType() *SystemVersionTupleType { return &SystemVersionTupleType{} }

func (o SystemVersionTupleType) String() string {
	var buffer bytes.Buffer
	if o.GenerationPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "generation", *o.GenerationPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("generation: nil\n"))
	}
	if o.MajorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "major", *o.MajorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("major: nil\n"))
	}
	if o.MinorPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "minor", *o.MinorPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("minor: nil\n"))
	}
	return buffer.String()
}

func (o *SystemVersionTupleType) Generation() int {
	r := *o.GenerationPtr
	return r
}

func (o *SystemVersionTupleType) SetGeneration(newValue int) *SystemVersionTupleType {
	o.GenerationPtr = &newValue
	return o
}

func (o *SystemVersionTupleType) Major() int {
	r := *o.MajorPtr
	return r
}

func (o *SystemVersionTupleType) SetMajor(newValue int) *SystemVersionTupleType {
	o.MajorPtr = &newValue
	return o
}

func (o *SystemVersionTupleType) Minor() int {
	r := *o.MinorPtr
	return r
}

func (o *SystemVersionTupleType) SetMinor(newValue int) *SystemVersionTupleType {
	o.MinorPtr = &newValue
	return o
}

type SubnetNameType string

type RoutingGroupType string

type UuidType string

type FailoverGroupType string

type DnsZoneType string

type DataProtocolType string

type IpAddressType string

type InitiatorInfoType struct {
	XMLName xml.Name `xml:"initiator-info"`

	InitiatorNamePtr *string `xml:"initiator-name"`
}

func (o *InitiatorInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewInitiatorInfoType() *InitiatorInfoType { return &InitiatorInfoType{} }

func (o InitiatorInfoType) String() string {
	var buffer bytes.Buffer
	if o.InitiatorNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-name", *o.InitiatorNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-name: nil\n"))
	}
	return buffer.String()
}

func (o *InitiatorInfoType) InitiatorName() string {
	r := *o.InitiatorNamePtr
	return r
}

func (o *InitiatorInfoType) SetInitiatorName(newValue string) *InitiatorInfoType {
	o.InitiatorNamePtr = &newValue
	return o
}

type InitiatorGroupInfoType struct {
	XMLName xml.Name `xml:"initiator-group-info"`

	InitiatorGroupAluaEnabledPtr           *bool               `xml:"initiator-group-alua-enabled"`
	InitiatorGroupNamePtr                  *string             `xml:"initiator-group-name"`
	InitiatorGroupOsTypePtr                *string             `xml:"initiator-group-os-type"`
	InitiatorGroupPortsetNamePtr           *string             `xml:"initiator-group-portset-name"`
	InitiatorGroupReportScsiNameEnabledPtr *bool               `xml:"initiator-group-report-scsi-name-enabled"`
	InitiatorGroupThrottleBorrowPtr        *bool               `xml:"initiator-group-throttle-borrow"`
	InitiatorGroupThrottleReservePtr       *int                `xml:"initiator-group-throttle-reserve"`
	InitiatorGroupTypePtr                  *string             `xml:"initiator-group-type"`
	InitiatorGroupUsePartnerPtr            *bool               `xml:"initiator-group-use-partner"`
	InitiatorGroupUuidPtr                  *string             `xml:"initiator-group-uuid"`
	InitiatorGroupVsaEnabledPtr            *bool               `xml:"initiator-group-vsa-enabled"`
	InitiatorsPtr                          []InitiatorInfoType `xml:"initiators>initiator-info"`
	LunIdPtr                               *int                `xml:"lun-id"`
	VserverPtr                             *string             `xml:"vserver"`
}

func (o *InitiatorGroupInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewInitiatorGroupInfoType() *InitiatorGroupInfoType { return &InitiatorGroupInfoType{} }

func (o InitiatorGroupInfoType) String() string {
	var buffer bytes.Buffer
	if o.InitiatorGroupAluaEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-alua-enabled", *o.InitiatorGroupAluaEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-alua-enabled: nil\n"))
	}
	if o.InitiatorGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-name", *o.InitiatorGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-name: nil\n"))
	}
	if o.InitiatorGroupOsTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-os-type", *o.InitiatorGroupOsTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-os-type: nil\n"))
	}
	if o.InitiatorGroupPortsetNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-portset-name", *o.InitiatorGroupPortsetNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-portset-name: nil\n"))
	}
	if o.InitiatorGroupReportScsiNameEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-report-scsi-name-enabled", *o.InitiatorGroupReportScsiNameEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-report-scsi-name-enabled: nil\n"))
	}
	if o.InitiatorGroupThrottleBorrowPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-throttle-borrow", *o.InitiatorGroupThrottleBorrowPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-throttle-borrow: nil\n"))
	}
	if o.InitiatorGroupThrottleReservePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-throttle-reserve", *o.InitiatorGroupThrottleReservePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-throttle-reserve: nil\n"))
	}
	if o.InitiatorGroupTypePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-type", *o.InitiatorGroupTypePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-type: nil\n"))
	}
	if o.InitiatorGroupUsePartnerPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-use-partner", *o.InitiatorGroupUsePartnerPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-use-partner: nil\n"))
	}
	if o.InitiatorGroupUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-uuid", *o.InitiatorGroupUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-uuid: nil\n"))
	}
	if o.InitiatorGroupVsaEnabledPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiator-group-vsa-enabled", *o.InitiatorGroupVsaEnabledPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiator-group-vsa-enabled: nil\n"))
	}
	if o.InitiatorsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "initiators", o.InitiatorsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("initiators: nil\n"))
	}
	if o.LunIdPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "lun-id", *o.LunIdPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("lun-id: nil\n"))
	}
	if o.VserverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver", *o.VserverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver: nil\n"))
	}
	return buffer.String()
}

func (o *InitiatorGroupInfoType) InitiatorGroupAluaEnabled() bool {
	r := *o.InitiatorGroupAluaEnabledPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupAluaEnabled(newValue bool) *InitiatorGroupInfoType {
	o.InitiatorGroupAluaEnabledPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupName() string {
	r := *o.InitiatorGroupNamePtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupName(newValue string) *InitiatorGroupInfoType {
	o.InitiatorGroupNamePtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupOsType() string {
	r := *o.InitiatorGroupOsTypePtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupOsType(newValue string) *InitiatorGroupInfoType {
	o.InitiatorGroupOsTypePtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupPortsetName() string {
	r := *o.InitiatorGroupPortsetNamePtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupPortsetName(newValue string) *InitiatorGroupInfoType {
	o.InitiatorGroupPortsetNamePtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupReportScsiNameEnabled() bool {
	r := *o.InitiatorGroupReportScsiNameEnabledPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupReportScsiNameEnabled(newValue bool) *InitiatorGroupInfoType {
	o.InitiatorGroupReportScsiNameEnabledPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupThrottleBorrow() bool {
	r := *o.InitiatorGroupThrottleBorrowPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupThrottleBorrow(newValue bool) *InitiatorGroupInfoType {
	o.InitiatorGroupThrottleBorrowPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupThrottleReserve() int {
	r := *o.InitiatorGroupThrottleReservePtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupThrottleReserve(newValue int) *InitiatorGroupInfoType {
	o.InitiatorGroupThrottleReservePtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupType() string {
	r := *o.InitiatorGroupTypePtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupType(newValue string) *InitiatorGroupInfoType {
	o.InitiatorGroupTypePtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupUsePartner() bool {
	r := *o.InitiatorGroupUsePartnerPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupUsePartner(newValue bool) *InitiatorGroupInfoType {
	o.InitiatorGroupUsePartnerPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupUuid() string {
	r := *o.InitiatorGroupUuidPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupUuid(newValue string) *InitiatorGroupInfoType {
	o.InitiatorGroupUuidPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) InitiatorGroupVsaEnabled() bool {
	r := *o.InitiatorGroupVsaEnabledPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiatorGroupVsaEnabled(newValue bool) *InitiatorGroupInfoType {
	o.InitiatorGroupVsaEnabledPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) Initiators() []InitiatorInfoType {
	r := o.InitiatorsPtr
	return r
}

func (o *InitiatorGroupInfoType) SetInitiators(newValue []InitiatorInfoType) *InitiatorGroupInfoType {
	newSlice := make([]InitiatorInfoType, len(newValue))
	copy(newSlice, newValue)
	o.InitiatorsPtr = newSlice
	return o
}

func (o *InitiatorGroupInfoType) LunId() int {
	r := *o.LunIdPtr
	return r
}

func (o *InitiatorGroupInfoType) SetLunId(newValue int) *InitiatorGroupInfoType {
	o.LunIdPtr = &newValue
	return o
}

func (o *InitiatorGroupInfoType) Vserver() string {
	r := *o.VserverPtr
	return r
}

func (o *InitiatorGroupInfoType) SetVserver(newValue string) *InitiatorGroupInfoType {
	o.VserverPtr = &newValue
	return o
}

type NodeNameType string

type NetInterfaceInfoType struct {
	XMLName xml.Name `xml:"net-interface-info"`

	AddressPtr                *IpAddressType     `xml:"address"`
	AddressFamilyPtr          *string            `xml:"address-family"`
	AdministrativeStatusPtr   *string            `xml:"administrative-status"`
	CommentPtr                *string            `xml:"comment"`
	CurrentNodePtr            *string            `xml:"current-node"`
	CurrentPortPtr            *string            `xml:"current-port"`
	DataProtocolsPtr          []DataProtocolType `xml:"data-protocols>data-protocol"`
	DnsDomainNamePtr          *DnsZoneType       `xml:"dns-domain-name"`
	FailoverGroupPtr          *FailoverGroupType `xml:"failover-group"`
	FailoverPolicyPtr         *string            `xml:"failover-policy"`
	FirewallPolicyPtr         *string            `xml:"firewall-policy"`
	ForceSubnetAssociationPtr *bool              `xml:"force-subnet-association"`
	HomeNodePtr               *string            `xml:"home-node"`
	HomePortPtr               *string            `xml:"home-port"`
	InterfaceNamePtr          *string            `xml:"interface-name"`
	IsAutoRevertPtr           *bool              `xml:"is-auto-revert"`
	IsHomePtr                 *bool              `xml:"is-home"`
	IsIpv4LinkLocalPtr        *bool              `xml:"is-ipv4-link-local"`
	LifUuidPtr                *UuidType          `xml:"lif-uuid"`
	ListenForDnsQueryPtr      *bool              `xml:"listen-for-dns-query"`
	NetmaskPtr                *IpAddressType     `xml:"netmask"`
	NetmaskLengthPtr          *int               `xml:"netmask-length"`
	OperationalStatusPtr      *string            `xml:"operational-status"`
	RolePtr                   *string            `xml:"role"`
	RoutingGroupNamePtr       *RoutingGroupType  `xml:"routing-group-name"`
	SubnetNamePtr             *SubnetNameType    `xml:"subnet-name"`
	UseFailoverGroupPtr       *string            `xml:"use-failover-group"`
	VserverPtr                *string            `xml:"vserver"`
	WwpnPtr                   *string            `xml:"wwpn"`
}

func (o *NetInterfaceInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		log.Errorf("error: %v", err)
	}
	return string(output), err
}

func NewNetInterfaceInfoType() *NetInterfaceInfoType { return &NetInterfaceInfoType{} }

func (o NetInterfaceInfoType) String() string {
	var buffer bytes.Buffer
	if o.AddressPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "address", *o.AddressPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("address: nil\n"))
	}
	if o.AddressFamilyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "address-family", *o.AddressFamilyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("address-family: nil\n"))
	}
	if o.AdministrativeStatusPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "administrative-status", *o.AdministrativeStatusPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("administrative-status: nil\n"))
	}
	if o.CommentPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "comment", *o.CommentPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("comment: nil\n"))
	}
	if o.CurrentNodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "current-node", *o.CurrentNodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("current-node: nil\n"))
	}
	if o.CurrentPortPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "current-port", *o.CurrentPortPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("current-port: nil\n"))
	}
	if o.DataProtocolsPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "data-protocols", o.DataProtocolsPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("data-protocols: nil\n"))
	}
	if o.DnsDomainNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "dns-domain-name", *o.DnsDomainNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("dns-domain-name: nil\n"))
	}
	if o.FailoverGroupPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "failover-group", *o.FailoverGroupPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("failover-group: nil\n"))
	}
	if o.FailoverPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "failover-policy", *o.FailoverPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("failover-policy: nil\n"))
	}
	if o.FirewallPolicyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "firewall-policy", *o.FirewallPolicyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("firewall-policy: nil\n"))
	}
	if o.ForceSubnetAssociationPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "force-subnet-association", *o.ForceSubnetAssociationPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("force-subnet-association: nil\n"))
	}
	if o.HomeNodePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "home-node", *o.HomeNodePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("home-node: nil\n"))
	}
	if o.HomePortPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "home-port", *o.HomePortPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("home-port: nil\n"))
	}
	if o.InterfaceNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "interface-name", *o.InterfaceNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("interface-name: nil\n"))
	}
	if o.IsAutoRevertPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-auto-revert", *o.IsAutoRevertPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-auto-revert: nil\n"))
	}
	if o.IsHomePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-home", *o.IsHomePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-home: nil\n"))
	}
	if o.IsIpv4LinkLocalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-ipv4-link-local", *o.IsIpv4LinkLocalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-ipv4-link-local: nil\n"))
	}
	if o.LifUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "lif-uuid", *o.LifUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("lif-uuid: nil\n"))
	}
	if o.ListenForDnsQueryPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "listen-for-dns-query", *o.ListenForDnsQueryPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("listen-for-dns-query: nil\n"))
	}
	if o.NetmaskPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "netmask", *o.NetmaskPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("netmask: nil\n"))
	}
	if o.NetmaskLengthPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "netmask-length", *o.NetmaskLengthPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("netmask-length: nil\n"))
	}
	if o.OperationalStatusPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "operational-status", *o.OperationalStatusPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("operational-status: nil\n"))
	}
	if o.RolePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "role", *o.RolePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("role: nil\n"))
	}
	if o.RoutingGroupNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "routing-group-name", *o.RoutingGroupNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("routing-group-name: nil\n"))
	}
	if o.SubnetNamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "subnet-name", *o.SubnetNamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("subnet-name: nil\n"))
	}
	if o.UseFailoverGroupPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "use-failover-group", *o.UseFailoverGroupPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("use-failover-group: nil\n"))
	}
	if o.VserverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver", *o.VserverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver: nil\n"))
	}
	if o.WwpnPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "wwpn", *o.WwpnPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("wwpn: nil\n"))
	}
	return buffer.String()
}

func (o *NetInterfaceInfoType) Address() IpAddressType {
	r := *o.AddressPtr
	return r
}

func (o *NetInterfaceInfoType) SetAddress(newValue IpAddressType) *NetInterfaceInfoType {
	o.AddressPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) AddressFamily() string {
	r := *o.AddressFamilyPtr
	return r
}

func (o *NetInterfaceInfoType) SetAddressFamily(newValue string) *NetInterfaceInfoType {
	o.AddressFamilyPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) AdministrativeStatus() string {
	r := *o.AdministrativeStatusPtr
	return r
}

func (o *NetInterfaceInfoType) SetAdministrativeStatus(newValue string) *NetInterfaceInfoType {
	o.AdministrativeStatusPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) Comment() string {
	r := *o.CommentPtr
	return r
}

func (o *NetInterfaceInfoType) SetComment(newValue string) *NetInterfaceInfoType {
	o.CommentPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) CurrentNode() string {
	r := *o.CurrentNodePtr
	return r
}

func (o *NetInterfaceInfoType) SetCurrentNode(newValue string) *NetInterfaceInfoType {
	o.CurrentNodePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) CurrentPort() string {
	r := *o.CurrentPortPtr
	return r
}

func (o *NetInterfaceInfoType) SetCurrentPort(newValue string) *NetInterfaceInfoType {
	o.CurrentPortPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) DataProtocols() []DataProtocolType {
	r := o.DataProtocolsPtr
	return r
}

func (o *NetInterfaceInfoType) SetDataProtocols(newValue []DataProtocolType) *NetInterfaceInfoType {
	newSlice := make([]DataProtocolType, len(newValue))
	copy(newSlice, newValue)
	o.DataProtocolsPtr = newSlice
	return o
}

func (o *NetInterfaceInfoType) DnsDomainName() DnsZoneType {
	r := *o.DnsDomainNamePtr
	return r
}

func (o *NetInterfaceInfoType) SetDnsDomainName(newValue DnsZoneType) *NetInterfaceInfoType {
	o.DnsDomainNamePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) FailoverGroup() FailoverGroupType {
	r := *o.FailoverGroupPtr
	return r
}

func (o *NetInterfaceInfoType) SetFailoverGroup(newValue FailoverGroupType) *NetInterfaceInfoType {
	o.FailoverGroupPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) FailoverPolicy() string {
	r := *o.FailoverPolicyPtr
	return r
}

func (o *NetInterfaceInfoType) SetFailoverPolicy(newValue string) *NetInterfaceInfoType {
	o.FailoverPolicyPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) FirewallPolicy() string {
	r := *o.FirewallPolicyPtr
	return r
}

func (o *NetInterfaceInfoType) SetFirewallPolicy(newValue string) *NetInterfaceInfoType {
	o.FirewallPolicyPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) ForceSubnetAssociation() bool {
	r := *o.ForceSubnetAssociationPtr
	return r
}

func (o *NetInterfaceInfoType) SetForceSubnetAssociation(newValue bool) *NetInterfaceInfoType {
	o.ForceSubnetAssociationPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) HomeNode() string {
	r := *o.HomeNodePtr
	return r
}

func (o *NetInterfaceInfoType) SetHomeNode(newValue string) *NetInterfaceInfoType {
	o.HomeNodePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) HomePort() string {
	r := *o.HomePortPtr
	return r
}

func (o *NetInterfaceInfoType) SetHomePort(newValue string) *NetInterfaceInfoType {
	o.HomePortPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) InterfaceName() string {
	r := *o.InterfaceNamePtr
	return r
}

func (o *NetInterfaceInfoType) SetInterfaceName(newValue string) *NetInterfaceInfoType {
	o.InterfaceNamePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) IsAutoRevert() bool {
	r := *o.IsAutoRevertPtr
	return r
}

func (o *NetInterfaceInfoType) SetIsAutoRevert(newValue bool) *NetInterfaceInfoType {
	o.IsAutoRevertPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) IsHome() bool {
	r := *o.IsHomePtr
	return r
}

func (o *NetInterfaceInfoType) SetIsHome(newValue bool) *NetInterfaceInfoType {
	o.IsHomePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) IsIpv4LinkLocal() bool {
	r := *o.IsIpv4LinkLocalPtr
	return r
}

func (o *NetInterfaceInfoType) SetIsIpv4LinkLocal(newValue bool) *NetInterfaceInfoType {
	o.IsIpv4LinkLocalPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) LifUuid() UuidType {
	r := *o.LifUuidPtr
	return r
}

func (o *NetInterfaceInfoType) SetLifUuid(newValue UuidType) *NetInterfaceInfoType {
	o.LifUuidPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) ListenForDnsQuery() bool {
	r := *o.ListenForDnsQueryPtr
	return r
}

func (o *NetInterfaceInfoType) SetListenForDnsQuery(newValue bool) *NetInterfaceInfoType {
	o.ListenForDnsQueryPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) Netmask() IpAddressType {
	r := *o.NetmaskPtr
	return r
}

func (o *NetInterfaceInfoType) SetNetmask(newValue IpAddressType) *NetInterfaceInfoType {
	o.NetmaskPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) NetmaskLength() int {
	r := *o.NetmaskLengthPtr
	return r
}

func (o *NetInterfaceInfoType) SetNetmaskLength(newValue int) *NetInterfaceInfoType {
	o.NetmaskLengthPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) OperationalStatus() string {
	r := *o.OperationalStatusPtr
	return r
}

func (o *NetInterfaceInfoType) SetOperationalStatus(newValue string) *NetInterfaceInfoType {
	o.OperationalStatusPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) Role() string {
	r := *o.RolePtr
	return r
}

func (o *NetInterfaceInfoType) SetRole(newValue string) *NetInterfaceInfoType {
	o.RolePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) RoutingGroupName() RoutingGroupType {
	r := *o.RoutingGroupNamePtr
	return r
}

func (o *NetInterfaceInfoType) SetRoutingGroupName(newValue RoutingGroupType) *NetInterfaceInfoType {
	o.RoutingGroupNamePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) SubnetName() SubnetNameType {
	r := *o.SubnetNamePtr
	return r
}

func (o *NetInterfaceInfoType) SetSubnetName(newValue SubnetNameType) *NetInterfaceInfoType {
	o.SubnetNamePtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) UseFailoverGroup() string {
	r := *o.UseFailoverGroupPtr
	return r
}

func (o *NetInterfaceInfoType) SetUseFailoverGroup(newValue string) *NetInterfaceInfoType {
	o.UseFailoverGroupPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) Vserver() string {
	r := *o.VserverPtr
	return r
}

func (o *NetInterfaceInfoType) SetVserver(newValue string) *NetInterfaceInfoType {
	o.VserverPtr = &newValue
	return o
}

func (o *NetInterfaceInfoType) Wwpn() string {
	r := *o.WwpnPtr
	return r
}

func (o *NetInterfaceInfoType) SetWwpn(newValue string) *NetInterfaceInfoType {
	o.WwpnPtr = &newValue
	return o
}

type SnapshotIdType string

type SnapshotInfoType struct {
	XMLName xml.Name `xml:"snapshot-info"`

	AccessTimePtr                        *int                `xml:"access-time"`
	BusyPtr                              *bool               `xml:"busy"`
	ContainsLunClonesPtr                 *bool               `xml:"contains-lun-clones"`
	CumulativePercentageOfTotalBlocksPtr *int                `xml:"cumulative-percentage-of-total-blocks"`
	CumulativePercentageOfUsedBlocksPtr  *int                `xml:"cumulative-percentage-of-used-blocks"`
	CumulativeTotalPtr                   *int                `xml:"cumulative-total"`
	DependencyPtr                        *string             `xml:"dependency"`
	Is7ModeSnapshotPtr                   *bool               `xml:"is-7-mode-snapshot"`
	IsConstituentSnapshotPtr             *bool               `xml:"is-constituent-snapshot"`
	NamePtr                              *string             `xml:"name"`
	PercentageOfTotalBlocksPtr           *int                `xml:"percentage-of-total-blocks"`
	PercentageOfUsedBlocksPtr            *int                `xml:"percentage-of-used-blocks"`
	SnapmirrorLabelPtr                   *string             `xml:"snapmirror-label"`
	SnapshotInstanceUuidPtr              *UUIDType           `xml:"snapshot-instance-uuid"`
	SnapshotOwnersListPtr                []SnapshotOwnerType `xml:"snapshot-owners-list>snapshot-owner"`
	SnapshotVersionUuidPtr               *UUIDType           `xml:"snapshot-version-uuid"`
	StatePtr                             *string             `xml:"state"`
	TotalPtr                             *int                `xml:"total"`
	VolumePtr                            *string             `xml:"volume"`
	VolumeProvenanceUuidPtr              *UUIDType           `xml:"volume-provenance-uuid"`
	VserverPtr                           *string             `xml:"vserver"`
}

func (o *SnapshotInfoType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} // TODO: handle better
	//fmt.Println(string(output))
	return string(output), err
}

func NewSnapshotInfoType() *SnapshotInfoType { return &SnapshotInfoType{} }

func (o SnapshotInfoType) String() string {
	var buffer bytes.Buffer
	if o.AccessTimePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "access-time", *o.AccessTimePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("access-time: nil\n"))
	}
	if o.BusyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "busy", *o.BusyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("busy: nil\n"))
	}
	if o.ContainsLunClonesPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "contains-lun-clones", *o.ContainsLunClonesPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("contains-lun-clones: nil\n"))
	}
	if o.CumulativePercentageOfTotalBlocksPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "cumulative-percentage-of-total-blocks", *o.CumulativePercentageOfTotalBlocksPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("cumulative-percentage-of-total-blocks: nil\n"))
	}
	if o.CumulativePercentageOfUsedBlocksPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "cumulative-percentage-of-used-blocks", *o.CumulativePercentageOfUsedBlocksPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("cumulative-percentage-of-used-blocks: nil\n"))
	}
	if o.CumulativeTotalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "cumulative-total", *o.CumulativeTotalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("cumulative-total: nil\n"))
	}
	if o.DependencyPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "dependency", *o.DependencyPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("dependency: nil\n"))
	}
	if o.Is7ModeSnapshotPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-7-mode-snapshot", *o.Is7ModeSnapshotPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-7-mode-snapshot: nil\n"))
	}
	if o.IsConstituentSnapshotPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "is-constituent-snapshot", *o.IsConstituentSnapshotPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("is-constituent-snapshot: nil\n"))
	}
	if o.NamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name", *o.NamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name: nil\n"))
	}
	if o.PercentageOfTotalBlocksPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-of-total-blocks", *o.PercentageOfTotalBlocksPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-of-total-blocks: nil\n"))
	}
	if o.PercentageOfUsedBlocksPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "percentage-of-used-blocks", *o.PercentageOfUsedBlocksPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("percentage-of-used-blocks: nil\n"))
	}
	if o.SnapmirrorLabelPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapmirror-label", *o.SnapmirrorLabelPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapmirror-label: nil\n"))
	}
	if o.SnapshotInstanceUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-instance-uuid", *o.SnapshotInstanceUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-instance-uuid: nil\n"))
	}
	if o.SnapshotOwnersListPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-owners-list", o.SnapshotOwnersListPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-owners-list: nil\n"))
	}
	if o.SnapshotVersionUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "snapshot-version-uuid", *o.SnapshotVersionUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("snapshot-version-uuid: nil\n"))
	}
	if o.StatePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "state", *o.StatePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("state: nil\n"))
	}
	if o.TotalPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "total", *o.TotalPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("total: nil\n"))
	}
	if o.VolumePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume", *o.VolumePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume: nil\n"))
	}
	if o.VolumeProvenanceUuidPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "volume-provenance-uuid", *o.VolumeProvenanceUuidPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("volume-provenance-uuid: nil\n"))
	}
	if o.VserverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver", *o.VserverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver: nil\n"))
	}
	return buffer.String()
}

func (o *SnapshotInfoType) AccessTime() int {
	r := *o.AccessTimePtr
	return r
}

func (o *SnapshotInfoType) SetAccessTime(newValue int) *SnapshotInfoType {
	o.AccessTimePtr = &newValue
	return o
}

func (o *SnapshotInfoType) Busy() bool {
	r := *o.BusyPtr
	return r
}

func (o *SnapshotInfoType) SetBusy(newValue bool) *SnapshotInfoType {
	o.BusyPtr = &newValue
	return o
}

func (o *SnapshotInfoType) ContainsLunClones() bool {
	r := *o.ContainsLunClonesPtr
	return r
}

func (o *SnapshotInfoType) SetContainsLunClones(newValue bool) *SnapshotInfoType {
	o.ContainsLunClonesPtr = &newValue
	return o
}

func (o *SnapshotInfoType) CumulativePercentageOfTotalBlocks() int {
	r := *o.CumulativePercentageOfTotalBlocksPtr
	return r
}

func (o *SnapshotInfoType) SetCumulativePercentageOfTotalBlocks(newValue int) *SnapshotInfoType {
	o.CumulativePercentageOfTotalBlocksPtr = &newValue
	return o
}

func (o *SnapshotInfoType) CumulativePercentageOfUsedBlocks() int {
	r := *o.CumulativePercentageOfUsedBlocksPtr
	return r
}

func (o *SnapshotInfoType) SetCumulativePercentageOfUsedBlocks(newValue int) *SnapshotInfoType {
	o.CumulativePercentageOfUsedBlocksPtr = &newValue
	return o
}

func (o *SnapshotInfoType) CumulativeTotal() int {
	r := *o.CumulativeTotalPtr
	return r
}

func (o *SnapshotInfoType) SetCumulativeTotal(newValue int) *SnapshotInfoType {
	o.CumulativeTotalPtr = &newValue
	return o
}

func (o *SnapshotInfoType) Dependency() string {
	r := *o.DependencyPtr
	return r
}

func (o *SnapshotInfoType) SetDependency(newValue string) *SnapshotInfoType {
	o.DependencyPtr = &newValue
	return o
}

func (o *SnapshotInfoType) Is7ModeSnapshot() bool {
	r := *o.Is7ModeSnapshotPtr
	return r
}

func (o *SnapshotInfoType) SetIs7ModeSnapshot(newValue bool) *SnapshotInfoType {
	o.Is7ModeSnapshotPtr = &newValue
	return o
}

func (o *SnapshotInfoType) IsConstituentSnapshot() bool {
	r := *o.IsConstituentSnapshotPtr
	return r
}

func (o *SnapshotInfoType) SetIsConstituentSnapshot(newValue bool) *SnapshotInfoType {
	o.IsConstituentSnapshotPtr = &newValue
	return o
}

func (o *SnapshotInfoType) Name() string {
	r := *o.NamePtr
	return r
}

func (o *SnapshotInfoType) SetName(newValue string) *SnapshotInfoType {
	o.NamePtr = &newValue
	return o
}

func (o *SnapshotInfoType) PercentageOfTotalBlocks() int {
	r := *o.PercentageOfTotalBlocksPtr
	return r
}

func (o *SnapshotInfoType) SetPercentageOfTotalBlocks(newValue int) *SnapshotInfoType {
	o.PercentageOfTotalBlocksPtr = &newValue
	return o
}

func (o *SnapshotInfoType) PercentageOfUsedBlocks() int {
	r := *o.PercentageOfUsedBlocksPtr
	return r
}

func (o *SnapshotInfoType) SetPercentageOfUsedBlocks(newValue int) *SnapshotInfoType {
	o.PercentageOfUsedBlocksPtr = &newValue
	return o
}

func (o *SnapshotInfoType) SnapmirrorLabel() string {
	r := *o.SnapmirrorLabelPtr
	return r
}

func (o *SnapshotInfoType) SetSnapmirrorLabel(newValue string) *SnapshotInfoType {
	o.SnapmirrorLabelPtr = &newValue
	return o
}

func (o *SnapshotInfoType) SnapshotInstanceUuid() UUIDType {
	r := *o.SnapshotInstanceUuidPtr
	return r
}

func (o *SnapshotInfoType) SetSnapshotInstanceUuid(newValue UUIDType) *SnapshotInfoType {
	o.SnapshotInstanceUuidPtr = &newValue
	return o
}

func (o *SnapshotInfoType) SnapshotOwnersList() []SnapshotOwnerType {
	r := o.SnapshotOwnersListPtr
	return r
}

func (o *SnapshotInfoType) SetSnapshotOwnersList(newValue []SnapshotOwnerType) *SnapshotInfoType {
	newSlice := make([]SnapshotOwnerType, len(newValue))
	copy(newSlice, newValue)
	o.SnapshotOwnersListPtr = newSlice
	return o
}

func (o *SnapshotInfoType) SnapshotVersionUuid() UUIDType {
	r := *o.SnapshotVersionUuidPtr
	return r
}

func (o *SnapshotInfoType) SetSnapshotVersionUuid(newValue UUIDType) *SnapshotInfoType {
	o.SnapshotVersionUuidPtr = &newValue
	return o
}

func (o *SnapshotInfoType) State() string {
	r := *o.StatePtr
	return r
}

func (o *SnapshotInfoType) SetState(newValue string) *SnapshotInfoType {
	o.StatePtr = &newValue
	return o
}

func (o *SnapshotInfoType) Total() int {
	r := *o.TotalPtr
	return r
}

func (o *SnapshotInfoType) SetTotal(newValue int) *SnapshotInfoType {
	o.TotalPtr = &newValue
	return o
}

func (o *SnapshotInfoType) Volume() string {
	r := *o.VolumePtr
	return r
}

func (o *SnapshotInfoType) SetVolume(newValue string) *SnapshotInfoType {
	o.VolumePtr = &newValue
	return o
}

func (o *SnapshotInfoType) VolumeProvenanceUuid() UUIDType {
	r := *o.VolumeProvenanceUuidPtr
	return r
}

func (o *SnapshotInfoType) SetVolumeProvenanceUuid(newValue UUIDType) *SnapshotInfoType {
	o.VolumeProvenanceUuidPtr = &newValue
	return o
}

func (o *SnapshotInfoType) Vserver() string {
	r := *o.VserverPtr
	return r
}

func (o *SnapshotInfoType) SetVserver(newValue string) *SnapshotInfoType {
	o.VserverPtr = &newValue
	return o
}

type SnapshotOwnerType struct {
	XMLName xml.Name `xml:"snapshot-owner"`

	OwnerPtr *string `xml:"owner"`
}

func (o *SnapshotOwnerType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} // TODO: handle better
	//fmt.Println(string(output))
	return string(output), err
}

func NewSnapshotOwnerType() *SnapshotOwnerType { return &SnapshotOwnerType{} }

func (o SnapshotOwnerType) String() string {
	var buffer bytes.Buffer
	if o.OwnerPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "owner", *o.OwnerPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("owner: nil\n"))
	}
	return buffer.String()
}

func (o *SnapshotOwnerType) Owner() string {
	r := *o.OwnerPtr
	return r
}

func (o *SnapshotOwnerType) SetOwner(newValue string) *SnapshotOwnerType {
	o.OwnerPtr = &newValue
	return o
}

type UUIDType string

type VolumeErrorType struct {
	XMLName xml.Name `xml:"volume-error"`

	ErrnoPtr   *int            `xml:"errno"`
	NamePtr    *VolumeNameType `xml:"name"`
	ReasonPtr  *string         `xml:"reason"`
	VserverPtr *string         `xml:"vserver"`
}

func (o *VolumeErrorType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} // TODO: handle better
	//fmt.Println(string(output))
	return string(output), err
}

func NewVolumeErrorType() *VolumeErrorType { return &VolumeErrorType{} }

func (o VolumeErrorType) String() string {
	var buffer bytes.Buffer
	if o.ErrnoPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "errno", *o.ErrnoPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("errno: nil\n"))
	}
	if o.NamePtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "name", *o.NamePtr))
	} else {
		buffer.WriteString(fmt.Sprintf("name: nil\n"))
	}
	if o.ReasonPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "reason", *o.ReasonPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("reason: nil\n"))
	}
	if o.VserverPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "vserver", *o.VserverPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("vserver: nil\n"))
	}
	return buffer.String()
}

func (o *VolumeErrorType) Errno() int {
	r := *o.ErrnoPtr
	return r
}

func (o *VolumeErrorType) SetErrno(newValue int) *VolumeErrorType {
	o.ErrnoPtr = &newValue
	return o
}

func (o *VolumeErrorType) Name() VolumeNameType {
	r := *o.NamePtr
	return r
}

func (o *VolumeErrorType) SetName(newValue VolumeNameType) *VolumeErrorType {
	o.NamePtr = &newValue
	return o
}

func (o *VolumeErrorType) Reason() string {
	r := *o.ReasonPtr
	return r
}

func (o *VolumeErrorType) SetReason(newValue string) *VolumeErrorType {
	o.ReasonPtr = &newValue
	return o
}

func (o *VolumeErrorType) Vserver() string {
	r := *o.VserverPtr
	return r
}

func (o *VolumeErrorType) SetVserver(newValue string) *VolumeErrorType {
	o.VserverPtr = &newValue
	return o
}

type BlockRangeType struct {
	XMLName xml.Name `xml:"block-range"`

	BlockCountPtr             *int `xml:"block-count"`
	DestinationBlockNumberPtr *int `xml:"destination-block-number"`
	SourceBlockNumberPtr      *int `xml:"source-block-number"`
}

func (o *BlockRangeType) ToXML() (string, error) {
	output, err := xml.MarshalIndent(o, " ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} // TODO: handle better
	//fmt.Println(string(output))
	return string(output), err
}

func NewBlockRangeType() *BlockRangeType { return &BlockRangeType{} }

func (o BlockRangeType) String() string {
	var buffer bytes.Buffer
	if o.BlockCountPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "block-count", *o.BlockCountPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("block-count: nil\n"))
	}
	if o.DestinationBlockNumberPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "destination-block-number", *o.DestinationBlockNumberPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("destination-block-number: nil\n"))
	}
	if o.SourceBlockNumberPtr != nil {
		buffer.WriteString(fmt.Sprintf("%s: %v\n", "source-block-number", *o.SourceBlockNumberPtr))
	} else {
		buffer.WriteString(fmt.Sprintf("source-block-number: nil\n"))
	}
	return buffer.String()
}

func (o *BlockRangeType) BlockCount() int {
	r := *o.BlockCountPtr
	return r
}

func (o *BlockRangeType) SetBlockCount(newValue int) *BlockRangeType {
	o.BlockCountPtr = &newValue
	return o
}

func (o *BlockRangeType) DestinationBlockNumber() int {
	r := *o.DestinationBlockNumberPtr
	return r
}

func (o *BlockRangeType) SetDestinationBlockNumber(newValue int) *BlockRangeType {
	o.DestinationBlockNumberPtr = &newValue
	return o
}

func (o *BlockRangeType) SourceBlockNumber() int {
	r := *o.SourceBlockNumberPtr
	return r
}

func (o *BlockRangeType) SetSourceBlockNumber(newValue int) *BlockRangeType {
	o.SourceBlockNumberPtr = &newValue
	return o
}

// IscsiServiceInfoType is a structure to represent a iscsi-service-info ZAPI object
type IscsiServiceInfoType struct {
	XMLName xml.Name `xml:"iscsi-service-info"`

	AliasNamePtr             *string `xml:"alias-name"`
	IsAvailablePtr           *bool   `xml:"is-available"`
	LoginTimeoutPtr          *int    `xml:"login-timeout"`
	MaxCmdsPerSessionPtr     *int    `xml:"max-cmds-per-session"`
	MaxConnPerSessionPtr     *int    `xml:"max-conn-per-session"`
	MaxErrorRecoveryLevelPtr *int    `xml:"max-error-recovery-level"`
	NodeNamePtr              *string `xml:"node-name"`
	RetainTimeoutPtr         *int    `xml:"retain-timeout"`
	TcpWindowSizePtr         *int    `xml:"tcp-window-size"`
	VserverPtr               *string `xml:"vserver"`
}

func (o *IscsiServiceInfoType) NodeName() string {
	r := *o.NodeNamePtr
	return r
}

func (o *IscsiServiceInfoType) Vserver() string {
	r := *o.VserverPtr
	return r
}
