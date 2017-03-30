// Copyright 2016 NetApp, Inc. All Rights Reserved.

package ontap

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/azgo"
)

// TODO externalize for testing against different configurations
func newConfig() (c *DriverConfig) {
	c = &DriverConfig{}
	c.ManagementLIF = "10.63.171.241"
	c.SVM = "CXE"
	c.Username = "admin"
	c.Password = "Netapp123"
	return
}

var initiatorGroupName = "api_test_igroup"
var lunPath = "/vol/v/lun0"
var volName = "api_test_vol"

func newZapiRunner(config DriverConfig) *azgo.ZapiRunner {
	zr := &azgo.ZapiRunner{
		ManagementLIF: config.ManagementLIF,
		SVM:           config.SVM,
		Username:      config.Username,
		Password:      config.Password,
		Secure:        true,
	}
	return zr
}

func TestSystemGetVersion(t *testing.T) {
	log.Debug("Running TestSystemGetVersion...")

	zr := newZapiRunner(*newConfig())

	r0, err0 := azgo.NewSystemGetVersionRequest().ExecuteUsing(zr)
	log.Debugf("r0.Result: %s\n", r0.Result)
	if err0 != nil {
		t.Error("Could not validate credentials")
	}

	// TODO add some sort of system version validation
	systemVersion := r0.Result
	if systemVersion.VersionPtr == nil {
		t.Error("Could not get system version")
	}
}

func TestSystemApiGetVersion(t *testing.T) {
	log.Debug("Running TestSystemApiGetVersion...")

	zr := newZapiRunner(*newConfig())

	r0, err0 := azgo.NewSystemGetOntapiVersionRequest().ExecuteUsing(zr)
	log.Debugf("r0.Result: %s\n", r0.Result)
	if err0 != nil {
		t.Error("Could not validate credentials")
	}

	systemOntapiVersion := r0.Result
	if systemOntapiVersion.MajorVersionPtr == nil || systemOntapiVersion.MinorVersionPtr == nil {
		t.Error("Could not get ontapi version")
	}
}

func TestErrorCode(t *testing.T) {
	log.Debug("Running TestErrorCode...")

	if azgo.EVDISK_ERROR_INITGROUP_MAPS_EXIST != "9029" {
		t.Error("Could not validate constant, something's wrong...")
	}
}

func TestIgroup(t *testing.T) {
	log.Debug("Running TestIgroup...")

	c := newConfig()
	d := NewDriver(*c)

	// check wrong os type fails
	response, err := d.IgroupCreate(initiatorGroupName, "iscsi", "leenux")
	if response.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive invalid input error for incorrect ostype 'leenux'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check wrong group type fails
	response, err = d.IgroupCreate(initiatorGroupName, "eyescsi", "linux")
	if response.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive invalid input error for incorrect group type 'eyescsi'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check create passes
	response, err = d.IgroupCreate(initiatorGroupName, "iscsi", "linux")
	if response.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to create igroup")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check double create fails
	response, err = d.IgroupCreate(initiatorGroupName, "iscsi", "linux")
	if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_EXISTS {
		t.Error("Expected to fail to create existing igroup that we should made")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check igroup add fails
	response2, err2 := d.IgroupAdd(initiatorGroupName, "bad")
	if response2.Result.ResultErrnoAttr != azgo.EAPIERROR {
		t.Error("Expected to fail to add an invalid initiator name")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check igroup add passes
	initiator := "iqn.1993-08.org.debian:01:9031309bbebd"
	response2, err2 = d.IgroupAdd(initiatorGroupName, initiator)
	if response2.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to add valid initiator name to igroup")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check igroup double add fails
	response2, err2 = d.IgroupAdd(initiatorGroupName, initiator)
	if response2.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_HAS_NODE {
		t.Error("Expected to fail to add an initiator twice to the same initiator group")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check igroup remove passes
	response3, err3 := d.IgroupRemove(initiatorGroupName, initiator, true)
	if response3.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to remove initiator from igroup")
	}
	if err3 != nil {
		t.Error("Unexpected error found")
	}

	// check igroup initiator remove fails if not in group
	response3, err3 = d.IgroupRemove(initiatorGroupName, initiator, true)
	if response3.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_NODE_NOT_IN_INITGROUP {
		t.Error("Expected to fail to remove a non existant initiator from an initiator group")
	}
	if err3 != nil {
		t.Error("Unexpected error found")
	}

	// check destroy passes
	response4, err4 := d.IgroupDestroy(initiatorGroupName)
	if response4.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to destroy igroup")
	}
	if err4 != nil {
		t.Error("Unexpected error found")
	}

	// check double destroy fails
	response4, err4 = d.IgroupDestroy(initiatorGroupName)
	if response4.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_NO_SUCH_INITGROUP {
		t.Error("Expected to fail to delete nonexisting igroup")
	}
	if err4 != nil {
		t.Error("Unexpected error found")
	}
}

func TestLun(t *testing.T) {
	log.Debug("Running TestLun...")

	c := newConfig()
	d := NewDriver(*c)

	// check wrong os type fails
	response, err := d.LunCreate(lunPath, 1, "leenux", false)
	if response.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive invalid input error for incorrect ostype 'leenux'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check missing volume fails
	response, err = d.LunCreate("/vol/baddddVolume/lun0", 1, "linux", false)
	if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_NO_SUCH_VOLUME {
		t.Error("Expected to receive invalid input error for nonexisting volume 'baddddVolume'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check invalid size fails
	response, err = d.LunCreate(lunPath, 1, "linux", false)
	if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_SIZE_TOO_SMALL {
		t.Error("Expected to receive invalid size error for lun of size '1'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check valid settings create lun of 1gb
	response, err = d.LunCreate(lunPath, 1048576*1024, "linux", false)
	if response.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to create lun")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check double create fails (which validates lun was created as a side effect)
	response, err = d.LunCreate(lunPath, 1048576*1024, "linux", false)
	if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_VDISK_EXISTS {
		t.Error("Expected to receive disk exists error for already created lun")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check lun offline passes
	response2, err2 := d.LunOffline(lunPath)
	if response2.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to offline lun")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check lun online passes
	response3, err3 := d.LunOnline(lunPath)
	if response3.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to offline lun")
	}
	if err3 != nil {
		t.Error("Unexpected error found")
	}

	// check lun online 2x fails
	response3, err3 = d.LunOnline(lunPath)
	if response3.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_VDISK_NOT_DISABLED {
		t.Error("Expected to receive disk not disabled error for already onlined lun")
	}
	if err3 != nil {
		t.Error("Unexpected error found")
	}

	// check lun offline passes
	response2, err2 = d.LunOffline(lunPath)
	if response2.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to offline lun")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check lun offline 2x fails
	response2, err2 = d.LunOffline(lunPath)
	if response2.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_VDISK_NOT_ENABLED {
		t.Error("Expected to receive disk not enabled error for already offlined lun")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check lun destroy passes
	response4, err4 := d.LunDestroy(lunPath)
	if response4.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to delete lun")
	}
	if err4 != nil {
		t.Error("Unexpected error found")
	}

	// check lun 2x destroy fails (lun already deleted)
	response4, err4 = d.LunDestroy(lunPath)
	if response4.Result.ResultErrnoAttr != azgo.EOBJECTNOTFOUND {
		t.Error("Expected to receive object missing error for already deleted lun")
	}
	if err4 != nil {
		t.Error("Unexpected error found")
	}
}

func TestLunMapping(t *testing.T) {
	log.Debug("Running TestLunMapping...")

	c := newConfig()
	d := NewDriver(*c)

	lunSize := 1048576 * 1024

	d.IgroupCreate(initiatorGroupName, "iscsi", "linux")
	d.LunCreate(lunPath, lunSize, "linux", false)
	d.LunCreate(lunPath+"b", lunSize, "linux", false)

	// check lun map passes
	response, err := d.LunMap(initiatorGroupName, lunPath, 0)
	if response.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to map lun to id '0'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check lun map fails if already mapped
	response, err = d.LunMap(initiatorGroupName, lunPath, 0)
	if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_HAS_VDISK {
		t.Error("Expected to error because LUN already mapped")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check lun map fails if in use
	response, err = d.LunMap(initiatorGroupName, lunPath+"b", 0)
	if response.Result.ResultErrnoAttr != azgo.EVDISK_ERROR_INITGROUP_HAS_LUN {
		t.Error("Expected to error because LUN id in use")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check if lun is NOT mapped behavior
	response2, err2 := d.LunMapListInfo(lunPath + "b")
	if response2.Result.ResultStatusAttr != "passed" &&
		response2.Result.InitiatorGroups() != nil {
		t.Error("Expected to pass and see no initiator groups")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check if lun IS mapped behavior
	response2, err2 = d.LunMapListInfo(lunPath)
	if response2.Result.ResultStatusAttr != "passed" ||
		response2.Result.InitiatorGroups() == nil {
		t.Error("Expected to pass and see initiator groups")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}
	if response2.Result.InitiatorGroups()[0].LunId() != 0 {
		t.Error("Expected to be mapped to lun id '0'")
	}

	d.LunDestroy(lunPath + "b")
	d.LunOffline(lunPath)
	d.LunDestroy(lunPath)
	d.IgroupDestroy(initiatorGroupName)
}

func TestVolumeAndSnapshot(t *testing.T) {
	log.Debug("Running TestVolumeAndSnapshot...")

	c := newConfig()
	d := NewDriver(*c)
	//tc := newTestConfig()

	//aggr := tc.Aggregate
	aggr := "VICE08_aggr1"
	unixPerms := "---rwxr-xr-x"
	exportPolicy := "default"

	// check bad volume name fails
	response, err := d.VolumeCreate("bad/bad", aggr, "1g", "none", "none", unixPerms, exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EAPIERROR {
		t.Error("Expected to receive invalid api error for name 'bad/bad'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check bad unix permissions fails
	response, err = d.VolumeCreate(volName, aggr, "1g", "none", "none", "bad", exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive invalid input error for invalid unix permissions 'bad'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check missing aggregate fails
	response, err = d.VolumeCreate(volName, "missingAggrBad", "1g", "none", "none", unixPerms, exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EAGGRDOESNOTEXIST {
		t.Error("Expected to receive aggr doesn't exist error for invalid aggregrate 'missingAggrBad'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check bad size fails
	response, err = d.VolumeCreate(volName, aggr, "badSize", "none", "none", unixPerms, exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive error for invalid size 'badSize'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check bad space reserve fails
	response, err = d.VolumeCreate(volName, aggr, "1g", "badSpaceReserve", "none", unixPerms, exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive error for invalid space reserve 'badSpaceReserve'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check bad snapshotPolicy fails
	response, err = d.VolumeCreate(volName, aggr, "1g", "none", "badSnapshotPolicy", unixPerms, exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EAPIERROR {
		t.Error("Expected to receive error for invalid snapshot policy 'badSnapshotPolicy'")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check create passes
	response, err = d.VolumeCreate(volName, aggr, "1g", "none", "none", unixPerms, exportPolicy, "unix")
	if response.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to create volume")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check double create fails
	response, err = d.VolumeCreate(volName, aggr, "1g", "none", "none", unixPerms, exportPolicy, "unix")
	if response.Result.ResultErrnoAttr != azgo.EONTAPI_EEXIST {
		t.Error("Expected to receive error for creating an already existing volume")
	}
	if err != nil {
		t.Error("Unexpected error found")
	}

	// check bad volume mount fails
	response2, err2 := d.VolumeMount(volName, "-"+volName)
	if response2.Result.ResultErrnoAttr != azgo.EINVALIDINPUTERROR {
		t.Error("Expected to receive error for mounting to a bad junction path '-' in name")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check volume mount passes
	response2, err2 = d.VolumeMount(volName, "/"+volName)
	if response2.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to mount volume to junction path")
	}
	if err2 != nil {
		t.Error("Unexpected error found")
	}

	// check disabling .snapshot directory access passes
	response3, err3 := d.VolumeDisableSnapshotDirectoryAccess(volName)
	if response3.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to be able to disable '.snapshot' directory access")
	}
	if err3 != nil {
		t.Error("Unexpected error found")
	}

	// check size lookup for missing volume fails
	response3b, err3b := d.VolumeSize("badVol")
	if response3b.Result.ResultErrnoAttr != azgo.EOBJECTNOTFOUND {
		t.Error("Expected to error looking up size for non existant volume")
	}
	if err3b != nil {
		t.Error("Unexpected error found")
	}

	// check size lookup for existing volume passes
	response3b, err3b = d.VolumeSize(volName)
	if response3b.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to be able to lookup volume size")
	}
	if err3b != nil {
		t.Error("Unexpected error found")
	}

	// check getting snapshots with missing volume fails
	response3c, err3c := d.SnapshotGetByVolume("badVol")
	if response3c.Result.ResultErrnoAttr != azgo.EOBJECTNOTFOUND {
		t.Error("Expected EOBJECTNOTFOUND getting snapshots for non-existent volume")
	}
	if err3c != nil {
		t.Error("Unexpected error found")
	}

	// check that a snapshot can be created on the existing volume
	response3d, err3d := d.SnapshotCreate("mysnap", volName)
	if response3d.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to create a snapshot")
	}
	if err3d != nil {
		t.Error("Unexpected error found")
	}

	// check getting snapshots for existing volume succeeds
	response3e, err3e := d.SnapshotGetByVolume(volName)
	if response3e.Result.ResultErrnoAttr != "passed" {
		t.Error("Expected to be able to get a list of snapshots")
	}
	if response3e.Result.AttributesList()[0].Name() != "mysnap" {
		t.Error("Expected the only snapshot to be the one that was just created")
	}
	if err3e != nil {
		t.Error("Unexpected error found")
	}

	// check volume offline fails (if still mounted)
	response4, err4 := d.VolumeOffline(volName)
	if response4.Result.ResultErrnoAttr != azgo.EONTAPI_EVOLOPNOTSUPP {
		t.Error("Expected to receive error for unmounting a mounted volume")
	}
	if err4 != nil {
		t.Error("Unexpected error found")
	}

	// check volume unmount succeeds
	response5, err5 := d.VolumeUnmount(volName, true)
	if response5.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to unmount volume")
	}
	if err5 != nil {
		t.Error("Unexpected error found")
	}

	// check volume offline passes (if no longer mounted)
	response4, err4 = d.VolumeOffline(volName)
	if response4.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to offline volume")
	}
	if err4 != nil {
		t.Error("Unexpected error found")
	}

	// check destroy passes
	response6, err6 := d.VolumeDestroy(volName, true)
	if response6.Result.ResultStatusAttr != "passed" {
		t.Error("Expected to destroy volume")
	}
	if err6 != nil {
		t.Error("Unexpected error found")
	}

	// check double destroy fails
	response6, err6 = d.VolumeDestroy(volName, true)
	if response6.Result.ResultErrnoAttr != azgo.EVOLUMEDOESNOTEXIST {
		t.Error("Expected to receive error for deleting a non existant volume")
	}
	if err6 != nil {
		t.Error("Unexpected error found")
	}
}
