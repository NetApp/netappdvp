# NetApp Docker Volume Plugin

The NetApp Docker Volume Plugin (nDVP) provides direct integration with the Docker ecosystem for NetApp's ONTAP, E-Series, and SolidFire storage platforms. The nDVP package supports the provisioning and management of storage resources from the storage platform to Docker hosts, with a robust framework for adding additional platforms in the future.

Multiple instances of the nDVP can run concurrently on the same host.  The allows simultaneous connections to multiple storage systems and storage types, with the ablity to customize the storage used for the Docker volume(s).

- [Quick Start](#quick-start)
- [Running Multiple nDVP Instances](#running-multiple-ndvp-instances)
- [Configuring your Docker host for NFS or iSCSI](#configuring-your-docker-host-for-nfs-or-iscsi)
    - [NFS](#nfs)
        - [RHEL / CentOS](#rhel-centos)
        - [Ubuntu / Debian](#ubuntu-debian)
    - [iSCSI](#iscsi)
        - [RHEL / CentOS](#rhel--centos-1)
        - [Ubuntu / Debian](#ubuntu--debian-1)
- [Global Config File Variables](#global-configuration-file-variables)
    - [Storage Prefix](#storage-prefix)
- [ONTAP Config File Variables](#ontap-config-file-variables)
    - [Example ONTAP Config Files](#example-ontap-config-files)
- [E-Series Config File Variables](#e-series-config-file-variables)
    - [Example E-Series Config Files](#example-e-series-config-files)
    - [E-Series Array Setup Notes](#e-series-array-setup-notes)
- [SolidFire Config File Variables](#solidfire-config-file-variables)
    - [Example Solidfire Config Files](#example-solidfire-config-files)

## Quick Start

1. Ensure you have Docker version 1.10 or above.

    ```bash
    docker --version
    ```

    If your version is out of date, update to the latest.

    ```bash
    curl -fsSL https://get.docker.com/ | sh
    ```

    Or, [follow the instructions for your distribution](https://docs.docker.com/engine/installation/).

2. After ensuring the correct version of Docker is installed, install and configure the NetApp Docker Volume Plugin.  Note, you will need to ensure that NFS and/or iSCSI is configured for your system.  See the installation instructions below for detailed information on how to do this.

	```bash
# download and unpack the application
wget -O /tmp/netappdvp-1.4.0.tar.gz https://github.com/NetApp/netappdvp/releases/download/v1.4.0/netappdvp-1.4.0.tar.gz 
cd /tmp
tar zxf netappdvp-1.4.0.tar.gz

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
   "svm": "docker_svm",
   "username": "vsadmin",
   "password": "password",
   "aggregate": "aggr",
   "storagePrefix": "netappdvp_",
}
EOF
	```    

3. After placing the binary and creating the configuration file(s), start the nDVP daemon using the desired configuration file.

    **Note:** Unless specified, the default name for the volume driver will be "netapp".

	```bash
    sudo netappdvp --config=/etc/netappdvp/ontap-nas.json
	```

4. Once the daemon is started, create and manage volumes using the Docker CLI interface.

    ```bash
    docker volume create -d netapp --name ndvp_1
    ```

    Provision Docker volume when starting a container:

    ```bash
    docker run --rm -it --volume-driver netapp --volume ndvp_2:/my_vol alpine ash
    ```

    Destroy docker volume:

    ```bash
    docker volume rm ndvp_1
    docker volume rm ndvp_2
    ```

## Running Multiple nDVP Instances

1. Launch the plugin with an NFS configuration using a custom driver ID:

    ```bash
    sudo netappdvp --volume-driver=netapp-nas --config=/path/to/config-nfs.json
    ```

2. Launch the plugin with an iSCSI configuration using a custom driver ID:

    ```bash
    sudo netappdvp --volume-driver=netapp-san --config=/path/to/config-iscsi.json
    ```

3. Provision Docker volumes each driver instance:

    **NFS**
    ```bash
    docker volume create -d netapp-nas --name my_nfs_vol
    ```

    **iSCSI**

    ```bash
    docker volume create -d netapp-san --name my_iscsi_vol
    ```
## Systemd service and startup scripts

1. Create the netappdvp systemd service

    **systemd**
    ```bash
    cat << EOF > /etc/systemd/system/netappdvp.service
    [Unit]
    Description=Start Netapp Docker Volume Plugin
    Before=docker.service
    After=network.target
    Requires=docker.service

    [Service]
    ExecStart=/etc/netappdvp/netappdvp_startup.sh
    
    [Install]
    WantedBy=multi-user.target
    EOF
    ```

	**configure startup options**
    ```bash	
    sudo systemctl enable /etc/systemd/system/netappdvp.service
    sudo systemctl daemon-reload
    ```

2. Create the startup script
  
    **bash script**
    ```bash
    cat << EOF > /etc/netappdvp/netappdvp_startup.sh
    #!/bin/bash
    /usr/local/bin/netappdvp --config=/etc/netappdvp/ontap-nas.json
    EOF
    
    chmod +x /etc/netappdvp/netappdvp_startup.sh
    ```

  **restart services**
    ```bash
	systemctl restart netappdvp
	systemctl restart docker
    ```
## Configuring your Docker host for NFS or iSCSI

### NFS

Install the following system packages:

#### RHEL / CentOS

```bash
sudo yum install -y nfs-utils
```

#### Ubuntu / Debian

```bash
sudo apt-get install -y nfs-common
```

### iSCSI

#### RHEL / CentOS

Other things needed
* Add in an iSCSI NIC (this only handles a single source NIC, so make it a failover NIC)  
* Configure an IP address on your iSCSI NIC  
* get the name of your iSCSI NIC 'connection' with nmcli con show  

Fill out the variables below and the rest is pretty much blind copy and paste.  

```bash  
#some variables to smooth out the process, change the areas contained in <foo>  
ISCSI_NIC_NAME="<interfacename>"  
SVM_NAME="<SVM name>"  
SVM_MGMT_IP="<SVM management IP>"  
SVM_ISCSI_IP1="<SVM iSCSI LIF IP1>"  
SVM_ISCSI_IP2="<SVM iSCSI LIF IP2>"  
VSADMIN_PASSWORD="<password>"  
AGGREGATE_NAME="<aggregate name>"  
STORAGE_PREFIX="<storage prefix>"
IGROUP="$(hostname -s)"  
  
#install dependencies  
sudo yum install -y lsscsi iscsi-initiator-utils sg3_utils device-mapper-multipath

#append hostname to IQN so its easy to identify
sudo printf "$(< /etc/iscsi/initiatorname.iscsi)-$(hostname -s)\n" > /etc/iscsi/initiatorname.iscsi

#configure iscsi virtual interface and bind it to the iscsi nics mac address
sudo iscsiadm -m iface -I ieth1 --op=new
sudo iscsiadm -m iface -I ieth1 --op=update -n iface.net_ifacename -v ${ISCSI_NIC_NAME}

#change iscsi interface to have a 9k mtu and not submit dns registration
sudo nmcli con mod ${ISCSI_NIC_CON_NAME} 802-3-ethernet.mtu 9000
sudo nmcli device mod ${ISCSI_NIC_NAME} ipv6.dhcp-send-hostname no
sudo nmcli device mod ${ISCSI_NIC_NAME} ipv4.dhcp-send-hostname no

#enable multipathd
sudo mpathconf --enable --with_multipathd y
sudo systemctl enable iscsid multipathd
sudo systemctl start iscsid multipathd

#change iscsi session timeout per netapps recommended best practice 
sudo sed -i 's/node.session.timeo.replacement_timeout = 120/node.session.timeo.replacement_timeout = 5/' /etc/iscsi/iscsid.conf

#configure discovery for both iscsi IPs
sudo iscsiadm -m discoverydb -t st -p ${SVM_ISCSI_IP1} --discover -I ieth1
sudo iscsiadm -m discoverydb -t st -p ${SVM_ISCSI_IP2} --discover -I ieth1
sudo iscsiadm -m node -p ${SVM_ISCSI_IP1} -I ieth1 --login
sudo iscsiadm -m node -p ${SVM_ISCSI_IP2} -I ieth1 --login

#enable iscsi
sudo systemctl enable iscsi
sudo systemctl start iscsi
```

#### Ubuntu / Debian

1. Install the following system packages:

    ```bash
    sudo apt-get install -y open-iscsi lsscsi sg3-utils multipath-tools scsitools
    ```

2. Enable multipathing:

    ```bash
    sudo tee /etc/multipath.conf <<-'EOF'
    defaults {
        user_friendly_names yes
        find_multipaths yes
    }
    EOF

    sudo service multipath-tools restart
    ```

3. Ensure that `iscsid` and `multipathd` are running:

    ```bash
    sudo service open-iscsi start
    sudo service multipath-tools start
    ```

4. Discover the iSCSI targets:

    ```bash
    sudo iscsiadm -m discoverydb -t st -p <DATA_LIF_IP> --discover
    ```

5. Login to the discovered iSCSI targets:

    ```bash
    sudo iscsiadm -m node -p <DATA_LIF_IP> --login
    ```

## Global Configuration File Variables

| Option            | Description                                                              | Example    |
| ----------------- | ------------------------------------------------------------------------ | ---------- |
| version           | Config file version number                                               | 1          |
| storageDriverName | `ontap-nas`, `ontap-san`, `eseries-iscsi`, or `solidfire-san`            | ontap-nas  |
| debug             | Turn debugging output on or off                                          | false      |
| storagePrefix     | Optional but recommended prefix for docker volume names.  Default: "netappdvp_"                 | netappdvp_ |

### Storage Prefix

A new config file variable has been added in v1.2 called "storagePrefix" that allows you to modify the prefix
applied to volume names by the plugin.  By default, when you run `docker volume create`, the volume name
supplied is prepended with "netappdvp_".  _("netappdvp-" for SolidFire.)_

If you wish to use a different prefix, you can specify it with this directive.  Alternatively, you can use
*pre-existing* volumes with the volume plugin by setting `storagePrefix` to an empty string, "".

*Ontap specific recommendation* Keeping the same storagePrefix across multiple docker hosts will show *all* docker volumes created with
the same storagePrefix, even if you use a different igroup. Use a different storagePrefix per docker host barring shared volumes. 
With iSCSI you are unlikely to ever want this scenario to happen.

*solidfire specific recommendation* do not use a storagePrefix (including the default)
By default the SolidFire driver will ignore this setting and not use a prefix.
We recommend using either a specific tenantID for docker volume mapping or
using the attribute data which is populated with the docker version, driver
info and raw name from docker in cases where any name munging may have been
used.

**A note of caution**: If you can see a volume with `docker volume ls` due to your storagePrefix name, then `docker volume rm <volumename>` will *delete* the specified volume, even if you aren't in the same igroup or if it was not created by the netappdvp plugin! Be very careful when using pre-existing volumes!

## ONTAP Specific Config File Variables

In addition to the global configuration values above, when using clustered Data ONTAP, these options are avaialble.

| Option            | Description                                                              | Example    |
| ----------------- | ------------------------------------------------------------------------ | ---------- |
| managementLIF     | IP address of clustered Data ONTAP management LIF                        | 10.0.0.1   |
| dataLIF           | IP address of protocol lif; will be derived if not specified             | 10.0.0.2   |
| svm               | Storage virtual machine to use (req, if management LIF is a cluster LIF) | svm_nfs    |
| username          | Username to connect to the storage device                                | vsadmin    |
| password          | Password to connect to the storage device                                | netapp123  |
| aggregate         | Aggregate to use for volume/LUN provisioning                             | aggr1      |
| IgroupName        | Optional but recommended prefix prefix for igroup name.  Default: "netappdvp"                 | netappdvp |

### Example ONTAP Config Files

**NFS Example for ontap-nas driver**

```json
{
    "version": 1,
    "storageDriverName": "ontap-nas",
    "managementLIF": "10.0.0.1",
    "dataLIF": "10.0.0.2",
    "svm": "svm_nfs",
    "username": "vsadmin",
    "password": "netapp123",
    "aggregate": "aggr1",
    "storagePrefix": "mydocker_nfs_"
}
```

**iSCSI Example for ontap-san driver**

```json
{
    "version": 1,
    "storageDriverName": "ontap-san",
    "managementLIF": "10.0.0.1",
    "dataLIF": "10.0.0.3",
    "svm": "svm_iscsi",
    "username": "vsadmin",
    "password": "netapp123",
    "aggregate": "aggr1",
    "IgroupName": "mydockerhost",
    "storagePrefix": "mydocker_iscsi_"
}
```

## E-Series Specific Config File Variables

In addition to the global configuration values above, when using E-Series, these options are available.

| Option                | Description                                                                             | Example       |
| --------------------- | --------------------------------------------------------------------------------------- | ------------- |
| webProxyHostname      | Hostname or IP address of Web Services Proxy                                            | localhost     |
| webProxyPort          | Port number of the Web Services Proxy (optional)                                        | 8443          |
| webProxyUseHTTP       | Use HTTP instead of HTTPS for Web Services Proxy (default = false)                      | true          |
| webProxyVerifyTLS     | Verify server's certificate chain and hostname (default = false)                        | true          |
| username              | Username for Web Services Proxy                                                         | rw            |
| password              | Password for Web Services Proxy                                                         | rw            |
| controllerA           | IP address of controller A                                                              | 10.0.0.5      |
| controllerB           | IP address of controller B                                                              | 10.0.0.6      |
| passwordArray         | Password for storage array if set                                                       | blank/empty   |
| hostDataIP            | Host iSCSI IP address (if multipathing just choose either one)                          | 10.0.0.101    |
| poolNameSearchPattern | Regular expression for matching storage pools available for nDVP volumes (default = .+) | docker.*      |
| hostType              | Type of E-series Host created by nDVP (default = linux_dm_mp)                           | linux_dm_mp   |
| accessGroupName       | Name of E-series Host Group to contain Hosts defined by nDVP (default = netappdvp)      | DockerHosts   |
 
### Example E-Series Config File

**Example for eseries-iscsi driver**

```json
{
	"version": 1,
	"storageDriverName": "eseries-iscsi",
	"debug": true,
	"webProxyHostname": "localhost",
	"webProxyPort": "8443",
	"webProxyUseHTTP": false,
	"webProxyVerifyTLS": true,
	"username": "rw",
	"password": "rw",
	"controllerA": "10.0.0.5",
	"controllerB": "10.0.0.6",
	"passwordArray": "",
	"hostDataIP": "10.0.0.101"
}
```

### E-Series Array Setup Notes

The E-Series Docker driver can provision Docker volumes in any storage pool on the array, including volume groups
and DDP pools. To limit the Docker driver to a subset of the storage pools, set the poolNameSearchPattern in the
configuration file to a regular expression that matches the desired pools.

When creating a docker volume you can specify the pool name, volume size, and disk media type using the '-o' option
and the tags 'pool, 'size', and 'mediaType', respectively. Valid values for media type are 'hdd' and 'ssd'. Note that
these are optional; if unspecified, the defaults will be a 1 GB volume allocated from an HDD pool. An example
of using these tags to create a 2 GB volume from any available SSD-based pool with sufficient space available:
 	
	docker volume create -d netapp --name my_vol -o size=2g -o mediaType=ssd

The E-series Docker driver will detect and use any preexisting Host definitions without modification, and
the driver will automatically define Host and Host Group objects as needed. The host type for hosts created
by the driver defaults to "linux_dm_mp", the native DM-MPIO multipath driver in Linux.

The current E-series Docker driver only supports iSCSI.

## SolidFire Specific Config File Variables

In addition to the global configuration values above, when using SolidFire, these options are available.

| Option            | Description                                                               | Example                    |
| ----------------- | --------------------------------------------------------------------------| -------------------------- |
| Endpoint          | Ex. https://<login>:<password>@<mvip>/json-rpc/<element-version>          |                            |
| SVIP              | iSCSI IP address and port                                                 | 10.0.0.7:3260              |
| TenantName        | SF Tenant to use (created if not found)                                   | "docker"                   |
| DefaultVolSz      | Volume size in GiB                                                        | 1                          |
| InitiatorIFace    | Specify interface when restricting iSCSI traffic to non-default interface | "default"                  |
| Types             | QoS specifications                                                        | See below                  |
| LegacyNamePrefix  | Prefix for upgraded NDVP installs                                         | "netappdvp-"               |


**LegacyNamePrefix** If you used a version of ndvp prior to 1.3.2 and perform an
upgrade with existing volumes, you'll need to set this value in order to access
your old volumes that were mapped via the volume-name method.

### Example Solidfire Config File

```json
{
    "version": 1,
    "storageDriverName": "solidfire-san",
    "debug": false,
    "Endpoint": "https://admin:admin@192.168.160.3/json-rpc/7.0",
    "SVIP": "10.0.0.7:3260",
    "TenantName": "docker",
    "DefaultVolSz": 1,
    "InitiatorIFace": "default",
    "Types": [
        {
            "Type": "Bronze",
            "Qos": {
                "minIOPS": 1000,
                "maxIOPS": 2000,
                "burstIOPS": 4000
            }
        },
        {
            "Type": "Silver",
            "Qos": {
                "minIOPS": 4000,
                "maxIOPS": 6000,
                "burstIOPS": 8000
            }
        },
        {
            "Type": "Gold",
            "Qos": {
                "minIOPS": 6000,
                "maxIOPS": 8000,
                "burstIOPS": 10000
            }
        }
    ]

}
```

##Known Issues and Limitations

1. Volume names must be a minimum of 2 characters in length

    This is a Docker client limitation. The client will interpret a single character name as being a Windows path. [Bug 25773]( https://github.com/docker/docker/issues/25773)
