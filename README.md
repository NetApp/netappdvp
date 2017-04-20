# NetApp Docker Volume Plugin

The NetApp Docker Volume Plugin (nDVP) provides direct integration with the Docker ecosystem for NetApp's ONTAP, E-Series, and SolidFire storage platforms. The nDVP package supports the provisioning and management of storage resources from the storage platform to Docker hosts, with a robust framework for adding additional platforms in the future.

Multiple instances of nDVP can run concurrently on the same host.  The allows simultaneous connections to multiple storage systems and storage types, with the ablity to customize the storage used for the Docker volume(s).

Documentation for nDVP can be found on [Read the Docs](http://netappdvp.readthedocs.io/en/latest/).  It can also be deployed directly from the [Docker Store](https://store.docker.com/plugins/netapp-docker-volume-plugin-ndvp).
