Configurable: The useless Go configuration package that doesn't do anything
===========================================================================

[GoDoc](http://godoc.org/github.com/hlandau/configurable)

Configurable is a Go library for managing program configuration information, no
matter whether it comes from command line arguments, configuration files,
environment variables, or anywhere else.

The most noteworthy feature of configurable is that it doesn't do anything. It
contains no functionality for examining or parsing command line arguments. It
doesn't do anything with environment variables. And it certainly can't read
configuration files.

The purpose of configurable is to act as an [integration
nexus](http://www.devever.net/~hl/nexuses), essentially a matchmaker between
application configuration and specific configuration interfaces. This creates
the important feature that your application's configuration can be expressed
completely independently of *how* that configuration is loaded.

Configurable doesn't implement any configuration loading logic because it
strives to be a neutral intermediary, which abstracts the interface between
configurable items and configurators.

In order to demonstrate the configurable way of doing things, a simple flag
parsing package is included. Use of this package is completely optional. If it
doesn't meet your needs, you can throw it out and use your own — but still
consume and configure all registered Configurables.

Included example packages demonstrate how an application or library might
register various configurable items, and then expose them for configuration via
the command line, configuration files or other means.

Background Reading
------------------

  - [On Nexuses](http://www.devever.net/~hl/nexuses)

Licence
-------

    © 2015 Hugo Landau <hlandau@devever.net>    MIT License

