Support
=======

Troubleshooting
---------------

The plugin has been built with a comprehensive logging capability that should help you diagnose most of the issues you are likely to come across. The method you use to access or tune those logs varies based on how you are running the plugin.

If you are running nDVP using the recommended managed plugin method (i.e., using ``docker plugin`` commands), the plugin is running in a container and the logs are available inside. Accessing those logs requires a little detective work because plugin containers are hidden from ``docker ps`` output:

  .. code-block:: bash

     # find the plugin container's abbreviated ID
     docker plugin ls
     
     # find the plugin container's full ID
     docker-runc list | grep <abbreviated ID>
     
     # view the logs in the container
     docker-runc exec -t <full ID> cat /var/log/netappdvp/netapp.log

The standard logging level should allow you to diagnose most issues. If you find that's not enough, you can enable debug logging:

  .. code-block:: bash

     # install the plugin with debug logging enabled
     docker plugin create netapp/ndvp-plugin:<version> --alias <alias> debug=true

     # or, enable debug logging on one that's already installed
     docker plugin disable <plugin>
     docker plugin set <plugin> debug=true
     docker plugin enable <plugin>

If you are not running as a managed plugin, the logs are available in the host's ``/var/log/netappdvp`` directory. If you need to enable debug logging, specify ``-debug`` when you run the plugin.

Getting Help
---------------

The nDVP is a supported NetApp product.  See the `find the support you need <http://mysupport.netapp.com/info/web/ECMLP2619434.html>`_ landing page on the Support site for options available to you.  To open a support case, use the serial number of the backend storage system and select containers and nDVP as the category you want help in.

There is also a vibrant discussion community of container users and engineers on the #containers channel of `Slack <http://netapp.io/slack>`_. This can be a great place to get answers and discuss with like-minded peers; highly recommended!



