// Copyright 2016 NetApp, Inc. All Rights Reserved.

package ontap

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/blang/semver"
	"github.com/netapp/netappdvp/azgo"
)

const defaultZapiRecords int = 100

// DriverConfig holds the configuration data for Driver objects
type DriverConfig struct {
	ManagementLIF   string
	SVM             string
	Username        string
	Password        string
	DebugTraceFlags map[string]bool
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
			ManagementLIF:   config.ManagementLIF,
			SVM:             config.SVM,
			Username:        config.Username,
			Password:        config.Password,
			Secure:          true,
			DebugTraceFlags: config.DebugTraceFlags,
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

// NewZapiError accepts the Response value from any AZGO call, extracts the status, reason, and errno values,
// and returns a ZapiError.  The interface passed in may either be a Response object, or the always-embedded
// Result object where the error info exists.
// TODO: Replace reflection with relevant enhancements in AZGO generator.
func NewZapiError(zapiResult interface{}) (err ZapiError) {

	defer func() {
		if r := recover(); r != nil {
			err = ZapiError{}
		}
	}()

	// A ZAPI Result struct works as-is, but a ZAPI Response struct must have its
	// embedded Result struct extracted via reflection.
	val := reflect.ValueOf(zapiResult)
	if testResult := val.FieldByName("Result"); testResult.IsValid() {
		zapiResult = testResult.Interface()
		val = reflect.ValueOf(zapiResult)
	}

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
func (e ZapiError) Reason() string {
	return e.reason
}
func (e ZapiError) Code() string {
	return e.code
}

// GetError accepts both an error and the Response value from an AZGO invocation.
// If error is non-nil, it is returned as is.  Otherwise, the Response value is
// probed for an error returned by ZAPI; if one is found, a ZapiError error object
// is returned.  If no failures are detected, the method returns nil.  The interface
// passed in may either be a Response object, or the always-embedded Result object
// where the error info exists.
func GetError(zapiResult interface{}, errorIn error) (errorOut error) {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Panic in ontap#GetError. %v", zapiResult)
			errorOut = ZapiError{}
		}
	}()

	// A ZAPI Result struct works as-is, but a ZAPI Response struct must have its
	// embedded Result struct extracted via reflection.
	val := reflect.ValueOf(zapiResult)
	if testResult := val.FieldByName("Result"); testResult.IsValid() {
		zapiResult = testResult.Interface()
	}

	errorOut = nil

	if errorIn != nil {
		errorOut = errorIn
	} else if zerr := NewZapiError(zapiResult); !zerr.IsPassed() {
		errorOut = zerr
	}

	return
}

/////////////////////////////////////////////////////////////////////////////
// API feature operations BEGIN

type ontapApiFeature string

// Define new version-specific feature constants here
const (
	MINIMUM_ONTAPI_VERSION   ontapApiFeature = "MINIMUM_ONTAPI_VERSION"
	VSERVER_SHOW_AGGR        ontapApiFeature = "VSERVER_SHOW_AGGR"
	FLEX_GROUPS              ontapApiFeature = "FLEX_GROUPS"
	NETAPP_VOLUME_ENCRYPTION ontapApiFeature = "NETAPP_VOLUME_ENCRYPTION"
)

// Indicate the minimum Ontapi version for each feature here
var ontapAPIFeatures = map[ontapApiFeature]semver.Version{
	MINIMUM_ONTAPI_VERSION:   semver.MustParse("1.30.0"),  // cDOT 8.3.0
	VSERVER_SHOW_AGGR:        semver.MustParse("1.100.0"), // cDOT 9.0.0
	FLEX_GROUPS:              semver.MustParse("1.100.0"), // cDOT 9.0.0
	NETAPP_VOLUME_ENCRYPTION: semver.MustParse("1.110.0"), // cDOT 9.1.0
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
		SetMaxRecords(defaultZapiRecords).
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

// LunMap maps a LUN in an initiator group, allowing ONTAP to choose an available LUN ID
// equivalent to filer::> lun map -vserver iscsi_vs -path /vol/v/lun1 -igroup docker
func (d Driver) LunMapAutoID(initiatorGroupName, lunPath string) (response azgo.LunMapResponse, err error) {
	response, err = azgo.NewLunMapRequest().
		SetInitiatorGroup(initiatorGroupName).
		SetPath(lunPath).
		ExecuteUsing(d.zr)
	return
}

func (d Driver) LunMapIfNotMapped(initiatorGroupName, lunPath string) (int, error) {

	// Read LUN maps to see if the LUN is already mapped to the igroup
	lunMapListResponse, err := d.LunMapListInfo(lunPath)
	if err != nil {
		return -1, fmt.Errorf("Problem reading maps for LUN %s: %v", lunPath, err)
	} else if lunMapListResponse.Result.ResultStatusAttr != "passed" {
		return -1, fmt.Errorf("Problem reading maps for LUN %s: %+v", lunPath, lunMapListResponse.Result)
	}

	lunID := 0
	alreadyMapped := false
	for _, igroup := range lunMapListResponse.Result.InitiatorGroups() {
		if igroup.InitiatorGroupName() == initiatorGroupName {

			lunID = igroup.LunId()
			alreadyMapped = true

			log.WithFields(log.Fields{
				"lun":    lunPath,
				"igroup": initiatorGroupName,
				"id":     lunID,
			}).Debug("LUN already mapped.")

			break
		}
	}

	// Map IFF not already mapped
	if !alreadyMapped {
		lunMapResponse, err := d.LunMapAutoID(initiatorGroupName, lunPath)
		if err != nil {
			return -1, fmt.Errorf("Problem mapping LUN %s: %v", lunPath, err)
		} else if lunMapResponse.Result.ResultStatusAttr != "passed" {
			return -1, fmt.Errorf("Problem mapping LUN %s: %+v", lunPath, lunMapResponse.Result)
		}

		lunID = lunMapResponse.Result.LunIdAssigned()

		log.WithFields(log.Fields{
			"lun":    lunPath,
			"igroup": initiatorGroupName,
			"id":     lunID,
		}).Debug("LUN mapped.")
	}

	return lunID, nil
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

// LunSetAttribute sets a named attribute for a given LUN.
func (d Driver) LunSetAttribute(lunPath, name, value string) (response azgo.LunSetAttributeResponse, err error) {
	response, err = azgo.NewLunSetAttributeRequest().
		SetPath(lunPath).
		SetName(name).
		SetValue(value).
		ExecuteUsing(d.zr)
	return
}

// LunGetAttribute gets a named attribute for a given LUN.
func (d Driver) LunGetAttribute(lunPath, name string) (response azgo.LunGetAttributeResponse, err error) {
	response, err = azgo.NewLunGetAttributeRequest().
		SetPath(lunPath).
		SetName(name).
		ExecuteUsing(d.zr)
	return
}

// LunGet returns all relevant details for a single LUN
// equivalent to filer::> lun show
func (d Driver) LunGet(path string) (azgo.LunInfoType, error) {

	// Limit the LUNs to the one matching the path
	query := azgo.NewLunInfoType().SetPath(path)

	// Limit the returned data to only the data relevant to containers
	desiredAttributes := azgo.NewLunInfoType().
		SetPath("").
		SetVolume("").
		SetSize(0)

	response, err := azgo.NewLunGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)

	if err != nil {
		return azgo.LunInfoType{}, err
	} else if response.Result.NumRecords() == 0 {
		return azgo.LunInfoType{}, fmt.Errorf("LUN %s not found.", path)
	} else if response.Result.NumRecords() > 1 {
		return azgo.LunInfoType{}, fmt.Errorf("More than one LUN %s found.", path)
	}

	return response.Result.AttributesList()[0], nil
}

// LunGetAll returns all relevant details for all LUNs whose paths match the supplied pattern
// equivalent to filer::> lun show
func (d Driver) LunGetAll(pathPattern string) (response azgo.LunGetIterResponse, err error) {

	// Limit the LUNs to those matching the path pattern
	query := azgo.NewLunInfoType().SetPath(pathPattern)

	// Limit the returned data to only the data relevant to containers
	desiredAttributes := azgo.NewLunInfoType().
		SetPath("").
		SetVolume("").
		SetSize(0)

	response, err = azgo.NewLunGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// LUN operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// VOLUME operations BEGIN

// VolumeCreate creates a volume with the specified options
// equivalent to filer::> volume create -vserver iscsi_vs -volume v -aggregate aggr1 -size 1g -state online -type RW -policy default -unix-permissions ---rwxr-xr-x -space-guarantee none -snapshot-policy none -security-style unix -encrypt false
func (d Driver) VolumeCreate(name, aggregateName, size, spaceReserve, snapshotPolicy, unixPermissions,
	exportPolicy, securityStyle string, encrypt *bool) (response azgo.VolumeCreateResponse, err error) {
	request := azgo.NewVolumeCreateRequest().
		SetVolume(name).
		SetContainingAggrName(aggregateName).
		SetSize(size).
		SetSpaceReserve(spaceReserve).
		SetSnapshotPolicy(snapshotPolicy).
		SetUnixPermissions(unixPermissions).
		SetExportPolicy(exportPolicy).
		SetVolumeSecurityStyle(securityStyle)

	// Don't send 'encrypt' unless needed, as pre-9.1 ONTAP won't accept it.
	if encrypt != nil {
		request.SetEncrypt(*encrypt)
	}

	response, err = request.ExecuteUsing(d.zr)
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

// VolumeExists tests for the existence of a Flexvol
func (d Driver) VolumeExists(name string) (bool, error) {
	response, err := azgo.NewVolumeSizeRequest().
		SetVolume(name).
		ExecuteUsing(d.zr)

	if err != nil {
		return false, err
	}

	if zerr := NewZapiError(response); !zerr.IsPassed() {
		switch zerr.Code() {
		case azgo.EOBJECTNOTFOUND, azgo.EVOLUMEDOESNOTEXIST:
			return false, nil
		default:
			return false, zerr
		}
	}

	return true, nil
}

// VolumeSize retrieves the size of the specified volume
func (d Driver) VolumeSize(name string) (response azgo.VolumeSizeResponse, err error) {
	response, err = azgo.NewVolumeSizeRequest().
		SetVolume(name).
		ExecuteUsing(d.zr)
	return
}

// SetVolumeSize sets the size of the specified volume
func (d Driver) SetVolumeSize(name, newSize string) (response azgo.VolumeSizeResponse, err error) {
	response, err = azgo.NewVolumeSizeRequest().
		SetVolume(name).
		SetNewSize(newSize).
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

// VolumeGet returns all relevant details for a single Flexvol
// equivalent to filer::> volume show
func (d Driver) VolumeGet(name string) (azgo.VolumeAttributesType, error) {

	// Limit the Flexvols to the one matching the name
	queryVolIdAttrs := azgo.NewVolumeIdAttributesType().SetName(azgo.VolumeNameType(name))
	if d.SupportsApiFeature(FLEX_GROUPS) {
		queryVolIdAttrs.SetStyleExtended("flexvol")
	}
	query := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*queryVolIdAttrs)

	response, err := azgo.NewVolumeGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		ExecuteUsing(d.zr)

	if err != nil {
		return azgo.VolumeAttributesType{}, err
	} else if response.Result.NumRecords() == 0 {
		return azgo.VolumeAttributesType{}, fmt.Errorf("Flexvol %s not found.", name)
	} else if response.Result.NumRecords() > 1 {
		return azgo.VolumeAttributesType{}, fmt.Errorf("More than one Flexvol %s found.", name)
	}

	return response.Result.AttributesList()[0], nil
}

// VolumeGetAll returns all relevant details for all FlexVols whose names match the supplied prefix
// equivalent to filer::> volume show
func (d Driver) VolumeGetAll(prefix string) (response azgo.VolumeGetIterResponse, err error) {

	// Limit the Flexvols to those matching the name prefix
	queryVolIdAttrs := azgo.NewVolumeIdAttributesType().SetName(azgo.VolumeNameType(prefix + "*"))
	if d.SupportsApiFeature(FLEX_GROUPS) {
		queryVolIdAttrs.SetStyleExtended("flexvol")
	}
	query := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*queryVolIdAttrs)

	// Limit the returned data to only the data relevant to containers
	desiredVolExportAttrs := azgo.NewVolumeExportAttributesType().
		SetPolicy("")
	desiredVolIdAttrs := azgo.NewVolumeIdAttributesType().
		SetName("").
		SetContainingAggregateName("")
	desiredVolSecurityUnixAttrs := azgo.NewVolumeSecurityUnixAttributesType().
		SetPermissions("")
	desiredVolSecurityAttrs := azgo.NewVolumeSecurityAttributesType().
		SetVolumeSecurityUnixAttributes(*desiredVolSecurityUnixAttrs)
	desiredVolSpaceAttrs := azgo.NewVolumeSpaceAttributesType().
		SetSize(0)
	desiredVolSnapshotAttrs := azgo.NewVolumeSnapshotAttributesType().
		SetSnapdirAccessEnabled(true).
		SetSnapshotPolicy("")

	desiredAttributes := azgo.NewVolumeAttributesType().
		SetVolumeExportAttributes(*desiredVolExportAttrs).
		SetVolumeIdAttributes(*desiredVolIdAttrs).
		SetVolumeSecurityAttributes(*desiredVolSecurityAttrs).
		SetVolumeSpaceAttributes(*desiredVolSpaceAttrs).
		SetVolumeSnapshotAttributes(*desiredVolSnapshotAttrs)

	response, err = azgo.NewVolumeGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// VolumeList returns the names of all Flexvols whose names match the supplied prefix
func (d Driver) VolumeList(prefix string) (response azgo.VolumeGetIterResponse, err error) {

	// Limit the Flexvols to those matching the name prefix
	queryVolIdAttrs := azgo.NewVolumeIdAttributesType().SetName(azgo.VolumeNameType(prefix + "*"))
	if d.SupportsApiFeature(FLEX_GROUPS) {
		queryVolIdAttrs.SetStyleExtended("flexvol")
	}
	query := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*queryVolIdAttrs)

	// Limit the returned data to only the Flexvol names
	desiredVolIdAttrs := azgo.NewVolumeIdAttributesType().SetName("")
	desiredAttributes := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*desiredVolIdAttrs)

	response, err = azgo.NewVolumeGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// VolumeListByAttrs returns the names of all Flexvols matching the specified attributes
func (d Driver) VolumeListByAttrs(
	prefix, aggregate, spaceReserve, snapshotPolicy string, snapshotDir bool, encrypt *bool,
) (response azgo.VolumeGetIterResponse, err error) {

	// Limit the Flexvols to those matching the specified attributes
	queryVolIdAttrs := azgo.NewVolumeIdAttributesType().
		SetName(azgo.VolumeNameType(prefix + "*")).
		SetContainingAggregateName(aggregate)
	if d.SupportsApiFeature(FLEX_GROUPS) {
		queryVolIdAttrs.SetStyleExtended("flexvol")
	}
	queryVolSpaceAttrs := azgo.NewVolumeSpaceAttributesType().
		SetSpaceGuarantee(spaceReserve)
	queryVolSnapshotAttrs := azgo.NewVolumeSnapshotAttributesType().
		SetSnapshotPolicy(snapshotPolicy).
		SetSnapdirAccessEnabled(snapshotDir)
	query := azgo.NewVolumeAttributesType().
		SetVolumeIdAttributes(*queryVolIdAttrs).
		SetVolumeSpaceAttributes(*queryVolSpaceAttrs).
		SetVolumeSnapshotAttributes(*queryVolSnapshotAttrs)

	if encrypt != nil {
		query.SetEncrypt(*encrypt)
	}

	// Limit the returned data to only the Flexvol names
	desiredVolIdAttrs := azgo.NewVolumeIdAttributesType().SetName("")
	desiredAttributes := azgo.NewVolumeAttributesType().SetVolumeIdAttributes(*desiredVolIdAttrs)

	response, err = azgo.NewVolumeGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// VolumeGetRootName gets the name of the root volume of a vserver
func (d Driver) VolumeGetRootName() (response azgo.VolumeGetRootNameResponse, err error) {
	response, err = azgo.NewVolumeGetRootNameRequest().
		ExecuteUsing(d.zr)
	return
}

// VOLUME operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// QTREE operations BEGIN

// QtreeCreate creates a qtree with the specified options
// equivalent to filer::> qtree create -vserver ndvp_vs -volume v -qtree q -export-policy default -unix-permissions ---rwxr-xr-x -security-style unix
func (d Driver) QtreeCreate(name, volumeName, unixPermissions, exportPolicy, securityStyle string) (response azgo.QtreeCreateResponse, err error) {
	response, err = azgo.NewQtreeCreateRequest().
		SetQtree(name).
		SetVolume(volumeName).
		SetMode(unixPermissions).
		SetSecurityStyle(securityStyle).
		SetExportPolicy(exportPolicy).
		ExecuteUsing(d.zr)
	return
}

// QtreeRename renames a qtree
// equivalent to filer::> volume qtree rename
func (d Driver) QtreeRename(path, newPath string) (response azgo.QtreeRenameResponse, err error) {
	response, err = azgo.NewQtreeRenameRequest().
		SetQtree(path).
		SetNewQtreeName(newPath).
		ExecuteUsing(d.zr)
	return
}

// QtreeDestroy destroys a qtree in the background
// equivalent to filer::> volume qtree delete -foreground false
func (d Driver) QtreeDestroyAsync(path string, force bool) (response azgo.QtreeDeleteAsyncResponse, err error) {
	response, err = azgo.NewQtreeDeleteAsyncRequest().
		SetQtree(path).
		SetForce(force).
		ExecuteUsing(d.zr)
	return
}

// QtreeList returns the names of all Qtrees whose names match the supplied prefix
// equivalent to filer::> volume qtree show
func (d Driver) QtreeList(prefix, volumePrefix string) (response azgo.QtreeListIterResponse, err error) {

	// Limit the qtrees to those matching the Flexvol and Qtree name prefixes
	query := azgo.NewQtreeInfoType().SetVolume(volumePrefix + "*").SetQtree(prefix + "*")

	// Limit the returned data to only the Flexvol and Qtree names
	desiredAttributes := azgo.NewQtreeInfoType().SetVolume("").SetQtree("")

	response, err = azgo.NewQtreeListIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// QtreeCount returns the number of Qtrees in the specified Flexvol, not including the Flexvol itself
func (d Driver) QtreeCount(volume string) (int, error) {

	// Limit the qtrees to those in the specified Flexvol
	query := azgo.NewQtreeInfoType().SetVolume(volume)

	// Limit the returned data to only the Flexvol and Qtree names
	desiredAttributes := azgo.NewQtreeInfoType().SetVolume("").SetQtree("")

	response, err := azgo.NewQtreeListIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)

	if err = GetError(response, err); err != nil {
		return 0, err
	}

	// There will always be one qtree for the Flexvol, so decrement by 1
	switch response.Result.NumRecords() {
	case 0:
		fallthrough
	case 1:
		return 0, nil
	default:
		return response.Result.NumRecords() - 1, nil
	}
}

// QtreeExists returns true if the named Qtree exists (and is unique in the matching Flexvols)
func (d Driver) QtreeExists(name, volumePrefix string) (bool, string, error) {

	// Limit the qtrees to those matching the Flexvol and Qtree name prefixes
	query := azgo.NewQtreeInfoType().SetVolume(volumePrefix + "*").SetQtree(name)

	// Limit the returned data to only the Flexvol and Qtree names
	desiredAttributes := azgo.NewQtreeInfoType().SetVolume("").SetQtree("")

	response, err := azgo.NewQtreeListIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)

	// Ensure the API call succeeded
	if err = GetError(response, err); err != nil {
		return false, "", err
	}

	// Ensure qtree is unique
	if response.Result.NumRecords() != 1 {
		return false, "", nil
	}

	// Get containing Flexvol
	flexvol := response.Result.AttributesList()[0].Volume()

	return true, flexvol, nil
}

// QtreeGet returns all relevant details for a single qtree
// equivalent to filer::> volume qtree show
func (d Driver) QtreeGet(name, volumePrefix string) (azgo.QtreeInfoType, error) {

	// Limit the qtrees to those matching the Flexvol and Qtree name prefixes
	query := azgo.NewQtreeInfoType().SetVolume(volumePrefix + "*").SetQtree(name)

	response, err := azgo.NewQtreeListIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		ExecuteUsing(d.zr)

	if err != nil {
		return azgo.QtreeInfoType{}, err
	} else if response.Result.NumRecords() == 0 {
		return azgo.QtreeInfoType{}, fmt.Errorf("Qtree %s not found.", name)
	} else if response.Result.NumRecords() > 1 {
		return azgo.QtreeInfoType{}, fmt.Errorf("More than one qtree %s found.", name)
	}

	return response.Result.AttributesList()[0], nil
}

// QtreeGetAll returns all relevant details for all qtrees whose Flexvol names match the supplied prefix
// equivalent to filer::> volume qtree show
func (d Driver) QtreeGetAll(volumePrefix string) (response azgo.QtreeListIterResponse, err error) {

	// Limit the qtrees to those matching the Flexvol name prefix
	query := azgo.NewQtreeInfoType().SetVolume(volumePrefix + "*")

	// Limit the returned data to only the data relevant to containers
	desiredAttributes := azgo.NewQtreeInfoType().
		SetVolume("").
		SetQtree("").
		SetSecurityStyle("").
		SetMode("").
		SetExportPolicy("")

	response, err = azgo.NewQtreeListIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// QuotaOn enables quotas on a Flexvol
// equivalent to filer::> volume quota on
func (d Driver) QuotaOn(volume string) (response azgo.QuotaOnResponse, err error) {
	response, err = azgo.NewQuotaOnRequest().
		SetVolume(volume).
		ExecuteUsing(d.zr)
	return
}

// QuotaOff disables quotas on a Flexvol
// equivalent to filer::> volume quota off
func (d Driver) QuotaOff(volume string) (response azgo.QuotaOffResponse, err error) {
	response, err = azgo.NewQuotaOffRequest().
		SetVolume(volume).
		ExecuteUsing(d.zr)
	return
}

// QuotaResize resizes quotas on a Flexvol
// equivalent to filer::> volume quota resize
func (d Driver) QuotaResize(volume string) (response azgo.QuotaResizeResponse, err error) {
	response, err = azgo.NewQuotaResizeRequest().
		SetVolume(volume).
		ExecuteUsing(d.zr)
	return
}

// QuotaStatus returns the quota status for a Flexvol
// equivalent to filer::> volume quota show
func (d Driver) QuotaStatus(volume string) (response azgo.QuotaStatusResponse, err error) {
	response, err = azgo.NewQuotaStatusRequest().
		SetVolume(volume).
		ExecuteUsing(d.zr)
	return
}

// QuotaSetEntry creates a new quota rule with an optional hard disk limit
// equivalent to filer::> volume quota policy rule create
func (d Driver) QuotaSetEntry(qtreeName, volumeName, quotaTarget, quotaType, diskLimit string) (response azgo.QuotaSetEntryResponse, err error) {

	request := azgo.NewQuotaSetEntryRequest().
		SetQtree(qtreeName).
		SetVolume(volumeName).
		SetQuotaTarget(quotaTarget).
		SetQuotaType(quotaType)

	// To create a default quota rule, pass an empty disk limit
	if diskLimit != "" {
		request.SetDiskLimit(diskLimit)
	}

	return request.ExecuteUsing(d.zr)
}

// QuotaEntryGet returns the disk limit for a single qtree
// equivalent to filer::> volume quota policy rule show
func (d Driver) QuotaEntryGet(target string) (azgo.QuotaEntryType, error) {

	query := azgo.NewQuotaEntryType().SetQuotaType("tree").SetQuotaTarget(target)

	// Limit the returned data to only the disk limit
	desiredAttributes := azgo.NewQuotaEntryType().SetDiskLimit("").SetQuotaTarget("")

	response, err := azgo.NewQuotaListEntriesIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)

	if err != nil {
		return azgo.QuotaEntryType{}, err
	} else if response.Result.NumRecords() == 0 {
		return azgo.QuotaEntryType{}, fmt.Errorf("Tree quota for %s not found.", target)
	} else if response.Result.NumRecords() > 1 {
		return azgo.QuotaEntryType{}, fmt.Errorf("More than one tree quota for %s found.", target)
	}

	return response.Result.AttributesList()[0], nil
}

// QuotaEntryList returns the disk limit quotas for a Flexvol
// equivalent to filer::> volume quota policy rule show
func (d Driver) QuotaEntryList(volume string) (response azgo.QuotaListEntriesIterResponse, err error) {

	query := azgo.NewQuotaEntryType().SetVolume(volume).SetQuotaType("tree")

	// Limit the returned data to only the disk limit
	desiredAttributes := azgo.NewQuotaEntryType().SetDiskLimit("").SetQuotaTarget("")

	response, err = azgo.NewQuotaListEntriesIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// QTREE operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// EXPORT POLICY operations BEGIN

// ExportPolicyCreate creates an export policy
// equivalent to filer::> vserver export-policy create
func (d Driver) ExportPolicyCreate(policy string) (response azgo.ExportPolicyCreateResponse, err error) {
	response, err = azgo.NewExportPolicyCreateRequest().
		SetPolicyName(azgo.ExportPolicyNameType(policy)).
		ExecuteUsing(d.zr)
	return
}

// ExportRuleCreate creates a rule in an export policy
// equivalent to filer::> vserver export-policy rule create
func (d Driver) ExportRuleCreate(
	policy, clientMatch string,
	protocols, roSecFlavors, rwSecFlavors, suSecFlavors []string,
) (response azgo.ExportRuleCreateResponse, err error) {

	protocolTypes := []azgo.AccessProtocolType{}
	for _, p := range protocols {
		protocolTypes = append(protocolTypes, azgo.AccessProtocolType(p))
	}

	roSecFlavorTypes := []azgo.SecurityFlavorType{}
	for _, f := range roSecFlavors {
		roSecFlavorTypes = append(roSecFlavorTypes, azgo.SecurityFlavorType(f))
	}

	rwSecFlavorTypes := []azgo.SecurityFlavorType{}
	for _, f := range rwSecFlavors {
		rwSecFlavorTypes = append(rwSecFlavorTypes, azgo.SecurityFlavorType(f))
	}

	suSecFlavorTypes := []azgo.SecurityFlavorType{}
	for _, f := range suSecFlavors {
		suSecFlavorTypes = append(suSecFlavorTypes, azgo.SecurityFlavorType(f))
	}

	response, err = azgo.NewExportRuleCreateRequest().
		SetPolicyName(azgo.ExportPolicyNameType(policy)).
		SetClientMatch(clientMatch).
		SetProtocol(protocolTypes).
		SetRoRule(roSecFlavorTypes).
		SetRwRule(rwSecFlavorTypes).
		SetSuperUserSecurity(suSecFlavorTypes).
		ExecuteUsing(d.zr)
	return
}

// ExportRuleGetIterRequest returns the export rules in an export policy
// equivalent to filer::> vserver export-policy rule show
func (d Driver) ExportRuleGetIterRequest(policy string) (response azgo.ExportRuleGetIterResponse, err error) {

	// Limit the qtrees to those matching the Flexvol and Qtree name prefixes
	query := azgo.NewExportRuleInfoType().SetPolicyName(azgo.ExportPolicyNameType(policy))

	response, err = azgo.NewExportRuleGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		ExecuteUsing(d.zr)
	return
}

// EXPORT POLICY operations END
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
		SetMaxRecords(defaultZapiRecords).
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
		SetMaxRecords(defaultZapiRecords).
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
		SetMaxRecords(defaultZapiRecords).
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

	response, err := azgo.NewVserverGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		SetQuery(*query).
		ExecuteUsing(d.zr)

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
		SetMaxRecords(defaultZapiRecords).
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
		SetMaxRecords(defaultZapiRecords).
		ExecuteUsing(zr)
	return
}

// AGGREGATE operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// SNAPMIRROR operations BEGIN

// SnapmirrorGetLoadSharingMirrors gets load-sharing SnapMirror relationships for a volume
// equivalent to filer::> snapmirror -type LS -source-volume
func (d Driver) SnapmirrorGetLoadSharingMirrors(volume string) (response azgo.SnapmirrorGetIterResponse, err error) {

	// Limit the mirrors to load-sharing mirrors matching the source FlexVol
	query := azgo.NewSnapmirrorInfoType().SetSourceVolume(volume).SetRelationshipType("load_sharing")

	// Limit the returned data to only the source location
	desiredAttributes := azgo.NewSnapmirrorInfoType().SetSourceLocation("").SetRelationshipStatus("")

	response, err = azgo.NewSnapmirrorGetIterRequest().
		SetQuery(*query).
		SetDesiredAttributes(*desiredAttributes).
		ExecuteUsing(d.zr)
	return
}

// SnapmirrorUpdateLoadSharingMirrors updates the destination volumes of a set of load-sharing mirrors
// equivalent to filer::> snapmirror update-ls-set -source-path
func (d Driver) SnapmirrorUpdateLoadSharingMirrors(
	sourceLocation string,
) (response azgo.SnapmirrorUpdateLsSetResponse, err error) {

	response, err = azgo.NewSnapmirrorUpdateLsSetRequest().
		SetSourceLocation(sourceLocation).
		ExecuteUsing(d.zr)
	return
}

// SNAPMIRROR operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// MISC operations BEGIN

// NetInterfaceGet returns the list of network interfaces with associated metadata
// equivalent to filer::> net interface list
func (d Driver) NetInterfaceGet() (response azgo.NetInterfaceGetIterResponse, err error) {
	response, err = azgo.NewNetInterfaceGetIterRequest().
		SetMaxRecords(defaultZapiRecords).
		ExecuteUsing(d.zr)
	return
}

func (d Driver) NetInterfaceGetDataLIFs(protocol string) ([]string, error) {
	lifResponse, err := d.NetInterfaceGet()
	if err = GetError(lifResponse, err); err != nil {
		return nil, fmt.Errorf("Error checking network interfaces. %v", err)
	}

	dataLIFs := make([]string, 0)
	for _, attrs := range lifResponse.Result.AttributesList() {
		for _, proto := range attrs.DataProtocols() {
			if proto == azgo.DataProtocolType(protocol) {
				dataLIFs = append(dataLIFs, string(attrs.Address()))
			}
		}
	}

	log.WithField("dataLIFs", dataLIFs).Debug("Data LIFs")
	return dataLIFs, nil
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
		if err = GetError(result, err); err != nil {
			return "", fmt.Errorf("Could not read ONTAPI version. %v", err)
		}

		major := result.Result.MajorVersion()
		minor := result.Result.MinorVersion()
		d.zr.OntapiVersion = fmt.Sprintf("%d.%d", major, minor)
	}

	return d.zr.OntapiVersion, nil
}

func (d Driver) ListNodeSerialNumbers() ([]string, error) {

	serialNumbers := make([]string, 0, 0)
	zr := d.GetNontunneledZapiRunner()

	// Limit the returned data to only the serial numbers
	desiredAttributes := azgo.NewNodeDetailsInfoType().SetNodeSerialNumber("")

	response, err := azgo.NewSystemNodeGetIterRequest().
		SetDesiredAttributes(*desiredAttributes).
		SetMaxRecords(defaultZapiRecords).
		ExecuteUsing(zr)

	if err = GetError(response, err); err != nil {
		return serialNumbers, err
	}

	if response.Result.NumRecords() == 0 {
		return serialNumbers, errors.New("Could not get node info.")
	}

	// Get the serial numbers
	for _, node := range response.Result.AttributesList() {
		serialNumber := node.NodeSerialNumber()
		if serialNumber != "" {
			serialNumbers = append(serialNumbers, serialNumber)
		}
	}

	if len(serialNumbers) == 0 {
		return serialNumbers, errors.New("Could not get node serial numbers.")
	}

	log.WithFields(log.Fields{
		"Count":         len(serialNumbers),
		"SerialNumbers": serialNumbers,
	}).Debug("Read serial numbers.")

	return serialNumbers, nil
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
