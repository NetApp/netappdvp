// Copyright 2016 NetApp, Inc. All Rights Reserved.

package ontap

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/blang/semver"
	"github.com/netapp/netappdvp/azgo"
)

const maxZapiRecords int = 0xfffffffe

// DriverConfig holds the configuration data for Driver objects
type DriverConfig struct {
	ManagementLIF string
	SVM           string
	Username      string
	Password      string
}

// Driver is the object to use for interacting with the Filer
type Driver struct {
	config DriverConfig
	zr     *azgo.ZapiRunner
	m      *sync.Mutex
}

// NewDriver is a factory method for creating a new instance
func NewDriver(config DriverConfig) *Driver {
	d := &Driver{
		config: config,
		zr: &azgo.ZapiRunner{
			ManagementLIF: config.ManagementLIF,
			SVM:           config.SVM,
			Username:      config.Username,
			Password:      config.Password,
			Secure:        true,
		},
		m: &sync.Mutex{},
	}
	return d
}

// GetClonedZapiRunner returns a clone of the ZapiRunner configured on this driver.
func (d Driver) GetClonedZapiRunner() *azgo.ZapiRunner {
	clone := new(azgo.ZapiRunner)
	*clone = *d.zr
	return clone
}

// GetNontunneledZapiRunner returns a clone of the ZapiRunner configured on this driver with the SVM field cleared so ZAPI calls
// made with the resulting runner aren't tunneled.  Note that the calls could still go directly to either a cluster or
// vserver management LIF.
func (d Driver) GetNontunneledZapiRunner() *azgo.ZapiRunner {
	clone := new(azgo.ZapiRunner)
	*clone = *d.zr
	clone.SVM = ""
	return clone
}

// NewZapiError accepts the Response value from any AZGO call, extracts the status, reason, and errno values, and returns a ZapiError.
// TODO: Replace reflection with relevant enhancements in AZGO generator.
func NewZapiError(response interface{}) (err ZapiError) {

	defer func() {
		if r := recover(); r != nil {
			err = ZapiError{}
		}
	}()

	val := reflect.ValueOf(response)

	err = ZapiError{
		val.FieldByName("ResultStatusAttr").String(),
		val.FieldByName("ResultReasonAttr").String(),
		val.FieldByName("ResultErrnoAttr").String(),
	}

	return
}

// ZapiError encapsulates the status, reason, and errno values from a ZAPI invocation, and it provides helper methods for detecting
// common error conditions.
type ZapiError struct {
	status string
	reason string
	code   string
}

func (e ZapiError) IsPassed() bool {
	return e.status == "passed"
}
func (e ZapiError) Error() string {
	if e.IsPassed() {
		return "API status: passed"
	}
	return fmt.Sprintf("API status: %s, Reason: %s, Code: %s", e.status, e.reason, e.code)
}
func (e ZapiError) IsPrivilegeError() bool {
	return e.code == azgo.EAPIPRIVILEGE
}
func (e ZapiError) IsScopeError() bool {
	return e.code == azgo.EAPIPRIVILEGE || e.code == azgo.EAPINOTFOUND
}

/////////////////////////////////////////////////////////////////////////////
// API feature operations BEGIN

type ontapApiFeature string

// Define new version-specific feature constants here
const (
	MINIMUM_ONTAPI_VERSION ontapApiFeature = "MINIMUM_ONTAPI_VERSION"
	VSERVER_SHOW_AGGR      ontapApiFeature = "VSERVER_SHOW_AGGR"
)

// Indicate the minimum Ontapi version for each feature here
var ontapAPIFeatures = map[ontapApiFeature]semver.Version{
	MINIMUM_ONTAPI_VERSION: semver.MustParse("1.30.0"), // cDOT 8.3.0
	VSERVER_SHOW_AGGR:      semver.MustParse("1.100.0"),
}

// SupportsApiFeature returns true if the Ontapi version supports the supplied feature
func (d Driver) SupportsApiFeature(feature ontapApiFeature) bool {

	ontapiVersion, err := d.SystemGetOntapiVersion()
	if err != nil {
		return false
	}

	ontapiSemVer, err := semver.Make(fmt.Sprintf("%s.0", ontapiVersion))
	if err != nil {
		return false
	}

	if minVersion, ok := ontapAPIFeatures[feature]; ok {
		return ontapiSemVer.GTE(minVersion)
	} else {
		return false
	}
}

// API feature operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// IGROUP operations BEGIN

// IgroupCreate creates the specified initiator group
// equivalent to filer::> igroup create docker -vserver iscsi_vs -protocol iscsi -ostype linux
func (d Driver) IgroupCreate(initiatorGroupName, initiatorGroupType, osType string) (response azgo.IgroupCreateResponse, err error) {
	response, err = azgo.NewIgroupCreateRequest().
		SetInitiatorGroupName(initiatorGroupName).
		SetInitiatorGroupType(initiatorGroupType).
		SetOsType(osType).
		ExecuteUsing(d.zr)
	return
}

// IgroupAdd adds an initiator to an initiator group
// equivalent to filer::> igroup add -vserver iscsi_vs -igroup docker -initiator iqn.1993-08.org.debian:01:9031309bbebd
func (d Driver) IgroupAdd(initiatorGroupName, initiator string) (response azgo.IgroupAddResponse, err error) {
	response, err = azgo.NewIgroupAddRequest().
		SetInitiatorGroupName(initiatorGroupName).
		SetInitiator(initiator).
		ExecuteUsing(d.zr)
	return
}

// IgroupRemove removes an initiator from an initiator group
func (d Driver) IgroupRemove(initiatorGroupName, initiator string, force bool) (response azgo.IgroupRemoveResponse, err error) {
	response, err = azgo.NewIgroupRemoveRequest().
		SetInitiatorGroupName(initiatorGroupName).
		SetInitiator(initiator).
		SetForce(force).
		ExecuteUsing(d.zr)
	return
}

// IgroupDestroy destroys an initiator group
func (d Driver) IgroupDestroy(initiatorGroupName string) (response azgo.IgroupDestroyResponse, err error) {
	response, err = azgo.NewIgroupDestroyRequest().
		SetInitiatorGroupName(initiatorGroupName).
		ExecuteUsing(d.zr)
	return
}

// IgroupList lists initiator groups
func (d Driver) IgroupList() (response azgo.IgroupGetIterResponse, err error) {
	response, err = azgo.NewIgroupGetIterRequest().
		SetMaxRecords(maxZapiRecords). // Is there any value in iterating?
		ExecuteUsing(d.zr)
	return
}

// IGROUP operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// LUN operations BEGIN

// LunCreate creates a lun with the specified attributes
// equivalent to filer::> lun create -vserver iscsi_vs -path /vol/v/lun1 -size 1g -ostype linux -space-reserve disabled
func (d Driver) LunCreate(lunPath string, sizeInBytes int, osType string, spaceReserved bool) (response azgo.LunCreateBySizeResponse, err error) {
	response, err = azgo.NewLunCreateBySizeRequest().
		SetPath(lunPath).
		SetSize(sizeInBytes).
		SetOstype(osType).
		SetSpaceReservationEnabled(spaceReserved).
		ExecuteUsing(d.zr)
	return
}

// LunGetSerialNumber returns the serial# for a lun
func (d Driver) LunGetSerialNumber(lunPath string) (response azgo.LunGetSerialNumberResponse, err error) {
	response, err = azgo.NewLunGetSerialNumberRequest().
		SetPath(lunPath).
		ExecuteUsing(d.zr)
	return
}

// LunMap maps a lun to an id in an initiator group
// equivalent to filer::> lun map -vserver iscsi_vs -path /vol/v/lun1 -igroup docker -lun-id 0
func (d Driver) LunMap(initiatorGroupName, lunPath string, lunID int) (response azgo.LunMapResponse, err error) {
	response, err = azgo.NewLunMapRequest().
		SetInitiatorGroup(initiatorGroupName).
		SetPath(lunPath).
		SetLunId(lunID).
		ExecuteUsing(d.zr)
	return
}

// LunMapListInfo returns lun mapping information for the specified lun
// equivalent to filer::> lun mapped show -vserver iscsi_vs -path /vol/v/lun0
func (d Driver) LunMapListInfo(lunPath string) (response azgo.LunMapListInfoResponse, err error) {
	response, err = azgo.NewLunMapListInfoRequest().
		SetPath(lunPath).
		ExecuteUsing(d.zr)
	return
}

// LunOffline offlines a lun
// equivalent to filer::> lun offline -vserver iscsi_vs -path /vol/v/lun0
func (d Driver) LunOffline(lunPath string) (response azgo.LunOfflineResponse, err error) {
	response, err = azgo.NewLunOfflineRequest().
		SetPath(lunPath).
		ExecuteUsing(d.zr)
	return
}

// LunOnline onlines a lun
// equivalent to filer::> lun online -vserver iscsi_vs -path /vol/v/lun0
func (d Driver) LunOnline(lunPath string) (response azgo.LunOnlineResponse, err error) {
	response, err = azgo.NewLunOnlineRequest().
		SetPath(lunPath).
		ExecuteUsing(d.zr)
	return
}

// LunDestroy destroys a lun
// equivalent to filer::> lun destroy -vserver iscsi_vs -path /vol/v/lun0
func (d Driver) LunDestroy(lunPath string) (response azgo.LunDestroyResponse, err error) {
	response, err = azgo.NewLunDestroyRequest().
		SetPath(lunPath).
		ExecuteUsing(d.zr)
	return
}

// LUN operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// VOLUME operations BEGIN

// VolumeCreate creates a volume with the specified options
// equivalent to filer::> volume create -vserver iscsi_vs -volume v -aggregate aggr1 -size 1g -state online -type RW -policy default -unix-permissions ---rwxr-xr-x -space-guarantee none -snapshot-policy none -security-style unix
func (d Driver) VolumeCreate(name, aggregateName, size, spaceReserve, snapshotPolicy, unixPermissions, exportPolicy, securityStyle string) (response azgo.VolumeCreateResponse, err error) {
	response, err = azgo.NewVolumeCreateRequest().
		SetVolume(name).
		SetContainingAggrName(aggregateName).
		SetSize(size).
		SetSpaceReserve(spaceReserve).
		SetSnapshotPolicy(snapshotPolicy).
		SetUnixPermissions(unixPermissions).
		SetExportPolicy(exportPolicy).
		SetVolumeSecurityStyle(securityStyle).
		ExecuteUsing(d.zr)
	return
}

// VolumeCloneCreate clones a volume from a snapshot
func (d Driver) VolumeCloneCreate(name, source, snapshot string) (response azgo.VolumeCloneCreateResponse, err error) {
	response, err = azgo.NewVolumeCloneCreateRequest().
		SetVolume(name).
		SetParentVolume(source).
		SetParentSnapshot(snapshot).
		ExecuteUsing(d.zr)
	return
}

// VolumeCloneSplitStart splits a cloned volume from its parent
func (d Driver) VolumeCloneSplitStart(name string) (response azgo.VolumeCloneSplitStartResponse, err error) {
	response, err = azgo.NewVolumeCloneSplitStartRequest().
		SetVolume(name).
		ExecuteUsing(d.zr)
	return
}

// VolumeDisableSnapshotDirectoryAccess disables access to the ".snapshot" directory
// Disable '.snapshot' to allow official mysql container's chmod-in-init to work
func (d Driver) VolumeDisableSnapshotDirectoryAccess(name string) (response azgo.VolumeModifyIterResponse, err error) {
	ssattr := azgo.NewVolumeSnapshotAttributesType().SetSnapdirAccessEnabled(false)
	volattr := azgo.NewVolumeAttributesType().SetVolumeSnapshotAttributes(*ssattr)
	volidattr := azgo.NewVolumeIdAttributesType().SetName(azgo.VolumeNameType(name))
	queryattr := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*volidattr)

	response, err = azgo.NewVolumeModifyIterRequest().
		SetQuery(*queryattr).
		SetAttributes(*volattr).
		ExecuteUsing(d.zr)
	return
}

// VolumeSize retrieves the size of the specified volume
func (d Driver) VolumeSize(name string) (response azgo.VolumeSizeResponse, err error) {
	response, err = azgo.NewVolumeSizeRequest().
		SetVolume(name).
		ExecuteUsing(d.zr)
	return
}

// VolumeMount mounts a volume at the specified junction
func (d Driver) VolumeMount(name, junctionPath string) (response azgo.VolumeMountResponse, err error) {
	response, err = azgo.NewVolumeMountRequest().
		SetVolumeName(name).
		SetJunctionPath(junctionPath).
		ExecuteUsing(d.zr)
	return
}

// VolumeUnmount unmounts a volume from the specified junction
func (d Driver) VolumeUnmount(name string, force bool) (response azgo.VolumeUnmountResponse, err error) {
	response, err = azgo.NewVolumeUnmountRequest().
		SetVolumeName(name).
		SetForce(force).
		ExecuteUsing(d.zr)
	return
}

// VolumeOffline offlines a volume
func (d Driver) VolumeOffline(name string) (response azgo.VolumeOfflineResponse, err error) {
	response, err = azgo.NewVolumeOfflineRequest().
		SetName(name).
		ExecuteUsing(d.zr)
	return
}

// VolumeDestroy destroys a volume
func (d Driver) VolumeDestroy(name string, force bool) (response azgo.VolumeDestroyResponse, err error) {
	response, err = azgo.NewVolumeDestroyRequest().
		SetName(name).
		SetUnmountAndOffline(force).
		ExecuteUsing(d.zr)
	return
}

// VolumeList lists volumes
func (d Driver) VolumeList(prefix string) (response azgo.VolumeGetIterResponse, err error) {
	viat := azgo.NewVolumeIdAttributesType().SetName(azgo.VolumeNameType(prefix + "*"))
	query := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*viat)

	response, err = azgo.NewVolumeGetIterRequest().
		SetMaxRecords(maxZapiRecords). // Is there any value in iterating?
		SetQuery(*query).
		ExecuteUsing(d.zr)
	return
}

// VOLUME operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// SNAPSHOT operations BEGIN

// SnapshotCreate creates a snapshot of a volume
func (d Driver) SnapshotCreate(name, volumeName string) (response azgo.SnapshotCreateResponse, err error) {
	response, err = azgo.NewSnapshotCreateRequest().
		SetSnapshot(name).
		SetVolume(volumeName).
		ExecuteUsing(d.zr)
	return
}

// SnapshotGetByVolume returns the list of snapshots associated with a volume
func (d Driver) SnapshotGetByVolume(volumeName string) (response azgo.SnapshotGetIterResponse, err error) {
	query := azgo.NewSnapshotInfoType().SetVolume(volumeName)

	response, err = azgo.NewSnapshotGetIterRequest().
		SetMaxRecords(maxZapiRecords). // Is there any value in iterating?
		SetQuery(*query).
		ExecuteUsing(d.zr)
	return
}

// SNAPSHOT operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// ISCSI operations BEGIN

// IscsiServiceGetIterRequest returns information about an iSCSI target
func (d Driver) IscsiServiceGetIterRequest() (response azgo.IscsiServiceGetIterResponse, err error) {
	response, err = azgo.NewIscsiServiceGetIterRequest().
		SetMaxRecords(maxZapiRecords). // Is there any value in iterating?
		ExecuteUsing(d.zr)
	return
}

// ISCSI operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// VSERVER operations BEGIN

// VserverGetIterRequest returns the vservers on the system
// equivalent to filer::> vserver show
func (d Driver) VserverGetIterRequest() (response azgo.VserverGetIterResponse, err error) {
	response, err = azgo.NewVserverGetIterRequest().
		SetMaxRecords(maxZapiRecords). // Is there any value in iterating?
		ExecuteUsing(d.zr)
	return
}

// GetVserverAggregateNames returns an array of names of the aggregates assigned to the configured vserver.
// The vserver-get-iter API works with either cluster or vserver scope, so the ZAPI runner may or may not
// be configured for tunneling; using the query parameter ensures we address only the configured vserver.
func (d Driver) GetVserverAggregateNames() ([]string, error) {

	// Get just the SVM of interest
	query := azgo.NewVserverInfoType()
	query.SetVserverName(d.config.SVM)

	response, err := azgo.NewVserverGetIterRequest().SetMaxRecords(maxZapiRecords).SetQuery(*query).ExecuteUsing(d.zr)
	if err != nil {
		return nil, err
	}
	if response.Result.NumRecords() != 1 {
		return nil, fmt.Errorf("Could not find SVM %s.", d.config.SVM)
	}

	// Get the aggregates assigned to the SVM
	aggrNames := make([]string, 0, 10)
	for _, vserver := range response.Result.AttributesList() {
		aggrList := vserver.VserverAggrInfoList()
		for _, aggr := range aggrList {
			aggrNames = append(aggrNames, string(aggr.AggrName()))
		}
	}

	return aggrNames, nil
}

// VserverShowAggrGetIterRequest returns the aggregates on the vserver.  Requires ONTAP 9 or later.
// equivalent to filer::> vserver show-aggregates
func (d Driver) VserverShowAggrGetIterRequest() (response azgo.VserverShowAggrGetIterResponse, err error) {

	response, err = azgo.NewVserverShowAggrGetIterRequest().
		SetMaxRecords(maxZapiRecords).
		ExecuteUsing(d.zr)
	return
}

// VSERVER operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// AGGREGATE operations BEGIN

// AggrGetIterRequest returns the aggregates on the system
// equivalent to filer::> storage aggregate show
func (d Driver) AggrGetIterRequest() (response azgo.AggrGetIterResponse, err error) {

	// If we tunnel to an SVM, which is the default case, this API will never work.
	// It will still fail if the non-tunneled ZapiRunner addresses a vserver management LIF,
	// but that possibility must be handled by the caller.
	zr := d.GetNontunneledZapiRunner()

	response, err = azgo.NewAggrGetIterRequest().
		SetMaxRecords(maxZapiRecords).
		ExecuteUsing(zr)
	return
}

// AGGREGATE operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// MISC operations BEGIN

// NetInterfaceGet returns the list of network interfaces with associated metadata
// equivalent to filer::> net interface list
func (d Driver) NetInterfaceGet() (response azgo.NetInterfaceGetIterResponse, err error) {
	response, err = azgo.NewNetInterfaceGetIterRequest().
		SetMaxRecords(maxZapiRecords). // Is there any value in iterating?
		ExecuteUsing(d.zr)
	return
}

// SystemGetVersion returns the system version
// equivalent to filer::> version
func (d Driver) SystemGetVersion() (response azgo.SystemGetVersionResponse, err error) {
	response, err = azgo.NewSystemGetVersionRequest().ExecuteUsing(d.zr)
	return
}

// GetOntapiVersion gets the ONTAPI version using the credentials, and caches & returns the result.
func (d Driver) SystemGetOntapiVersion() (string, error) {

	if d.zr.OntapiVersion == "" {
		result, err := azgo.NewSystemGetOntapiVersionRequest().ExecuteUsing(d.zr)
		if err != nil {
			return "", err
		} else if result.Result.ResultStatusAttr != "passed" {
			return "", fmt.Errorf("Could not read ONTAPI version. %s", result.Result.ResultReasonAttr)
		}

		major := result.Result.MajorVersion()
		minor := result.Result.MinorVersion()
		d.zr.OntapiVersion = fmt.Sprintf("%d.%d", major, minor)
	}

	return d.zr.OntapiVersion, nil
}

// EmsAutosupportLog generates an auto support message with the supplied parameters
func (d Driver) EmsAutosupportLog(
	appVersion string,
	autoSupport bool,
	category string,
	computerName string,
	eventDescription string,
	eventID int,
	eventSource string,
	logLevel int) (response azgo.EmsAutosupportLogResponse, err error) {

	response, err = azgo.NewEmsAutosupportLogRequest().
		SetAutoSupport(autoSupport).
		SetAppVersion(appVersion).
		SetCategory(category).
		SetComputerName(computerName).
		SetEventDescription(eventDescription).
		SetEventId(eventID).
		SetEventSource(eventSource).
		SetLogLevel(logLevel).
		ExecuteUsing(d.zr)
	return
}

// MISC operations END
/////////////////////////////////////////////////////////////////////////////
