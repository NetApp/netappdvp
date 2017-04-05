.. _ontap_vol_opts:

ONTAP Volume Options
====================

Volume create options for both NFS and iSCSI:

* ``size`` - the size of the volume, defaults to 1 GiB
* ``spaceReserve`` - thin or thick provision the volume, defaults to thin. Valid values are "none" (thin provisioned) or "volume" (thick provisioned).
* ``snapshotPolicy`` - this will set the snapshot policy to the desired value. The default is "none", meaning no snapshots will automatically be created for the volume. Unless modified by your storage administrator, a policy named "default" exists on all ONTAP systems which creates and retains six hourly, two daily, and two weekly snapshots. The data preserved in a snapshot can be recovered by browsing to the .snapshot directory in any directory in the volume.

NFS has two additional options that aren't relevant when using iSCSI:

* ``unixPermissions`` - this controls the permission set for the volume itself. By default the permissions will be set to ``---rwxr-xr-x``, or in numerical notation ``0755``, and root will be the owner. Either the text or numerical format will work.
* ``snapshotDir`` - setting this to ``true`` will make the .snapshot directory visible to clients accessing the volume. The default value is ``false``, meaning that access to snapshot data is disabled by default.  Some images, for example the official MySQL image, don't function as expected when the .snapshot directory is visible.
* ``exportPolicy`` - sets the export policy to be used for the volume.  The default is ``default``.
* ``securityStyle`` - sets the security style to be used for access to the volume.  The default is ``unix``, value should be ``unix`` or ``mixed``.

Using these options during the docker volume create operation is super simple, just provide the option and the value using the ``-o`` operator during the CLI operation.  These override any equivalent vales from the JSON configuration file.

.. code-block:: bash

   # create a 10GiB volume
   docker volume create -d netapp --name demo -o size=10G

   # create a 100GiB volume with snapshots
   docker volume create -d netapp --name demo -o size=100G -o snapshotPolicy=default

   # create a volume which has the setUID bit enabled
   docker volume create -d netapp --name demo -o unixPermissions=4755

The minimum volume size is 20MiB.