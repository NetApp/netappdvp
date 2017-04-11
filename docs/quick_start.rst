.. _quick_start:

Quick Start
===========

This quick start is targeted at the Docker Managed Plugin method (Docker >= 1.13 / 17.03).  If you're using an earlier version of Docker, please refer to the documentation: :ref:`host-configuration`.

#. nDVP is supported on the following operating systems:

   * Debian
   * Ubuntu, 14.04+ if not using iSCSI multipathing, 15.10+ with iSCSI multipathing.
   * CentOS, 7.0+
   * RHEL, 7.0+

#. Verify your storage system meets the minimum requirements:

   * ONTAP: 8.3 or greater
   * SolidFire: ElementOS 7 or greater
   * E-Series: Web Services Proxy

#. Ensure you have Docker Engine 17.03 (nee 1.13) or above installed.

   .. code-block:: bash
   
     docker --version
   
   If your version is out of date, `follow the instructions for your distribution <https://docs.docker.com/engine/installation/>`_ to install or update.
   

#. Verify that the protocol prerequesites are installed and configured on your host.  See :ref:`host-configuration`.
   
   

#. Create a configuration file.  The default location is ``/etc/netappdvp/config.json``.  Be sure to use the correct options for your storage system.

   .. code-block:: bash
   
     # create a location for the config files
     sudo mkdir -p /etc/netappdvp
 
     # create the configuration file, see below for more configuration examples
     cat << EOF > /etc/netappdvp/config.json
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
     EOF

#. Start nDVP using the managed plugin system.

   .. code-block:: bash
   
     docker plugin install netapp/ndvp-plugin:17.04 --alias netapp --grant-all-permissions

#. Begin using nDVP to consume storage from the configured system.

   .. code-block:: bash
   
     # create a volume named "firstVolume"
     docker volume create -d netapp --name firstVolume
     
     # create a default volume at container instantiation
     docker run --rm -it --volume-driver netapp --volume secondVolume:/my_vol alpine ash
     
     # remove the volume "firstVolume"
     docker volume rm firstVolume


