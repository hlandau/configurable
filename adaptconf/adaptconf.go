// Package adaptconf adapts registered configurables to configuration file
// formats.
package adaptconf

import "os"
import "fmt"
import "strings"
import "path/filepath"
import "gopkg.in/hlandau/configurable.v0"
import "gopkg.in/hlandau/configurable.v0/cflag"
import "gopkg.in/hlandau/service.v1/exepath"
import "github.com/BurntSushi/toml"

var confFlag = cflag.String(nil, "conf", "Configuration file path", "")
var lastConfPath string

func LastConfPath() string {
	return lastConfPath
}

func LoadPath(confFilePath string) error {
	var m map[string]interface{}
	_, err := toml.DecodeFile(confFilePath, &m)
	if err != nil {
		return err
	}

	lastConfPath = confFilePath

	configurable.Visit(func(c configurable.Configurable) error {
		applyChild(c, m)
		return nil
	})

	return nil
}

func LoadPaths(paths []string) error {
	confPath := confFlag.Value()

	if confPath == "" {
		for _, path := range paths {
			path = expandPath(path)

			if !pathExists(path) {
				continue
			}

			confPath = path
		}
	}

	if confPath == "" {
		return nil
	}

	return LoadPath(confPath)
}

func Load(programName string) error {
	return LoadPaths([]string{
		fmt.Sprintf("/etc/%s/%s.conf", programName, programName),
		fmt.Sprintf("/etc/%s.conf", programName),
		fmt.Sprintf("etc/%s.conf", programName),
		fmt.Sprintf("$BIN/%s.conf", programName),
		fmt.Sprintf("$BIN/../etc/%s/%s.conf", programName, programName),
		fmt.Sprintf("$BIN/../etc/%s.conf", programName),
	})
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return false
}

func expandPath(path string) string {
	if !strings.HasPrefix(path, "$BIN/") {
		return path
	}

	return filepath.Join(filepath.Dir(exepath.Abs), path[5:])
}

func apply(c configurable.Configurable, v interface{}) error {
	cch, ok := c.(interface {
		CfChildren() []configurable.Configurable
	})
	if ok {
		children := cch.CfChildren()
		if len(children) > 0 {
			return applyChildren(children, v)
		}
	}

	csv, ok := c.(interface {
		CfSetValue(x interface{}) error
	})
	if !ok {
		return nil
	}

	cprio, ok := c.(interface {
		CfSetPriority(priority configurable.Priority)
		CfGetPriority() configurable.Priority
	})
	if ok {
		prio := cprio.CfGetPriority()
		if prio <= configurable.ConfigPriority {
			err := csv.CfSetValue(v)
			if err != nil {
				return nil
			}

			cprio.CfSetPriority(configurable.ConfigPriority)
		}

		return nil
	} else {
		return csv.CfSetValue(v)
	}
}

func applyChildren(chs []configurable.Configurable, v interface{}) error {
	vm, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	for _, ch := range chs {
		applyChild(ch, vm)
	}

	return nil
}

func applyChild(ch configurable.Configurable, vm map[string]interface{}) error {
	name, ok := name(ch)
	if !ok {
		return nil
	}

	vch, ok := vm[name]
	if !ok {
		return nil
	}

	return apply(ch, vch)
}

func name(c configurable.Configurable) (name string, ok bool) {
	v, ok := c.(interface {
		CfName() string
	})
	if !ok {
		return
	}

	return v.CfName(), true
}
