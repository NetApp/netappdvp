ONTAP Configuration
===================

User Permissions
----------------

nDVP does not need full permissions on the ONTAP cluster and should not be used with the cluster-level admin account.  Below are the ONTAP CLI comands to create a dedicated user for nDVP with specific permissions.

.. code-block:: bash

  # create a new nDVP role
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname DEFAULT -access none
  
  # grant common nDVP permissions
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "event generate-autosupport-log" -access all
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "network interface" -access readonly
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "version" -access readonly
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "vserver" -access readonly
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "vserver nfs show" -access readonly
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "volume" -access all
  
  # grant iscsi nDVP permissions
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "vserver iscsi show" -access readonly
  security login role create -vserver [VSERVER] -role ndvp_role -cmddirname "lun" -access all
  
  # create a new nDVP user with nDVP role
  security login create -vserver [VSERVER] -username ndvp_user -role ndvp_role -application ontapi -authmethod password

Configuration File Options
--------------------------

In addition to the global configuration values above, when using clustered Data ONTAP, these top level options are available.

+-----------------------+--------------------------------------------------------------------------+------------+
| Option                | Description                                                              | Example    |
+=======================+==========================================================================+============+
| ``managementLIF``     | IP address of clustered Data ONTAP management LIF                        | 10.0.0.1   |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``dataLIF``           | IP address of protocol lif; will be derived if not specified             | 10.0.0.2   |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``svm``               | Storage virtual machine to use (req, if management LIF is a cluster LIF) | svm_nfs    |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``username``          | Username to connect to the storage device                                | vsadmin    |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``password``          | Password to connect to the storage device                                | netapp123  |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``aggregate``         | Aggregate to use for volume/LUN provisioning                             | aggr1      |
+-----------------------+--------------------------------------------------------------------------+------------+

For the ontap-nas driver, an additional top level option is available. For NFS host configuration, see also: http://www.netapp.com/us/media/tr-4067.pdf

+-----------------------+--------------------------------------------------------------------------+------------+
| Option                | Description                                                              | Example    |
+=======================+==========================================================================+============+
| ``nfsMountOptions``   | Fine grained control of NFS mount options; defaults to "-o nfsvers=3"    |-o nfsvers=4|
+-----------------------+--------------------------------------------------------------------------+------------+

Also, when using clustered Data ONTAP, these default option settings are available to avoid having to specify them on every volume create.

+-----------------------+--------------------------------------------------------------------------+------------+
| Defaults Option       | Description                                                              | Example    |
+=======================+==========================================================================+============+
| ``spaceReserve``      | Space reservation mode; "none" (thin provisioned) or "volume" (thick)    | none       |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``snapshotPolicy``    | Snapshot policy to use, default is "none"                                | none       |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``splitOnClone``      | Split a clone from its parent upon creation, defaults to "false"         | false      |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``unixPermissions``   | NAS option for provisioned NFS volumes, defaults to "777"                | 777        |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``snapshotDir``       | NAS option for access to the .snapshot directory, defaults to "false"    | false      |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``exportPolicy``      | NAS option for the NFS export policy to use, defaults to "default"       | default    |
+-----------------------+--------------------------------------------------------------------------+------------+
| ``securityStyle``     | NAS option for access to the provisioned NFS volume, defaults to "unix"  | mixed      |
+-----------------------+--------------------------------------------------------------------------+------------+

Example ONTAP Config Files
--------------------------

**NFS Example for ontap-nas driver**

.. code-block:: json

    {
        "version": 1,
        "storageDriverName": "ontap-nas",
        "managementLIF": "10.0.0.1",
        "dataLIF": "10.0.0.2",
        "svm": "svm_nfs",
        "username": "vsadmin",
        "password": "netapp123",
        "aggregate": "aggr1",
        "defaults": {
          "size": "10G",
          "spaceReserve": "none",
          "exportPolicy": "default"
        }
    }

**iSCSI Example for ontap-san driver**

.. code-block:: json

    {
        "version": 1,
        "storageDriverName": "ontap-san",
        "managementLIF": "10.0.0.1",
        "dataLIF": "10.0.0.3",
        "svm": "svm_iscsi",
        "username": "vsadmin",
        "password": "netapp123",
        "aggregate": "aggr1"
    }
