package cstruct_test

import "gopkg.in/hlandau/configurable.v0"
import "gopkg.in/hlandau/configurable.v0/cstruct"
import "gopkg.in/hlandau/configurable.v0/adaptflag"
import flag "github.com/ogier/pflag"
import "fmt"

func Example() {
	type Config struct {
		Bind    string `usage:"Address to bind server to (e.g. :80)" default:":80"`
		Foo     string `usage:"Some flag" default:""`
		Bar     int    `usage:"Some other flag" default:"42"`
		DoStuff bool   `usage:"Do stuff?" default:"false"`
	}

	cfg := &Config{}
	configurable.Register(cstruct.MustNew(cfg))
	adaptflag.Adapt()
	flag.Parse()

	fmt.Printf("Bind: %s\n", cfg.Bind)
	fmt.Printf("Foo:  %s\n", cfg.Foo)
	fmt.Printf("Bar:  %d\n", cfg.Bar)
	fmt.Printf("Do Stuff: %v\n", cfg.DoStuff)
}
