Multiple Instances of nDVP
==========================

Multiple instances of nDVP are needed when you desire to have multiple storage configurations available simultaneously.  The key to multiple instances is to give them different names using the ``--alias`` option with the containerized plugin, or ``--volume-driver`` option when instantiating the nDVP driver on the host.

**Docker Managed Plugin (Docker >= 1.13 / 17.03)**

#. Launch the first instance specifying an alias and configuration file
   
   .. code-block:: bash
   
      docker plugin install store/netapp/ndvp-plugin:17.04.0 --alias silver --config silver.json --grant-all-permissions
   
#. Launch the second instance, specifying a different alias and configuration file
   
   .. code-block:: bash
   
      docker plugin install store/netapp/ndvp-plugin:17.04.0 --alias gold --config gold.json --grant-all-permissions

#. Create volumes specifying the alias as the driver name
   
   .. code-block:: bash
      
      # gold volume
      docker volume create -d gold --name ntapGold
      
      # silver volume
      docker volume create -d silver --name ntapSilver


**Traditional (Docker <=1.12)**

#. Launch the plugin with an NFS configuration using a custom driver ID:

    .. code-block:: bash
    
       sudo netappdvp --volume-driver=netapp-nas --config=/path/to/config-nfs.json
       
#. Launch the plugin with an iSCSI configuration using a custom driver ID:

    .. code-block:: bash
    
       sudo netappdvp --volume-driver=netapp-san --config=/path/to/config-iscsi.json

#. Provision Docker volumes each driver instance:

   * NFS
     
     .. code-block:: bash
     
        docker volume create -d netapp-nas --name my_nfs_vol

   * iSCSI
   
     .. code-block:: bash
     
        docker volume create -d netapp-san --name my_iscsi_vol
