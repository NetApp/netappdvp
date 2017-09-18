# Change Log

[Releases](https://github.com/NetApp/netappdvp/releases)

## Changes since 17.07.0

**Fixes:**
- Changed the Solidfire driver to tear down the iSCSI connection as part of the docker volume delete operation
- Added pagination to ONTAP API calls, the lack of which in rare cases could cause nDVP to fail
- Fixed issue where ONTAP NAS volumes were not mountable immediately after creation when using load-sharing mirrors for the SVM root volume

**Enhancements:**
- Added controller serial numbers to logs
- Added ontap-nas-economy driver
- Added aggregate validation to the ONTAP drivers
- Added support for xfs and ext3 to the SolidFire, ONTAP SAN, and E-series drivers

## 17.07.0

**Fixes:**
- Increased wait time for SolidFire devices to appear

**Enhancements:**

- Allow customizing NFS mount options via config file
- Updated builds to use Go 1.8
- Allow ONTAP to split a clone from its parent upon creation
- Improved efficiency of ONTAP LUN ID selection
- Improved efficiency of ONTAP volume list
- Solidfire volumes now have 512e enabled by default (previously defaulted to 4k block size)
- Added options to toggle Solidfire's 512e setting in config file and at volume create time
- Added the ability to override QoS values when cloning a volume or snapshot on Solidfire

## 17.04.0

**Fixes:**

- Stopped using redundant debug parameter from config file.
- Prevent creating volumes with duplicate names on SolidFire.
- Fixed cloning from a snapshot in SolidFire.
- Fixed a rare SliceNotRegistered error when cloning volumes with SolidFire.
- Resolved an iSCSI attachment bug when using multipathing.

**Enhancements:**

- Added release notes (CHANGELOG.md).
- Added minimum ONTAP version check.
- Logging enhancements: log rotation, log level control, simultaneous logging
to console and file.
- Added default ONTAP volume creation options to config.json that were previously only available via '-o'.
- Added ONTAP securityStyle option handling.
- Added default volume size to config file.
- Standardized volume size format across ONTAP, SolidFire, and E-series drivers.
- Reduced log spam when debug logging is enabled with SolidFire.
- Moved README.md contents to [Read the Docs](http://netappdvp.readthedocs.io/en/latest/).
- `.snapshots` directory is now hidden in ONTAP NFS mounts by default.
