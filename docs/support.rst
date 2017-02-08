Support
=======

Support for the NetApp Docker Volume Plugin is provided on a best effort basis through the `NetApp Community <http://community.netapp.com>`_ and using `Slack <http://netapp.io/slack>`_.

Troubleshooting
---------------

* Ensure debug logging is turned on (it is enabled by default) for the ndvp daemon instance using the ``-debug=true`` command line option when starting.
* For the traditional install (Docker <= 1.12), check the logs at ``/var/log/netappdvp`` on the host.
* To retrive the logs from the container for the managed plugin (Docker >= 1.13 / 17.03):
  
  .. code-block:: bash
     
     # find the container
     docker-runc list
     
     # view the logs in the container
     docker-runc exec -t <container id> cat /var/log/netappdvp/netapp.log