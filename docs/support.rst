Support
=======

Troubleshooting
---------------

The most common problem new users run into is a misconfiguration that prevents
the plugin from initializing. When this happens you will likely see a message
like this when you try to install or enable the plugin:

  .. code-block:: console

    Error response from daemon: dial unix /run/docker/plugins/<id>/netapp.sock: connect: no such file or directory

This simply means that the plugin failed to start. Luckily, the plugin has been
built with a comprehensive logging capability that should help you diagnose
most of the issues you are likely to come across.

The method you use to access or tune those logs varies based on how you are
running the plugin.

Managed plugin method
^^^^^^^^^^^^^^^^^^^^^

If you are running nDVP using the recommended managed plugin method (i.e.,
using ``docker plugin`` commands), the logs are passed through Docker itself,
so they will be interleaved with Docker's own logging.

To view them, simply run:

  .. code-block:: console

     # docker plugin ls
     ID                  NAME                DESCRIPTION                          ENABLED
     4fb97d2b956b        netapp:latest       nDVP - NetApp Docker Volume Plugin   false

     # journalctl -u docker | grep 4fb97d2b956b

The standard logging level should allow you to diagnose most issues. If you
find that's not enough, you can enable debug logging:

  .. code-block:: bash

     # install the plugin with debug logging enabled
     docker plugin install netapp/ndvp-plugin:<version> --alias <alias> debug=true

     # or, enable debug logging when the plugin is already installed
     docker plugin disable <plugin>
     docker plugin set <plugin> debug=true
     docker plugin enable <plugin>

Binary method
^^^^^^^^^^^^^

If you are not running as a managed plugin, you are running the binary itself
on the host. The logs are available in the host's ``/var/log/netappdvp``
directory. If you need to enable debug logging, specify ``-debug`` when you run
the plugin.

Getting Help
---------------

The nDVP is a supported NetApp project. See the
`find the support you need <http://mysupport.netapp.com/info/web/ECMLP2619434.html>`_
landing page on the Support site for options available to you. To open a
support case, use the serial number of the backend storage system and select
``containers`` and ``nDVP`` as the category you want help in. In most cases
the serial number you need is logged when the plugin first starts.

There is also a vibrant discussion community of container users and engineers
on the #containers channel in `NetApp's thePub Slack team <http://netapp.io/slack>`_.
This can be a great place to get answers and discuss with like-minded peers;
highly recommended!
