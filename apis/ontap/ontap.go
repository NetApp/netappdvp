// Copyright 2016 NetApp, Inc. All Rights Reserved.

package ontap

import (
	"sync"

	"github.com/netapp/netappdvp/azgo"
)

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
// equivalent to filer::> volume create -vserver iscsi_vs -volume v -aggregate aggr1 -size 1g -state online -type RW -policy default -unix-permissions ---rwxr-xr-x -space-guarantee none -snapshot-policy none
func (d Driver) VolumeCreate(name, aggregateName, size, spaceReserve, snapshotPolicy, unixPermissions, exportPolicy string) (response azgo.VolumeCreateResponse, err error) {
	response, err = azgo.NewVolumeCreateRequest().
		SetVolume(name).
		SetContainingAggrName(aggregateName).
		SetSize(size).
		SetSpaceReserve(spaceReserve).
		SetSnapshotPolicy(snapshotPolicy).
		SetUnixPermissions(unixPermissions).
		SetExportPolicy(exportPolicy).
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
		SetQuery(*query).
		ExecuteUsing(d.zr)
	return
}

// SNAPSHOT operations END
/////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// MISC operations BEGIN

// NetInterfaceGet returns the list of network interfaces with associated metadata
// equivalent to filer::> net interface list
func (d Driver) NetInterfaceGet() (response azgo.NetInterfaceGetIterResponse, err error) {
	response, err = azgo.NewNetInterfaceGetIterRequest().ExecuteUsing(d.zr)
	return
}

// SystemGetVersion returns the system version
// equivalent to filer::> version
func (d Driver) SystemGetVersion() (response azgo.SystemGetVersionResponse, err error) {
	response, err = azgo.NewSystemGetVersionRequest().ExecuteUsing(d.zr)
	return
}

// VserverGetIterRequest returns the vservers on the system
// equivalent to filer::> vserver show
func (d Driver) VserverGetIterRequest() (response azgo.VserverGetIterResponse, err error) {
	response, err = azgo.NewVserverGetIterRequest().ExecuteUsing(d.zr)
	return
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
