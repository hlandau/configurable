package cflag_test

import "gopkg.in/hlandau/configurable.v0/cflag"
import "gopkg.in/hlandau/configurable.v0/adaptflag"
import flag "github.com/ogier/pflag"
import "fmt"

func Example() {
	var (
		g           = cflag.NewGroup(nil, "Program Options")
		bindFlag    = cflag.String(g, "bind", "Address to bind server to (e.g. :80)", ":80")
		fooFlag     = cflag.String(g, "foo", "Some flag", "")
		barFlag     = cflag.Int(g, "bar", "Some other flag", 42)
		doStuffFlag = cflag.Bool(g, "doStuff", "Do stuff?", false)
	)

	adaptflag.Adapt()
	flag.Parse()

	fmt.Printf("Bind: %s\n", bindFlag.Value())
	fmt.Printf("Foo:  %s\n", fooFlag.Value())
	fmt.Printf("Bar:  %d\n", barFlag.Value())
	fmt.Printf("Do Stuff: %v\n", doStuffFlag.Value())
}
