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

In addition to the global configuration values above, when using clustered Data ONTAP, these options are available.

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
      "aggregate": "aggr1"
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
