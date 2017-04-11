.. _host-configuration:

Host Configuration
==================

Prerequesites
-------------

* nDVP is supported on the following operating systems:

  * Debian
  * Ubuntu, 14.04+ if not using iSCSI multipathing, 15.10+ with iSCSI multipathing.
  * CentOS, 7.0+
  * RHEL, 7.0+

* Verify your storage system meets the minimum requirements:

  * ONTAP: 8.3 or greater
  * SolidFire: ElementOS 7 or greater
  * E-Series: Santricity

NFS
---

Install the following system packages:

* RHEL / CentOS

  .. code-block:: bash

    sudo yum install -y nfs-utils

* Ubuntu / Debian

  .. code-block:: bash

    sudo apt-get install -y nfs-common

iSCSI
-----

* RHEL / CentOS

  #. Install the following system packages:

     .. code-block:: bash

       sudo yum install -y lsscsi iscsi-initiator-utils sg3_utils device-mapper-multipath

  #. Start the multipathing daemon:

     .. code-block:: bash

       sudo mpathconf --enable --with_multipathd y

  #. Ensure that `iscsid` and `multipathd` are enabled and running:

     .. code-block:: bash

       sudo systemctl enable iscsid multipathd
       sudo systemctl start iscsid multipathd

  #. Discover the iSCSI targets:

     .. code-block:: bash

       sudo iscsiadm -m discoverydb -t st -p <DATA_LIF_IP> --discover

  #. Login to the discovered iSCSI targets:

     .. code-block:: bash

       sudo iscsiadm -m node -p <DATA_LIF_IP> --login

  #. Start and enable ``iscsi``:

     .. code-block:: bash

       sudo systemctl enable iscsi
       sudo systemctl start iscsi

* Ubuntu / Debian

  #. Install the following system packages:

     .. code-block:: bash

       sudo apt-get install -y open-iscsi lsscsi sg3-utils multipath-tools scsitools

  #. Enable multipathing:

     .. code-block:: bash

       sudo tee /etc/multipath.conf <<-'EOF'
       defaults {
           user_friendly_names yes
           find_multipaths yes
       }
       EOF

       sudo service multipath-tools restart

  #. Ensure that ``iscsid`` and ``multipathd`` are running:

     .. code-block:: bash

       sudo service open-iscsi start
       sudo service multipath-tools start


  #. Discover the iSCSI targets:

     .. code-block:: bash

       sudo iscsiadm -m discoverydb -t st -p <DATA_LIF_IP> --discover

  #. Login to the discovered iSCSI targets:

     .. code-block:: bash

       sudo iscsiadm -m node -p <DATA_LIF_IP> --login

Traditional Install Method (Docker <= 1.12)
-------------------------------------------

#. Ensure you have Docker version 1.10 or above
   
   .. code-block:: bash
   
      docker --version

   If your version is out of date, update to the latest.
   
   .. code-block:: bash
   
      curl -fsSL https://get.docker.com/ | sh

   Or, `follow the instructions for your distribution <https://docs.docker.com/engine/installation/>`_.
   

#. After ensuring the correct version of Docker is installed, install and configure the NetApp Docker Volume Plugin.  Note, you will need to ensure that NFS and/or iSCSI is configured for your system.  See the installation instructions below for detailed information on how to do this.

   .. code-block:: bash

      # download and unpack the application
      wget https://github.com/NetApp/netappdvp/releases/download/v17.04.0/netappdvp-17.04.0.tar.gz
      tar zxf netappdvp-17.04.0.tar.gz

      # move to a location in the bin path
      sudo mv netappdvp /usr/local/bin
      sudo chown root:root /usr/local/bin/netappdvp
      sudo chmod 755 /usr/local/bin/netappdvp

      # create a location for the config files
      sudo mkdir -p /etc/netappdvp

      # create the configuration file, see below for more configuration examples
      cat << EOF > /etc/netappdvp/ontap-nas.json
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

#. After placing the binary and creating the configuration file(s), start the nDVP daemon using the desired configuration file.

   **Note:** Unless specified, the default name for the volume driver will be "netapp".

   .. code-block:: bash

     sudo netappdvp --config=/etc/netappdvp/ontap-nas.json 


#. Once the daemon is started, create and manage volumes using the Docker CLI interface.

   .. code-block:: bash

      docker volume create -d netapp --name ndvp_1


   Provision Docker volume when starting a container:

   .. code-block:: bash

      docker run --rm -it --volume-driver netapp --volume ndvp_2:/my_vol alpine ash

   Destroy docker volume:

   .. code-block:: bash

      docker volume rm ndvp_1
      docker volume rm ndvp_2

Starting nDVP at System Startup
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

A sample unit file for systemd based systems can be found at ``contrib/netappdvp.service.example`` in the git repo.  To use the file, with CentOS/RHEL:

.. code-block:: bash
   
   # copy the file to the correct location.  you must use unique names for the
   # unit files if you have more than one instance running
   cp contrib/netappdvp.service.example /usr/lib/systemd/system/netappdvp.service
   
   # edit the file, change the description (line 2) to match the driver name and the
   # configuration file path (line 9) to reflect your envrionment.
   
   # reload systemd for it to ingest changes
   systemctl daemon-reload
   
   # enable the service, note this name will change depending on what you named the
   # file in the /usr/lib/systemd/system directory
   systemctl enable netappdvp
   
   # start the service, see note above about service name
   systemctl start netappdvp
   
   # view the status
   systemctl status netappdvp

Note that anytime the unit file is modified you will need to issue the command ``systemctl daemon-reload`` for it to be aware of the changes.

Docker Managed Plugin Method (Docker >= 1.13 / 17.03)
-----------------------------------------------------

**Note:** If you have used nDVP pre-1.13/17.03 in the traditional daemon method, please ensure that you stop the nDVP process and restart your Docker daemon before using the managed plugin method.

.. code-block:: bash

   # stop all running instances
   pkill /usr/local/bin/netappdvp
   
   # restart docker
   systemctl restart docker

**nDVP Specific Plugin Startup Options**

* ``--config`` - Specify the configuration file the plugin will use.  Only the file name should be specified, e.g. ``gold.json``, the location must be ``/etc/netappdvp`` on the host system.  The default is ``config.json``.
* ``--log-level`` - Specify the logging level (``debug``, ``info``, ``warn``, ``error``, ``fatal``).  The default is ``info``.
* ``--debug`` - Specify whether debug logging is enabled.  Default is false.  Overrides --log-level if true.

**Installing the Managed Plugin**
   
#. Ensure you have Docker Engine 17.03 (nee 1.13) or above installed.

   .. code-block:: bash
   
     docker --version
   
   If your version is out of date, `follow the instructions for your distribution <https://docs.docker.com/engine/installation/>`_ to install or update.

#. Create a configuration file.  The config file must be located in the ``/etc/netappdvp`` directory.  The default filename is ``config.json``, however you can use any name you choose by specifying the ``--config`` option with the file name.  Be sure to use the correct options for your storage system.

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
