package cflag

import "fmt"
import "strconv"
import "regexp"
import "strings"
import "gopkg.in/hlandau/configurable.v0"

// Group

type Registerable interface {
	Register(configurable configurable.Configurable)
}

type noReg struct{}

var NoReg noReg

func (r *noReg) Register(configurable configurable.Configurable) {

}

func register(r Registerable, c configurable.Configurable) {
	if r == nil {
		configurable.Register(c)
	} else {
		r.Register(c)
	}
}

type Group struct {
	configurables []configurable.Configurable
	name          string
}

func (ig *Group) CfName() string {
	return ig.name
}

func (ig *Group) CfChildren() []configurable.Configurable {
	return ig.configurables
}

func (ig *Group) String() string {
	return fmt.Sprintf("%s", ig.name)
}

func (ig *Group) Register(cfg configurable.Configurable) {
	ig.configurables = append(ig.configurables, cfg)
}

func NewGroup(reg Registerable, name string) *Group {
	ig := &Group{
		name: name,
	}
	register(reg, ig)
	return ig
}

// String

type SimpleFlag struct {
	name, curValue, summaryLine, defaultValue string
	curValuep                                 *string
}

func (sf *SimpleFlag) CfChildren() []configurable.Configurable {
	return nil
}

func (sf *SimpleFlag) String() string {
	return fmt.Sprintf("SimpleFlag(%s: %#v)", sf.name, *sf.curValuep)
}

func (sf *SimpleFlag) CfSetValue(v interface{}) error {
	vs, ok := v.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	*sf.curValuep = vs
	return nil
}

func (sf *SimpleFlag) CfValue() interface{} {
	return *sf.curValuep
}

func (sf *SimpleFlag) CfName() string {
	return sf.name
}

func (sf *SimpleFlag) CfUsageSummaryLine() string {
	return sf.summaryLine
}

func (sf *SimpleFlag) CfDefaultValue() interface{} {
	return sf.defaultValue
}

func (sf *SimpleFlag) Value() string {
	return *sf.curValuep
}

func NewSimpleFlag(reg Registerable, name, summaryLine, defaultValue string) *SimpleFlag {
	sf := &SimpleFlag{
		name:         name,
		summaryLine:  summaryLine,
		defaultValue: defaultValue,
		curValue:     defaultValue,
	}
	sf.curValuep = &sf.curValue

	register(reg, sf)
	return sf
}

// Int

type SimpleFlagInt struct {
	name, summaryLine      string
	curValue, defaultValue int
	curValuep              *int
}

func (sf *SimpleFlagInt) CfChildren() []configurable.Configurable {
	return nil
}

func (sf *SimpleFlagInt) String() string {
	return fmt.Sprintf("SimpleFlagInt(%s: %#v)", sf.name, *sf.curValuep)
}

func (sf *SimpleFlagInt) CfSetValue(v interface{}) error {
	vi, ok := v.(int)
	if ok {
		*sf.curValuep = vi
		return nil
	}

	vs, ok := v.(string)
	if ok {
		vs = strings.TrimSpace(vs)
		n, err := strconv.ParseInt(vs, 0, 32)
		if err != nil {
			return err
		}

		*sf.curValuep = int(n)
		return nil
	}

	return fmt.Errorf("invalid value for configurable %#v, expecting int: %v", sf.name, v)
}

func (sf *SimpleFlagInt) CfValue() interface{} {
	return sf.curValue
}

func (sf *SimpleFlagInt) CfName() string {
	return sf.name
}

func (sf *SimpleFlagInt) CfUsageSummaryLine() string {
	return sf.summaryLine
}

func (sf *SimpleFlagInt) CfDefaultValue() interface{} {
	return sf.defaultValue
}

func (sf *SimpleFlagInt) Value() int {
	return sf.curValue
}

func NewSimpleFlagInt(reg Registerable, name, summaryLine string, defaultValue int) *SimpleFlagInt {
	sf := &SimpleFlagInt{
		name:         name,
		summaryLine:  summaryLine,
		defaultValue: defaultValue,
		curValue:     defaultValue,
	}
	sf.curValuep = &sf.curValue

	register(reg, sf)
	return sf
}

// Bool

type SimpleFlagBool struct {
	name, summaryLine      string
	curValue, defaultValue bool
	curValuep              *bool
}

func (sf *SimpleFlagBool) CfChildren() []configurable.Configurable {
	return nil
}

func (sf *SimpleFlagBool) String() string {
	return fmt.Sprintf("SimpleFlagBool(%s: %#v)", sf.name, sf.curValue)
}

var re_no = regexp.MustCompilePOSIX(`^(00?|no?|f(alse)?)$`)

func (sf *SimpleFlagBool) CfSetValue(v interface{}) error {
	vb, ok := v.(bool)
	if ok {
		*sf.curValuep = vb
		return nil
	}

	vi, ok := v.(int)
	if ok {
		*sf.curValuep = (vi != 0)
		return nil
	}

	vs, ok := v.(string)
	if ok {
		vs = strings.TrimSpace(vs)
		*sf.curValuep = !re_no.MatchString(vs)
		return nil
	}

	return fmt.Errorf("invalid value for configurable %#v, expecting bool: %v", sf.name, v)
}

func (sf *SimpleFlagBool) CfValue() interface{} {
	return sf.curValue
}

func (sf *SimpleFlagBool) CfName() string {
	return sf.name
}

func (sf *SimpleFlagBool) CfUsageSummaryLine() string {
	return sf.summaryLine
}

func (sf *SimpleFlagBool) CfDefaultValue() interface{} {
	return sf.defaultValue
}

func (sf *SimpleFlagBool) Value() bool {
	return sf.curValue
}

func NewSimpleFlagBool(reg Registerable, name, summaryLine string, defaultValue bool) *SimpleFlagBool {
	sf := &SimpleFlagBool{
		name:         name,
		summaryLine:  summaryLine,
		defaultValue: defaultValue,
		curValue:     defaultValue,
	}
	sf.curValuep = &sf.curValue

	register(reg, sf)
	return sf
}
