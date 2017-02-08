Using nDVP
==========

Creating and consuming storage from ONTAP, SolidFire, and/or E-Series systems is easy with nDVP.  Simply use the standard ``docker volume`` commands with the nDVP driver name specified when needed.

Volume Driver CLI Options
-------------------------

Each storage driver has a different set of options which can be provided at volume creation time to customize the outcome.  Refer to the documentation below for your configured storage system to determine which options apply.

* :ref:`ontap_vol_opts`
* :ref:`sf_vol_opts`
* :ref:`es_vol_opts`

Create a Volume
---------------

.. code-block:: bash

   # create a volume with an nDVP driver using the default name
   docker volume create -d netapp --name firstVolume
   
   # create a volume with a specific nDVP instance
   docker volume create -d ntap_bronze --name bronzeVolume

If no options are specified, the defaults for the driver are used.  The defaults are documented on the page for the storage driver you're using below.

Destroy a Volume
-------------------

.. code-block:: bash

   # destroy the volume just like any other Docker volume
   docker volume rm firstVolume

Volume Cloning
--------------

For ONTAP and SolidFire only, volumes can be cloned using the Docker Volume Plugin.

.. code-block:: bash

   # inspect the volume to enumerate snapshots
   docker volume inspect <volume_name>
   
   # create a new volume from an existing volume.  this will result in a new snapshot being created
   docker volume create -d <driver_name> --name <new_name> -o from=<source_docker_volume>
   
   # create a new volume from an existing snapshot on a volume.  this will not create a new snapshot
   docker volume create -d <driver_name> --name <new_name> -o from=<source_docker_volume> -o fromSnapshot=<source_snap_name>

Here is an example of that in action:

.. code-block:: bash

   [me@host ~]$ docker volume inspect firstVolume
   
   [
       {
           "Driver": "ontap-nas",
           "Labels": null,
           "Mountpoint": "/var/lib/docker-volumes/ontap-nas/netappdvp_firstVolume",
           "Name": "firstVolume",
           "Options": {},
           "Scope": "global",
           "Status": {
               "Snapshots": [
                   {
                       "Created": "2017-02-10T19:05:00Z",
                       "Name": "hourly.2017-02-10_1505"
                   }
               ]
           }
       }
   ]
    
   [me@host ~]$ docker volume create -d ontap-nas --name clonedVolume -o from=firstVolume
   clonedVolume
   
   [me@host ~]$ docker volume rm clonedVolume
   [me@host ~]$ docker volume create -d ontap-nas --name volFromSnap -o from=firstVolume -o fromSnapshot=hourly.2017-02-10_1505
   volFromSnap
   
   [me@host ~]$ docker volume rm volFromSnap


