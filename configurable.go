package configurable

import "sync"

type Configurable interface {
	CfChildren() []Configurable
}

var configurablesMutex sync.Mutex
var configurables []Configurable

func Register(configurable Configurable) {
	configurablesMutex.Lock()
	defer configurablesMutex.Unlock()

  if configurable == nil {
    panic("cannot register nil configurable")
  }

	configurables = append(configurables, configurable)
}

func Visit(do func(configurable Configurable) error) error {
	configurablesMutex.Lock()
	defer configurablesMutex.Unlock()

	for _, configurable := range configurables {
		err := do(configurable)
		if err != nil {
			return err
		}
	}

	return nil
}
