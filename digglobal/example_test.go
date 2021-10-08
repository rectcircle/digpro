package digglobal_test

import (
	"fmt"
	"os"
	"reflect"

	"github.com/rectcircle/digpro"
	"github.com/rectcircle/digpro/digglobal"
	"go.uber.org/dig"
)

type Foo struct {
	A       string
	B       int
	C       int  `name:"c"`
	private bool //lint:ignore U1000 for test
	ignore  int  `digpro:"ignore"`
}

func init() {
	// register object
	digglobal.Supply("a")
	digglobal.Supply(1)
	digglobal.Supply(2, dig.Name("c"))
	digglobal.Supply(true)

	// register a struct
	digglobal.Struct(Foo{ignore: 3})
}

func Example() {

	foo, err := digglobal.Extract(Foo{})
	if err != nil {
		digpro.QuickPanic(err)
	}
	c, err := digglobal.Extract(int(0), digpro.ExtractByName("c"))
	if err != nil {
		digpro.QuickPanic(err)
	}

	fmt.Println("*** foo ***")
	fmt.Printf("%#v\n", foo)
	fmt.Println("*** int[name=\"c\"] ***")
	fmt.Printf("%#v\n", c)
	fmt.Println("*** type of digglobal.Unwrap() ***")
	fmt.Println(reflect.TypeOf(digglobal.Unwrap()))
	fmt.Println("*** inspect node and value <see stderr> ***")
	os.Stderr.WriteString(digglobal.String())
	fmt.Println("*** inspect dot graph <see stderr> ***")
	digglobal.Visualize(os.Stderr)
	os.Stderr.WriteString("\n")
	// Output:
	// *** foo ***
	// digglobal_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
	// *** int[name="c"] ***
	// 2
	// *** type of digglobal.Unwrap() ***
	// *dig.Container
	// *** inspect node and value <see stderr> ***
	// *** inspect dot graph <see stderr> ***
}
