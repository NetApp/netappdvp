Installing nDVP
===============

The first step to installing the NetApp Docker Volume Plugin is to ensure that your host is configured for the protocol you intend to use, NFS or iSCSI.  Once that is complete you will want to create a configuration file which details the nDVP instance configuration, and finally instantiate the daemon on the host.

Follow the directions in the :ref:`quick_start` for downloading the nDVP binary and installing it to your host.

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   host_config
   ndvp_global_config
   ndvp_ontap_config
   ndvp_sf_config
   ndvp_e_config
   multi_instance

