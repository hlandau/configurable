package adaptflag

import "fmt"
import "flag"
import "github.com/ogier/pflag"
import "gopkg.in/hlandau/configurable.v0"

func name(c configurable.Configurable) (name string, ok bool) {
	v, ok := c.(interface {
		CfName() string
	})
	if !ok {
		return
	}

	return v.CfName(), true
}

func usageSummaryLine(c configurable.Configurable) (s string, ok bool) {
	v, ok := c.(interface {
		CfUsageSummaryLine() string
	})
	if !ok {
		return
	}

	return v.CfUsageSummaryLine(), true
}

var errNotSupported = fmt.Errorf("not supported")

type value struct {
	c configurable.Configurable
}

// The flag package uses this to get the default value.
func (v *value) String() string {
	cs, ok := v.c.(interface {
		CfDefaultValue() interface{}
	})
	if !ok {
		return "[configurable]"
	}

	return fmt.Sprintf("%#v", cs.CfDefaultValue())
}

func (v *value) Set(s string) error {
	cs, ok := v.c.(interface {
		CfSetValue(v interface{}) error
	})
	if !ok {
		return errNotSupported
	}

	return cs.CfSetValue(s)
}

func (v *value) Get() interface{} {
	cg, ok := v.c.(interface {
		CfGetValue() interface{}
	})
	if !ok {
		return nil // ...
	}

	return cg.CfGetValue()
}

var adapted = map[interface{}]struct{}{}

func adapt(c configurable.Configurable, f AdaptFunc) error {
	_, ok := adapted[c]
	if ok {
		return nil
	}

	name, ok := name(c)
	if !ok {
		return errNotSupported
	}

	_, ok = c.(interface {
		CfSetValue(v interface{}) error
	})
	if !ok {
		return errNotSupported
	}

	v := &value{c: c}
	usage, _ := usageSummaryLine(c)

	f(v, name, usage)

	adapted[c] = struct{}{}
	return nil
}

type AdaptFunc func(v Value, name, usage string)

func recursiveAdapt(c configurable.Configurable, f AdaptFunc) error {
	adapt(c, f)
	for _, ch := range c.CfChildren() {
		recursiveAdapt(ch, f)
	}
	return nil
}

type Value interface {
	String() string
	Set(x string) error
	//Get() interface{}
}

func AdaptWithFunc(f AdaptFunc) {
	configurable.Visit(func(c configurable.Configurable) error {
		return recursiveAdapt(c, f)
	})
}

func Adapt() {
	AdaptWithFunc(func(v Value, name, usage string) {
		flag.Var(v, name, usage)
		pflag.Var(v, name, usage)
	})
}
