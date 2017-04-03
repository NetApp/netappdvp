.. _es_vol_opts:

E-Series Volume Options
=======================

Media Type
----------

The E-Series driver offers the ability to specify the type of disk which will be used to back the volume and,
like the other drivers, the ability to set the size of the volume at creation time.

Currently only two values for ``mediaType`` are supported:  ``ssd`` and ``hdd``.

.. code-block:: bash

   # create a 10GiB SSD backed volume
   docker volume create -d eseries --name eseriesSsd -o mediaType=ssd -o size=10G

   # create a 100GiB HDD backed volume
   docker volume create -d eseries --name eseriesHdd -o mediaType=hdd -o size=100G

Pool
----

The user can specify the pool name to use for creating the volume.  

``poolName`` is optional, if no pool is specified, then the default is to use all pools available.

.. code-block:: bash

  # create a volume using the "testme" pool
  docker volume create -d eseries --name testmePoolVolume -o poolName=testme -o size=100G

