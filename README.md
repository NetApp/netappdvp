# NetApp Docker Volume Plugin

The NetApp Docker Volume Plugin (nDVP) provides direct integration with the Docker ecosystem for the clustered Data ONTAP platform.  The nDVP package supports the provisioning and management of storage resources from the storage platform to Docker hosts, with a robust framework for adding additional platforms in the future.

Multiple instances of the nDVP can run concurrently on the same host.  The allows simultaneous connections to multiple storage systems and storage types, with the ablity to customize the storage used for the Docker volume(s).

- [Quick Start](#quick-start)
- [Running Multiple nDVP Instances](#running-multiple-ndvp-instances)
- [Configuring your Docker host for NFS or iSCSI](#configuring-your-docker-host-for-nfs-or-iscsi)
    - [NFS](#nfs)
        - [RHEL / CentOS](#rhel-centos)
        - [Ubuntu / Debian](#ubuntu-debian)
    - [iSCSI](#iscsi)
        - [RHEL / CentOS](#rhel-centos-1)
        - [Ubuntu / Debian](#ubuntu-debian-1)
    - [ONTAP Config File Variables](#ontap-config-file-variables)
        - [Example ONTAP Config Files](#example-ontap-config-files)

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
    wget https://github.com/netapp/?????/releases/download/v1.0/???.tar.gz
    tar xf ???.tar.gz

    # move to a location in the bin path
    sudo mv netappdvp/netappdvp /usr/local/bin
    sudo chown root:root /usr/local/bin/netappdvp
    sudo chmod 755 /usr/local/bin/netappdvp

    # create a location for the config files
    sudo mkdir -p /etc/netappdvp

    # create the configuration file, see below for Data ONTAP NFS and iSCSI configuration examples
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

1. Install the following system packages:

    ```bash
    sudo yum install -y lsscsi iscsi-initiator-utils sg3_utils device-mapper-multipath
    ```

2. Start the multipathing daemon:

    ```bash
    sudo mpathconf --enable --with_multipathd y
    ```

3. Ensure that `iscsid` and `multipathd` are enabled and running:

    ```bash
    sudo systemctl enable iscsid multipathd
    sudo systemctl start iscsid multipathd
    ```

4. Discover the iSCSI targets:

    ```bash
    sudo iscsiadm -m discoverydb -t st -p <DATA_LIF_IP> --discover
    ```

5. Login to the discovered iSCSI targets:

    ```bash
    sudo iscsiadm -m node -p <DATA_LIF_IP> --login
    ```

6. Start and enable `iscsi`:

    ```bash
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

## ONTAP Config File Variables

| Option            | Description                                                              | Example   |
| ----------------- | ------------------------------------------------------------------------ | --------- |
| version           | Config file version number                                               | 1         |
| storageDriverName | `ontap-nas` or `ontap-san`                                               | ontap-nas |
| managementLIF     | IP address of clustered Data ONTAP management LIF                        | 10.0.0.1  |
| dataLIF           | IP address of protocol lif; will be derived if not specified             | 10.0.0.2  |
| svm               | Storage virtual machine to use (req, if management LIF is a cluster LIF) | svm_nfs   |
| username          | Username to connect to the storage device                                | vsadmin   |
| password          | Password to connect to the storage device                                | netapp123 |
| aggregate         | Aggregate to use for volume/LUN provisioning                             | aggr1     |

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
    "aggregate": "aggr1"
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
    "aggregate": "aggr1"
}
```