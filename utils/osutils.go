// Copyright 2016 NetApp, Inc. All Rights Reserved.

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

const ISCSI_ERR_NO_OBJS_FOUND int = 21

// DFInfo data structure for wrapping the parsed output from the 'df' command
type DFInfo struct {
	Target string
	Source string
}

// GetDFOutput returns parsed DF output
func GetDFOutput() ([]DFInfo, error) {
	log.Debug("Begin osutils.GetDFOutput")
	var result []DFInfo
	out, err := exec.Command("df", "--output=target,source").Output()
	if err != nil {
		// df returns an error if there's a stale file handle that we can
		// safely ignore. There may be other reasons. Consider it a warning if
		// it printed anything to stdout.
		if len(out) == 0 {
			log.Errorf("Error encountered gathering df output: %v.", err)
			return nil, err
		}
	}
	//log.Debugf("out==%v", string(out))
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		//log.Debugf("l==%v", l)
		a := strings.Fields(l)
		if len(a) > 1 {
			result = append(result, DFInfo{
				Target: a[0],
				Source: a[1],
			})
		}
	}
	if len(result) > 1 {
		return result[1:], nil
	}
	return result, nil
}

func Stat(fileName string) (string, error) {
	statCmd := fmt.Sprintf("stat %v", fileName)
	log.Debugf("running 'sh -c %v'", statCmd)
	out, err := exec.Command("sh", "-c", statCmd).CombinedOutput()
	log.Debugf("out==%v", string(out))
	if err == nil {
		return string(out), err
	} else {
		return string(out), err
	}
}

// GetInitiatorIqns returns parsed contents of /etc/iscsi/initiatorname.iscsi
func GetInitiatorIqns() ([]string, error) {
	log.Debug("Begin osutils.GetInitiatorIqns")
	var iqns []string
	out, err := exec.Command("cat", "/etc/iscsi/initiatorname.iscsi").CombinedOutput()
	if err != nil {
		log.Errorf("Error gathering initiator names: %v. %v", err, string(out))
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if strings.Contains(l, "InitiatorName=") {
			iqns = append(iqns, strings.Split(l, "=")[1])
		}
	}
	return iqns, nil
}

// WaitForPathToExist retries every second, up to numTries times, with increasing backoff, for the specified fileName to show up
func WaitForPathToExist(fileName string, numTries int) bool {
	log.Debugf("Begin osutils.waitForPathToExist fileName: %v", fileName)
	for i := 0; i < numTries; i++ {
		_, err := Stat(fileName)
		if err == nil {
			log.Debugf("path found for fileName on attempt %v of %v: %v", i, numTries, fileName)
			return true
		}
		time.Sleep(time.Second * time.Duration(2+i))
	}
	log.Warnf("osutils.waitForPathToExist giving up looking for fileName: %v", fileName)
	return false
}

// ScsiDeviceInfo contains information about SCSI devices
type ScsiDeviceInfo struct {
	Host            string
	Channel         string
	Target          string
	LUN             string
	Device          string
	MultipathDevice string
	Filesystem      string
	IQN             string
}

// LsscsiCmd executes and parses the output from the 'lsscsi' command
func LsscsiCmd(args []string) ([]ScsiDeviceInfo, error) {
	/*
		# lsscsi
		[0:0:0:0]    disk    ATA      VBOX HARDDISK    1.0   /dev/sda
		[5:0:0:0]    disk    NETAPP   LUN C-Mode       8200  /dev/sdb
		[6:0:0:0]    disk    NETAPP   LUN C-Mode       8200  /dev/sdc
		[7:0:0:0]    disk    NETAPP   LUN C-Mode       8200  /dev/sdd
		[8:0:0:0]    disk    NETAPP   LUN C-Mode       8200  /dev/sde

		# lsscsi -t
		[0:0:0:0]    disk    sata:                           /dev/sda
		[5:0:0:0]    disk    iqn.1992-08.com.netapp:sn.afbb1784f77411e582f8080027e22798:vs.3,t,0x404  /dev/sdb
		[6:0:0:0]    disk    iqn.1992-08.com.netapp:sn.afbb1784f77411e582f8080027e22798:vs.3,t,0x405  /dev/sdc
		[7:0:0:0]    disk    iqn.1992-08.com.netapp:sn.d724e00bfa0311e582f8080027e22798:vs.4,t,0x407  /dev/sdd
		[8:0:0:0]    disk    iqn.1992-08.com.netapp:sn.d724e00bfa0311e582f8080027e22798:vs.4,t,0x408  /dev/sde
	*/
	hasArgs := args != nil && len(args) > 0

	log.Debugf("Begin osutils.LsscsiCmd: %v", args)
	out, err := exec.Command("lsscsi", args...).CombinedOutput()
	if err != nil {
		log.Errorf("Error listing iSCSI devices: %v. %v", err, string(out))
		return nil, err
	}

	var info []ScsiDeviceInfo

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	log.Debugf("Found output lines: %v", lines)
	for _, line := range lines {
		log.Debugf("processing line: %v", line)

		d := strings.Fields(line)
		if d == nil || len(d) < 1 {
			log.Debugf("could not parse output, skipping line: %v", line)
			continue
		}

		log.Debugf("Found d: %v", d)
		s := d[0]
		s = s[1 : len(s)-1]
		scsiBusInfo := strings.Split(s, ":")

		scsiHost := scsiBusInfo[0]
		scsiChannel := scsiBusInfo[1]
		scsiTarget := scsiBusInfo[2]
		scsiLun := scsiBusInfo[3]

		// the device is the last field in the output
		devFile := d[len(d)-1]
		log.Debugf("devFile: %v", devFile)
		if !strings.HasPrefix(devFile, "/dev") {
			log.Debugf("could not find device in output, skipping line: %v", line)
			continue
		}

		// we only have iqn info if '-t' specified
		iqn := ""
		if hasArgs {
			iqn = d[2]
		}
		log.Debugf("iqn: %v", iqn)

		// check to see if there's a multipath device
		multipathDevFile := ""
		if !MultipathDetected() {
			log.Debug("Skipping multipath check, /sbin/multipath doesn't exist")
		} else {

			lsblkCmd := fmt.Sprintf("lsblk %v -n -o name,type -r | grep mpath | cut -f1 -d\\ ", devFile)
			log.Debugf("running 'sh -c %v'", lsblkCmd)
			out2, err2 := exec.Command("sh", "-c", lsblkCmd).CombinedOutput()
			if err2 != nil {
				// this can be fine, for instance could be a floppy or cd-rom, later logic will error if we never find our device
				log.Debugf("Error running multipath check for device %v: %v. %v", devFile, err2, string(out2))
			} else {
				md := strings.Split(strings.TrimSpace(string(out2)), " ")
				if md != nil && len(md) > 0 && len(md[0]) > 0 {
					if strings.HasPrefix(md[0], "lsblk") || strings.HasSuffix(string(out2), "failed to get device path") {
						return nil, fmt.Errorf("Problem checking device path while running multipath check for device: %v: output: %v", devFile, string(out2))
					}
					log.Debug("Found md: ", md)

					multipathDevFileToCheck := "/dev/mapper/" + md[0]
					_, err3 := Stat(multipathDevFileToCheck)
					if err3 == nil {
						multipathDevFile = multipathDevFileToCheck
					}
				}
			}
		}

		fsType := ""
		if multipathDevFile != "" {
			fsType = GetFSType(multipathDevFile)
		} else {
			fsType = GetFSType(devFile)
		}

		log.WithFields(log.Fields{
			"scsiHost":         scsiHost,
			"scsiChannel":      scsiChannel,
			"scsiTarget":       scsiTarget,
			"scsiLun":          scsiLun,
			"multipathDevFile": multipathDevFile,
			"devFile":          devFile,
			"fsType":           fsType,
			"iqn":              iqn,
		}).Debug("Found")

		info = append(info, ScsiDeviceInfo{
			Host:            scsiHost,
			Channel:         scsiChannel,
			Target:          scsiTarget,
			LUN:             scsiLun,
			MultipathDevice: multipathDevFile,
			Device:          devFile,
			Filesystem:      fsType,
			IQN:             iqn,
		})
	}

	return info, nil
}

// GetDeviceInfoForLuns parses 'lsscsi' to find NetApp LUNs
func GetDeviceInfoForLuns() ([]ScsiDeviceInfo, error) {
	log.Debug("Begin osutils.getDeviceInfoForLuns: ")

	// first, list w/out iSCSI target info
	var info1 []ScsiDeviceInfo
	info1, err1 := LsscsiCmd(nil)
	if err1 != nil {
		return nil, err1
	}

	// now, list w/ iSCSI target info
	var info2 []ScsiDeviceInfo
	info2, err2 := LsscsiCmd([]string{"-t"})
	if err2 != nil {
		return nil, err2
	}

	// finally, glue the 2 outputs together
	for j, e1 := range info1 {
	innerLoop:
		for k, e2 := range info2 {
			if e1.MultipathDevice == "" || e2.MultipathDevice == "" {
				// no multipath device info, skipping
				if e1.Device == e2.Device {
					// no multipath device info, skipping multipath compare but we still need the IQN info
					log.Debugf("Matched, setting IQN to: %v", e2.IQN)
					info1[j].IQN = info2[k].IQN
				}
				continue
			}

			log.Debugf("Comparing d: %v and d: %v", e1.Device, e2.Device)
			log.Debugf("Comparing md: %v and md: %v", e1.MultipathDevice, e2.MultipathDevice)
			if (e1.Device == e2.Device) &&
				(e1.MultipathDevice == e2.MultipathDevice) {
				log.Debugf("Matched, setting IQN to: %v", e2.IQN)
				info1[j].IQN = info2[k].IQN
				break innerLoop
			}
		}
	}

	return info1, nil
}

// GetDeviceFileFromIscsiPath returns the /dev device for the supplied iscsiPath
func GetDeviceFileFromIscsiPath(iscsiPath string) (devFile string) {
	log.Debug("Begin osutils.GetDeviceFileFromIscsiPath: ", iscsiPath)
	out, err := exec.Command("ls", "-la", iscsiPath).CombinedOutput()
	if err != nil {
		log.Errorf("Error getting device file from iSCSI path %v: %v. %v", iscsiPath, err, string(out))
		return
	}
	d := strings.Split(string(out), "../../")
	log.Debugf("Found device: %v for iscsiPath: %v", d, iscsiPath)
	devFile = "/dev/" + d[1]
	log.Debugf("using device file: %v", devFile)
	devFile = strings.TrimSpace(devFile)
	return
}

// IscsiSupported returns true if iscsiadm is installed and in the PATH
func IscsiSupported() bool {
	out, err := exec.Command("iscsiadm", "-h").CombinedOutput()
	if err != nil {
		log.Debugf("iscsiadm tools not found on this host: %v. %v", err, string(out))
		return false
	}
	return true
}

// IscsiDiscoveryInfo contains information about discovered iSCSI targets
type IscsiDiscoveryInfo struct {
	Portal     string
	PortalIP   string
	TargetName string
}

// IscsiDiscovery uses the 'iscsiadm' command to perform discovery
func IscsiDiscovery(portal string) ([]IscsiDiscoveryInfo, error) {
	log.Debugf("Begin osutils.IscsiDiscovery (portal: %s)", portal)

	out, err := IscsiadmCmd([]string{"-m", "discovery", "-t", "sendtargets", "-p", portal})
	if err != nil {
		log.Errorf("Error encountered in sendtargets cmd: %v. %v", err, string(out))
		return nil, err
	}

	/*
			   iscsiadm -m discovery -t st -p 10.63.152.249:3260

		           10.63.152.249:3260,1 iqn.1992-08.com.netapp:2752.600a0980006074c20000000056b32c4d
		           10.63.152.250:3260,2 iqn.1992-08.com.netapp:2752.600a0980006074c20000000056b32c4d

		           a[0]==10.63.152.249:3260,1
		           a[1]==iqn.1992-08.com.netapp:2752.600a0980006074c20000000056b32c4d
	*/

	var discoveryInfo []IscsiDiscoveryInfo

	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		a := strings.Fields(l)
		if len(a) >= 2 {

			portalIP := strings.Split(a[0], ":")[0]

			discoveryInfo = append(discoveryInfo, IscsiDiscoveryInfo{
				Portal:     a[0],
				PortalIP:   portalIP,
				TargetName: a[1],
			})

			log.WithFields(log.Fields{
				"Portal":     a[0],
				"PortalIP":   portalIP,
				"TargetName": a[1],
			}).Debug("Adding iSCSI discovery info")
		}
	}
	return discoveryInfo, nil
}

// IscsiSessionInfo contains information about iSCSI sessions
type IscsiSessionInfo struct {
	SID        string
	Portal     string
	PortalIP   string
	TargetName string
}

// GetIscsiSessionInfo parses output from 'iscsiadm -m session' and returns the parsed output
func GetIscsiSessionInfo() ([]IscsiSessionInfo, error) {
	log.Debugf("Begin osutils.GetIscsiSessionInfo")

	out, err := IscsiadmCmd([]string{"-m", "session"})
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() == ISCSI_ERR_NO_OBJS_FOUND {
			log.Debug("No iSCSI session found.")
			return []IscsiSessionInfo{}, nil
		} else {
			log.Errorf("Problem checking iSCSI sessions. %v", err)
			return nil, err
		}
	}

	/*
	   # iscsiadm -m session
	   tcp: [3] 10.0.207.7:3260,1028 iqn.1992-08.com.netapp:sn.afbb1784f77411e582f8080027e22798:vs.3 (non-flash)
	   tcp: [4] 10.0.207.9:3260,1029 iqn.1992-08.com.netapp:sn.afbb1784f77411e582f8080027e22798:vs.3 (non-flash)

	   a[0]==tcp:
	   a[1]==[4]
	   a[2]==10.0.207.9:3260,1029
	   a[3]==iqn.1992-08.com.netapp:sn.afbb1784f77411e582f8080027e22798:vs.3
	   a[4]==(non-flash)
	*/

	var sessionInfo []IscsiSessionInfo

	//log.Debugf("out==%v", string(out))
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		//log.Debugf("l==%v", l)
		a := strings.Fields(l)
		if len(a) > 3 {
			sid := a[1]
			sid = sid[1 : len(sid)-1]

			//log.Debugf("a[2]==%v", strings.Split(a[2], ":"))
			portalIP := strings.Split(a[2], ":")[0]
			sessionInfo = append(sessionInfo, IscsiSessionInfo{
				SID:        sid,
				Portal:     a[2],
				PortalIP:   portalIP,
				TargetName: a[3],
			})

			log.WithFields(log.Fields{
				"SID":        sid,
				"Portal":     a[2],
				"PortalIP":   portalIP,
				"TargetName": a[3],
			}).Debug("Adding iSCSI session info")

		}
	}

	return sessionInfo, nil
}

// IscsiTargetInfo structure for usage with the iscsiadm command
type IscsiTargetInfo struct {
	IP        string
	Port      string
	Portal    string
	Iqn       string
	Lun       string
	Device    string
	Discovery string
}

// IscsiDisableDelete logout from the supplied target and remove the iSCSI device
func IscsiDisableDelete(tgt *IscsiTargetInfo) (err error) {
	log.Debugf("Begin osutils.IscsiDisableDelete: %v", tgt)
	out, err := exec.Command("iscsiadm", "-m", "node", "-T", tgt.Iqn, "--portal", tgt.IP, "-u").CombinedOutput()
	if err != nil {
		log.Debugf("Error during iSCSI logout: %v. %v", err, string(out))
	}
	_, err = exec.Command("iscsiadm", "-m", "node", "-o", "delete", "-T", tgt.Iqn).CombinedOutput()
	return
}

// IscsiSessionExists checks to see if a session exists to the sepecified portal
func IscsiSessionExists(portal string) (bool, error) {
	log.Debugf("Begin osutils.IscsiSessionExists")

	sessionInfo, err := GetIscsiSessionInfo()
	if err != nil {
		log.Errorf("Problem checking iSCSI sessions error: %v", err)
		return false, err
	}

	for _, e := range sessionInfo {
		if e.PortalIP == portal {
			return true, nil
		}
	}

	return false, nil
}

// IscsiRescan uses the 'rescan-scsi-bus' command to perform rescanning of the SCSI bus
func IscsiRescan() (err error) {
	log.Debugf("Begin osutils.IscsiRescan")
	defer UdevSettle()

	// look for version of rescan-scsi-bus in known locations
	var rescanCommands []string = []string{"/sbin/rescan-scsi-bus", "/sbin/rescan-scsi-bus.sh", "/bin/rescan-scsi-bus.sh", "/usr/bin/rescan-scsi-bus.sh"}
	for _, rescanCommand := range rescanCommands {
		_, err = os.Lstat(rescanCommand)
		// The command exists in this location
		if !os.IsNotExist(err) {
			out, rescanErr := exec.Command(rescanCommand, "-a", "-r").CombinedOutput()
			// We encountered an error condition
			if rescanErr != nil {
				log.Errorf("Could not rescan SCSI bus: %v. %v", rescanErr, string(out))
				return rescanErr
			} else {
				// The command was successful
				return
			}
		}

	}

	// Attempt to find the binary on the path
	out, err := exec.Command("rescan-scsi-bus.sh", "-a", "-r").CombinedOutput()
	if err != nil {
		log.Errorf("Could not rescan SCSI bus: %v. %v", err, string(out))
		return
	}

	log.Warn("Unable to find rescan-scsi-bus command!")
	return
}

// UdevSettle invokes the 'udevadm settle' command
func UdevSettle() error {
	// creating new storage and attaching it to a host can trigger a ripple of udev activity.
	log.Debug("Begin osutils.udevSettle")

	// back-to-back invoke /sbin/multipath, if it exists, to make sure that device discovery has settled down post rescan
	for i := 0; i < 2; i++ {
		Multipath()
		Multipath()

		// attempt to wait for inflight udev events to complete (should eventually timeout if they never complete)
		out, err := exec.Command("udevadm", "settle").CombinedOutput()
		if err != nil {
			// nothing to really do if it generates an error but log and return it
			log.Debugf("Error encountered in udevadm settle cmd: %v. %v", err, string(out))
			return err
		}
	}

	return nil
}

// Multipath invokes the 'multipath' command
func Multipath() (err error) {
	log.Debug("Begin osutils.multipath")

	if !MultipathDetected() {
		log.Debug("Skipping multipath command, /sbin/multipath doesn't exist.")
		return
	}

	out, err := exec.Command("multipath").CombinedOutput()
	if err != nil {
		// nothing to really do if it generates an error but log and return it
		log.Debugf("Error encountered in multipath cmd: %v. %v", err, string(out))
		return
	}
	return
}

// MultipathFlush invokes the 'multipath' commands to flush paths that have been removed
func MultipathFlush() (err error) {
	log.Debug("Begin osutils.multipathFlush")
	out, err := exec.Command("multipath", "-F").CombinedOutput()
	if err != nil {
		// nothing to really do if it generates an error but log and return it
		log.Debugf("Error encountered in multipath flush unused paths cmd: %v. %v", err, string(out))
		return
	}
	return
}

var mpChecked = false
var mpDetected = false
var mpMutex sync.Mutex

// MultipathDetected returns true if /sbin/multipath is installed and in the PATH
func MultipathDetected() bool {
	mpMutex.Lock()
	defer mpMutex.Unlock()

	if !mpChecked {
		_, errStat := Stat("/sbin/multipath")
		if errStat != nil {
			mpDetected = false
		} else {
			mpDetected = true
		}
		mpChecked = true
	}
	return mpDetected
}

// GetFSType returns the filesystem for the supplied device
func GetFSType(device string) string {
	log.Debugf("Begin osutils.GetFSType: %s", device)
	fsType := ""
	out, err := exec.Command("blkid", device).CombinedOutput()
	if err != nil {
		log.Debugf("Could not get FSType for device %v: %v. %v", device, err, string(out))
		return fsType
	}

	if strings.Contains(string(out), "TYPE=") {
		for _, v := range strings.Split(string(out), " ") {
			if strings.Contains(v, "TYPE=") {
				fsType = strings.Split(v, "=")[1]
				fsType = strings.Replace(fsType, "\"", "", -1)
			}
		}
	}
	return fsType
}

// FormatVolume creates a filesystem for the supplied device of the supplied type
func FormatVolume(device, fsType string) error {
	log.Debugf("Begin osutils.FormatVolume: %s, %s", device, fsType)
	cmd := "mkfs.ext4"
	if fsType == "xfs" {
		cmd = "mkfs.xfs"
	}
	log.Debug("Perform ", cmd, " on device: ", device)
	out, err := exec.Command(cmd, "-F", device).CombinedOutput()
	log.Debug("Result of mkfs cmd: ", string(out))
	return err
}

// Mount attaches the supplied device at the supplied location
func Mount(device, mountpoint string) error {
	log.Debugf("Begin osutils.Mount device: %s on: %s", device, mountpoint)
	out, err := exec.Command("mkdir", mountpoint).CombinedOutput()
	out, err = exec.Command("mount", device, mountpoint).CombinedOutput()
	log.Debug("Response from mount ", device, " at ", mountpoint, ": ", string(out))
	if err != nil {
		log.Errorf("Error in mount: %v.", err)
	}
	return err
}

// Umount detaches from the supplied location
func Umount(mountpoint string) error {
	log.Debugf("Begin osutils.Umount: %s", mountpoint)
	out, err := exec.Command("umount", mountpoint).CombinedOutput()
	log.Debug("Response from umount ", mountpoint, ": ", string(out))
	if err != nil {
		log.Errorf("Error in unmount: %v.", err)
	}
	return err
}

// IscsiadmCmd uses the 'iscsiadm' command to perform operations
func IscsiadmCmd(args []string) ([]byte, error) {
	log.Debugf("Begin osutils.iscsiadmCmd: iscsiadm %+v", args)
	out, err := exec.Command("iscsiadm", args...).CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprint("Error encountered running iscsiadm ", args)
		log.Errorf("%s: %v. %v", errMsg, err, string(out))
	}
	return out, err
}

// Login to iSCSI target
func LoginIscsiTarget(iqn, portal string) error {

	log.WithFields(log.Fields{
		"IQN":    iqn,
		"Portal": portal,
	}).Debug("Logging in to iSCSI target.")

	args := []string{"-m", "node", "-T", iqn, "-l", "-p", portal + ":3260"}

	if _, err := IscsiadmCmd(args); err != nil {
		log.Errorf("Error logging in to iSCSI target: %v.", err)
		return err
	}
	return nil
}

// LoginWithChap will login to the iSCSI target with the supplied credentials
func LoginWithChap(tiqn, portal, username, password, iface string) error {
	log.Debugf("Begin osutils.LoginWithChap: iqn: %s, portal: %s, username: %s, password=xxxx, iface: %s", tiqn, portal, username, iface)
	args := []string{"-m", "node", "-T", tiqn, "-p", portal + ":3260"}
	createArgs := append(args, []string{"--interface", iface, "--op", "new"}...)

	if out, err := exec.Command("iscsiadm", createArgs...).CombinedOutput(); err != nil {
		log.Errorf("Error running iscsiadm node create: %v. %v", err, string(out))
		return err
	}

	authMethodArgs := append(args, []string{"--op=update", "--name", "node.session.auth.authmethod", "--value=CHAP"}...)
	if out, err := exec.Command("iscsiadm", authMethodArgs...).CombinedOutput(); err != nil {
		log.Errorf("Error running iscsiadm set authmethod: %v. %v", err, string(out))
		return err
	}

	authUserArgs := append(args, []string{"--op=update", "--name", "node.session.auth.username", "--value=" + username}...)
	if out, err := exec.Command("iscsiadm", authUserArgs...).CombinedOutput(); err != nil {
		log.Errorf("Error running iscsiadm set authuser: %v. %v", err, string(out))
		return err
	}
	authPasswordArgs := append(args, []string{"--op=update", "--name", "node.session.auth.password", "--value=" + password}...)
	if out, err := exec.Command("iscsiadm", authPasswordArgs...).CombinedOutput(); err != nil {
		log.Errorf("Error running iscsiadm set authpassword: %v. %v", err, string(out))
		return err
	}
	loginArgs := append(args, []string{"--login"}...)
	if out, err := exec.Command("iscsiadm", loginArgs...).CombinedOutput(); err != nil {
		log.Errorf("Error running iscsiadm login: %v. %v", err, string(out))
		return err
	}
	return nil
}

func EnsureIscsiSession(hostDataIP string) error {

	// Ensure iSCSI is supported on system
	if !IscsiSupported() {
		return errors.New("iSCSI support not detected.")
	}

	// Ensure iSCSI session exists for the specified iSCSI portal
	sessionExists, err := IscsiSessionExists(hostDataIP)
	if err != nil {
		return fmt.Errorf("Could not check for iSCSI session. %v", err)
	}
	if !sessionExists {

		// Run discovery in case we haven't seen this target from this host
		targets, err := IscsiDiscovery(hostDataIP)
		if err != nil {
			return fmt.Errorf("Could not run iSCSI discovery. %v", err)
		}
		if len(targets) == 0 {
			return errors.New("iSCSI discovery found no targets.")
		}

		log.WithFields(log.Fields{
			"Targets": targets,
		}).Debugf("Found matching iSCSI targets.")

		// Determine which target matches the portal we requested
		targetIndex := -1
		for i, target := range targets {
			if target.PortalIP == hostDataIP {
				targetIndex = i
				break
			}
		}

		if targetIndex == -1 {
			return fmt.Errorf("iSCSI discovery found no targets with portal %s.", hostDataIP)
		}

		// To enable multipath, log in to each discovered target with the same IQN (target name)
		targetName := targets[targetIndex].TargetName
		for _, target := range targets {
			if target.TargetName == targetName {

				// Log in to target
				err = LoginIscsiTarget(target.TargetName, target.PortalIP)
				if err != nil {
					return fmt.Errorf("Login to iSCSI target failed. %v", err)
				}
			}
		}

		// Recheck to ensure a session is now open
		sessionExists, err = IscsiSessionExists(hostDataIP)
		if err != nil {
			return fmt.Errorf("Could not recheck for iSCSI session. %v", err)
		}
		if !sessionExists {
			return fmt.Errorf("Expected iSCSI session %v NOT found, please login to the iSCSI portal.", hostDataIP)
		}
	}

	log.Debugf("Found session to iSCSI portal %s.", hostDataIP)

	return nil
}
