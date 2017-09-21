# Change Log

[Releases](https://github.com/NetApp/netappdvp/releases)

## Changes since 17.10.0

**Fixes:**
- Added delete idempotency to the ontap-nas-economy driver as needed by Trident.
- Fixed an issue where qtrees with names near the 64-character limit could not
  be deleted.
- Fixed an issue with volume creation on all-flash E-series arrays.

**Enhancements:**
- Improved iSCSI rescan performance for ONTAP SAN and E-series plugins.

## 17.10.0

**Fixes:**
- Changed the SolidFire driver to tear down the iSCSI connection as part
  of the docker volume delete operation (Issue [#93](https://github.com/NetApp/netappdvp/issues/93)).
- nDVP plugin no longer requires /etc/iscsi bind mount even for NFS-only
  installs (Issue [#82](https://github.com/NetApp/netappdvp/issues/82)).
- Added pagination to ONTAP API calls, the lack of which in rare cases
  could cause nDVP to fail.
- Fixed issue where ONTAP NAS volumes were not mountable immediately
  after creation when using load-sharing mirrors for the SVM root
  volume (Issue [#84](https://github.com/NetApp/netappdvp/issues/84)).
- Fixed a bug where port names on E-Series could exceed the 30 character
  system limitation.
- Fixed an issue where large data sets could cause the log system to
  crash the plugin.

**Enhancements:**
- Added controller serial numbers to logs.
- Added ontap-nas-economy driver.
- Added aggregate validation to the ONTAP drivers
  (Issue [#92](https://github.com/NetApp/netappdvp/issues/92)).
- Added support for xfs and ext3 to the SolidFire, ONTAP SAN, and
  E-series drivers (Issue [#73](https://github.com/NetApp/netappdvp/issues/73)).
- Added iSCSI multipath support to SolidFire driver (Issue [#49](https://github.com/NetApp/netappdvp/issues/49)).
- Added support for NetApp Volume Encryption to the ONTAP drivers.

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
