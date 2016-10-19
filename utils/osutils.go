// Copyright 2016 NetApp, Inc. All Rights Reserved.

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

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
			log.Error("Error encountered gathering df output: ", err)
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

// GetInitiatorIqns returns parsed contents of /etc/iscsi/initiatorname.iscsi
func GetInitiatorIqns() ([]string, error) {
	log.Debug("Begin osutils.GetInitiatorIqns")
	var iqns []string
	out, err := exec.Command("cat", "/etc/iscsi/initiatorname.iscsi").CombinedOutput()
	if err != nil {
		log.Error("Error encountered gathering initiator names: ", err)
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

// WaitForPathToExist retries every second, up to numTries times, for the specified fileName to show up
func WaitForPathToExist(fileName string, numTries int) bool {
	log.Debugf("Begin osutils.waitForPathToExist fileName: %v", fileName)
	for i := 0; i < numTries; i++ {
		_, err := os.Stat(fileName)
		if err == nil {
			log.Debugf("path found for fileName: %v", fileName)
			return true
		}
		if err != nil && !os.IsNotExist(err) {
			return false
		}
		time.Sleep(time.Second)
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
		_, err := os.Lstat("/sbin/multipath")
		if os.IsNotExist(err) {
			log.Debug("Skipping multipath check, /sbin/multipath doesn't exist")
		} else {
			lsblkCmd := fmt.Sprintf("lsblk %v -n -o name,type -r | grep mpath | cut -f1 -d\\ ", devFile)
			log.Debugf("running 'sh -c %v'", lsblkCmd)
			out2, err2 := exec.Command("sh", "-c", lsblkCmd).CombinedOutput()
			if err2 != nil {
				// this can be fine, for instance could be a floppy or cd-rom, later logic will error if we never find our device
				log.Debugf("could not run multipath check against device: %v error: %v", devFile, err2)
			} else {
				md := strings.Split(strings.TrimSpace(string(out2)), " ")
				if md != nil && len(md) > 0 && len(md[0]) > 0 {
					log.Debug("Found md: ", md)
					multipathDevFileToCheck := "/dev/mapper/" + md[0]
					_, err := os.Lstat(multipathDevFileToCheck)
					if !os.IsNotExist(err) {
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

	// first, list w/out iscsi target info
	var info1 []ScsiDeviceInfo
	info1, err1 := LsscsiCmd(nil)
	if err1 != nil {
		return nil, err1
	}

	// now, list w/ iscsi target info
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
	_, err := exec.Command("iscsiadm", "-h").CombinedOutput()
	if err != nil {
		log.Debug("iscsiadm tools not found on this host")
		return false
	}
	return true
}

// IscsiDiscovery uses the 'iscsiadm' command to perform discovery
func IscsiDiscovery(portal string) (targets []string, err error) {
	log.Debugf("Begin osutils.IscsiDiscovery (portal: %s)", portal)
	out, err := exec.Command("iscsiadm", "-m", "discovery", "-t", "sendtargets", "-p", portal).CombinedOutput()
	if err != nil {
		log.Error("Error encountered in sendtargets cmd: ", out)
		return
	}
	targets = strings.Split(string(out), "\n")
	return
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
		log.Errorf("Problem checking iscsi sessions error: %v", err)
		return nil, err
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

// IscsiDisableDelete logout from the supplied target and remove the iscsi device
func IscsiDisableDelete(tgt *IscsiTargetInfo) (err error) {
	log.Debugf("Begin osutils.IscsiDisableDelete: %v", tgt)
	_, err = exec.Command("sudo", "iscsiadm", "-m", "node", "-T", tgt.Iqn, "--portal", tgt.IP, "-u").CombinedOutput()
	if err != nil {
		log.Debugf("Error during iscsi logout: ", err)
		//return
	}
	_, err = exec.Command("sudo", "iscsiadm", "-m", "node", "-o", "delete", "-T", tgt.Iqn).CombinedOutput()
	return
}

// IscsiSessionExists checks to see if a session exists to the sepecified portal
func IscsiSessionExists(portal string) (bool, error) {
	log.Debugf("Begin osutils.IscsiSessionExists")

	sessionInfo, err := GetIscsiSessionInfo()
	if err != nil {
		log.Errorf("Problem checking iscsi sessions error: %v", err)
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

	// look for version of rescan-scsi-bus in known locations
	var rescanCommands []string = []string{"/sbin/rescan-scsi-bus", "/sbin/rescan-scsi-bus.sh", "/bin/rescan-scsi-bus.sh", "/usr/bin/rescan-scsi-bus.sh"}
	for _, rescanCommand := range rescanCommands {
		_, err = os.Lstat(rescanCommand)
		// The command exists in this location
		if !os.IsNotExist(err) {
			out, rescanErr := exec.Command(rescanCommand, "-a", "-r").CombinedOutput()
			// We encountered an error condition
			if rescanErr != nil {
				log.Error("Error encountered in rescan-scsi-bus cmd: ", out)
				return rescanErr
				// The command was successful
			} else {
				return
			}
		}

	}

	//Attempt to find the binary on the path
	out, err := exec.Command("rescan-scsi-bus.sh", "-a", "-r").CombinedOutput()
	if err != nil {
		log.Error("Error encountered in rescan-scsi-bus cmd: ", out)
		return
	}

	log.Warn("Unable to find rescan-scsi-bus command!")
	return
}

// MultipathFlush uses the 'multipath' commands to flush paths that have been removed
func MultipathFlush() (err error) {
	log.Debugf("Begin osutils.multipathFlush")
	out, err := exec.Command("multipath", "-F").CombinedOutput()
	if err != nil {
		// nothing to really do if it generates an error but log and return it
		log.Debugf("Error encountered in multipath flush unused paths cmd: ", out)
		return
	}
	return
}

// GetFSType returns the filesystem for the supplied device
func GetFSType(device string) string {
	log.Debugf("Begin osutils.GetFSType: %s", device)
	fsType := ""
	out, err := exec.Command("blkid", device).CombinedOutput()
	if err != nil {
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
		log.Error("Error in mount: ", err)
	}
	return err
}

// Umount detaches from the supplied location
func Umount(mountpoint string) error {
	log.Debugf("Begin osutils.Umount: %s", mountpoint)
	out, err := exec.Command("umount", mountpoint).CombinedOutput()
	log.Debug("Response from umount ", mountpoint, ": ", out)
	if err != nil {
		log.Error("Error in unmount: ", err)
	}
	/*
		out, _ = exec.Command("rmdir", mountpoint).CombinedOutput()
		log.Debug("Response from rmdir ", mountpoint, ": ", out)
	*/
	return err
}

// IscsiadmCmd uses the 'iscsiadm' command to perform operations
func IscsiadmCmd(args []string) ([]byte, error) {
	log.Debugf("Begin osutils.iscsiadmCmd: iscsiadm %+v", args)
	resp, err := exec.Command("iscsiadm", args...).CombinedOutput()
	if err != nil {
		log.Error("Error encountered running iscsiadm ", args, ": ", resp)
		log.Error("Error message: ", err)
	}
	return resp, err
}

// LoginWithChap will login to the iscsi target with the supplied credentials
func LoginWithChap(tiqn, portal, username, password, iface string) error {
	log.Debugf("Begin osutils.LoginWithChap: iqn: %s, portal: %s, username: %s, password=xxxx, iface: %s", tiqn, portal, username, iface)
	args := []string{"-m", "node", "-T", tiqn, "-p", portal + ":3260"}
	createArgs := append(args, []string{"--interface", iface, "--op", "new"}...)

	if _, err := exec.Command("iscsiadm", createArgs...).CombinedOutput(); err != nil {
		log.Error("Error running iscsiadm node create: ", err)
		return err
	}

	authMethodArgs := append(args, []string{"--op=update", "--name", "node.session.auth.authmethod", "--value=CHAP"}...)
	if out, err := exec.Command("iscsiadm", authMethodArgs...).CombinedOutput(); err != nil {
		log.Error("Error running iscsiadm set authmethod: ", err, "{", out, "}")
		return err
	}

	authUserArgs := append(args, []string{"--op=update", "--name", "node.session.auth.username", "--value=" + username}...)
	if _, err := exec.Command("iscsiadm", authUserArgs...).CombinedOutput(); err != nil {
		log.Error("Error running iscsiadm set authuser: ", err)
		return err
	}
	authPasswordArgs := append(args, []string{"--op=update", "--name", "node.session.auth.password", "--value=" + password}...)
	if _, err := exec.Command("iscsiadm", authPasswordArgs...).CombinedOutput(); err != nil {
		log.Error("Error running iscsiadm set authpassword: ", err)
		return err
	}
	loginArgs := append(args, []string{"--login"}...)
	if _, err := exec.Command("iscsiadm", loginArgs...).CombinedOutput(); err != nil {
		log.Error("Error running iscsiadm login: ", err)
		return err
	}
	return nil
}
