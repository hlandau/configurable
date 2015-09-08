package adaptenv

import "gopkg.in/hlandau/configurable.v0"
import "os"

// Loads values from environment variables into any configurables which expose
// CfEnvVarName() string. Priorities are checked.
func Adapt() {
	configurable.Visit(func(c configurable.Configurable) error {
		adaptRecursive(c)
		return nil
	})
}

func adaptRecursive(c configurable.Configurable) {
	cc, ok := c.(interface {
		CfChildren() []configurable.Configurable
	})
	if ok {
		for _, ch := range cc.CfChildren() {
			adaptRecursive(ch)
		}
	}

	adapt(c)
}

func adapt(c configurable.Configurable) {
	cenv, ok := c.(interface {
		CfEnvVarName() string
		CfSetValue(x interface{}) error
	})
	if !ok {
		return
	}

	envVarName := cenv.CfEnvVarName()
	if envVarName == "" {
		return
	}

	v, ok := os.LookupEnv(envVarName)
	if !ok {
		return
	}

	cprio, ok := c.(interface {
		CfGetPriority() configurable.Priority
		CfSetPriority(priority configurable.Priority)
	})
	if ok {
		if cprio.CfGetPriority() > configurable.EnvPriority {
			return
		}
	}

	err := cenv.CfSetValue(v)
	if err != nil {
		return
	}

	if ok {
		cprio.CfSetPriority(configurable.EnvPriority)
	}
}
