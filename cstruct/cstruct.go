package cstruct

import "fmt"
import "reflect"
import "strings"
import "regexp"
import "strconv"
import "gopkg.in/hlandau/configurable.v0"

type group struct {
	configurables []configurable.Configurable
}

func (g *group) CfChildren() []configurable.Configurable {
	return g.configurables
}

type value struct {
	name, usageSummaryLine string
	v                      reflect.Value
	defaultValue           interface{}
}

func (v *value) CfChildren() []configurable.Configurable {
	return nil
}

func (v *value) CfName() string {
	return v.name
}

func (v *value) String() string {
	return fmt.Sprintf("cstruct-value(%s)", v.CfName())
}

func (v *value) CfGetValue() interface{} {
	return v.v.Interface()
}

func (v *value) CfSetValue(x interface{}) error {
	xv := reflect.ValueOf(x)
	if !xv.Type().AssignableTo(v.v.Type()) {
		if xv.Type().Kind() != reflect.String {
			return fmt.Errorf("not assignable with that type")
		}

		pv, err := parseString(xv.String(), v.v.Type())
		if err != nil {
			return err
		}

		xv = reflect.ValueOf(pv)
		if !xv.Type().AssignableTo(v.v.Type()) {
			return fmt.Errorf("still not assignable with type after string conversion")
		}
	}

	v.v.Set(xv)
	return nil
}

func (v *value) CfDefaultValue() interface{} {
	return v.defaultValue
}

func (v *value) CfUsageSummaryLine() string {
	return v.usageSummaryLine
}

var re_no = regexp.MustCompilePOSIX(`^(00?|no?|f(alse)?)$`)

func parseString(s string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Int:
		n, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return nil, err
		}

		return int(n), nil

	case reflect.Bool:
		on := (s != "" && !re_no.MatchString(s))

		return on, nil

	default:
		return s, nil
	}
}

func MustNew(target interface{}) (c configurable.Configurable) {
	c, err := New(target)
	if err != nil {
		panic(err)
	}

	return c
}

func New(target interface{}) (c configurable.Configurable, err error) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = reflect.Indirect(v)
	}

	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("target interface is not a struct: %v", t)
		return
	}

	g := &group{}
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		name := strings.ToLower(field.Name)
		usage := field.Tag.Get("usage")
		dflt := field.Tag.Get("default")

		if usage == "" && dflt == "" {
			continue
		}

		vf := v.FieldByIndex(field.Index)

		var dfltv interface{}
		dfltv, err = parseString(dflt, vf.Type())
		if err != nil {
			err = fmt.Errorf("invalid default value: %#v", dflt)
			return
		}

		vv := &value{
			v:                vf,
			name:             name,
			usageSummaryLine: usage,
			defaultValue:     dfltv,
		}

		if !vf.CanSet() {
			err = fmt.Errorf("field not assignable")
			return
		}

		g.configurables = append(g.configurables, vv)

		// Do the type check now
		switch field.Type.Kind() {
		case reflect.Int:
		case reflect.String:
		case reflect.Bool:
		default:
			err = fmt.Errorf("unsupported field type: %v", field.Type)
			return
		}
	}

	return g, nil
}
