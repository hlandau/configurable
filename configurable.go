// Package configurable provides an integration nexus for program and library
// configuration items.
//
// Configurable is a Go library for managing program configuration information,
// no matter whether it comes from command line arguments, configuration files,
// environment variables, or anywhere else.
//
// The most noteworthy feature of configurable is that it doesn't do anything.
// It contains no functionality for examining or parsing command line
// arguments. It doesn't do anything with environment variables. And it
// certainly can't read configuration files.
//
// The purpose of configurable is to act as an [integration
// nexus](http://www.devever.net/~hl/nexuses), essentially a matchmaker between
// application configuration and specific configuration interfaces. This
// creates the important feature that your application's configuration can be
// expressed completely independently of *how* that configuration is loaded.
//
// Configurable doesn't implement any configuration loading logic because it
// strives to be a neutral intermediary, which abstracts the interface between
// configurable items and configurators.
//
// Pursuant to this, package configurable is this and only this: an interface
// Configurable which all configuration items must implement, and a facility
// for registering top-level Configurables and visiting them.
//
// In v1, the Configurable interface has no methods and is thus considered to
// be implemented by anything.
package configurable // import "gopkg.in/hlandau/configurable.v1"

import "sync"

// Configurable is the interface which must be implemented by any configuration
// item to be used with package configurable. In the current version, v1, it
// contains no methods and is thus satisfied by anything. All functionality
// must be obtained via interface upgrades.
type Configurable interface{}

var configurablesMutex sync.RWMutex
var configurables []Configurable

// Registers a top-level Configurable.
func Register(configurable Configurable) {
	configurablesMutex.Lock()
	defer configurablesMutex.Unlock()

	if configurable == nil {
		panic("cannot register nil configurable")
	}

	configurables = append(configurables, configurable)
}

// Visits all registered top-level Configurables.
//
// Returning a non-nil error short-circuits the iteration process and returns
// that error.
func Visit(do func(configurable Configurable) error) error {
	configurablesMutex.RLock()
	defer configurablesMutex.RUnlock()

	for _, configurable := range configurables {
		err := do(configurable)
		if err != nil {
			return err
		}
	}

	return nil
}

// Priority values are used to determine whether values should be overridden.
type Priority int

const (
	// The priority of default-set values.
	DefaultPriority Priority = 0

	// The recommended priority for values set from environment variables.
	EnvPriority = 1000

	// The recommended priority for values loaded from a config file.
	ConfigPriority = 2000

	// The recommended priority for values set from command line flags.
	FlagPriority = 3000
)
