Configurable: The useless Go configuration package that doesn't do anything
===========================================================================

[![GoDoc](https://godoc.org/gopkg.in/hlandau/configurable.v1?status.svg)](https://godoc.org/gopkg.in/hlandau/configurable.v1)

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

**Import as:** `gopkg.in/hlandau/configurable.v1`

Configurable
------------

A Configurable is an object that represents some configurable thing. It is
obliged only to implement the following interface:

```go
type Configurable interface{}
```

Configurable is designed around interface upgrades. If you want to actually do
anything with a Configurable, you must attempt to cast it to an interface with
the methods you need. A Configurable is not obliged to implement any interface
besides Configurable, but almost always will.

Here are some common interfaces implemented by Configurables, in descending
order of importance:

  - `CfSetValue(interface{}) error` — attempt to set the Configurable to a value.

  - `CfName() string` — get the Configurable's name.

  - `CfDefaultValue() interface{}` — get the Configurable's default value.

  - `CfGetValue() interface{}` — get the Configurable's value.

  - `CfChildren() []Configurable` — return the children of this Configurable, if any.

  - `CfUsageSummaryLine() string` — get a one-line usage summary suitable for
    use as  command line usage information.

  - `String() string` — the standard Go `String()` interface.

  - `CfGetPriority() Priority` — retrieves the priority of the value, used to
    determine whether it should be overridden.

  - `CfSetPriority(priority Priority)` — sets the priority of the value.

  - `CfEnvVarName() string` — if a non-empty string, an environment variable
    that maps to this Configurable.

Configurable-specific methods should always be prefixed with `Cf` so that it is clear
that they are intended for consumption by Configurable consumers.

A command line parsing adapter should typically be able to make do with a Configurable
which implements just `CfSetValue` and `CfName`.

The Standard Bindings
---------------------

For a package which makes it easy to register and consume configurables, see
the [easyconfig](https://github.com/hlandau/easyconfig) package.

Of course, nothing requires you to use the easyconfig package. You are free to
eschew it and make your own.

Background Reading
------------------

  - [On Nexuses](http://www.devever.net/~hl/nexuses)
  - See also: [Measurable](https://github.com/hlandau/measurable)

Licence
-------

    © 2015 Hugo Landau <hlandau@devever.net>    MIT License

