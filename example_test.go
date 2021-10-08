package digpro_test

import (
	"fmt"

	"github.com/rectcircle/digpro"
	"go.uber.org/dig"
)

func Example_lowerLevelAPI() {
	// Lower Level API can be used quickly, but the error is not friendly enough, such as:
	//   ... missing dependencies for function "reflect".makeFuncStub (/usr/local/Cellar/go/1.17.1/libexec/src/reflect/asm_amd64.s:30) ...

	type Foo struct {
		A       string
		B       int
		C       int  `name:"c"`
		private bool //lint:ignore U1000 for test
		ignore  int  `digpro:"ignore"`
	}
	// new a *dig.Container
	c := dig.New()
	// provide some constructor
	digpro.QuickPanic(
		// register object
		c.Provide(digpro.Supply("a")),
		c.Provide(digpro.Supply(1)),
		c.Provide(digpro.Supply(true)),
		c.Provide(digpro.Supply(2), dig.Name("c")),
		// register a struct
		c.Provide(digpro.Struct(Foo{
			ignore: 3,
		})),
	)
	// extract object from container
	foo, err := digpro.Extract(c, Foo{})
	if err != nil {
		digpro.QuickPanic(err)
	}
	fmt.Printf("%#v", foo)
	// Output: digpro_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
}

func Example_highLevelAPI() {
	// High Level API has better error output

	type Foo struct {
		A       string
		B       int
		C       int  `name:"c"`
		private bool //lint:ignore U1000 for test
		ignore  int  `digpro:"ignore"`
	}
	// new a *dig.Container wrapper - *digpro.ContainerWrapper
	c := digpro.New()
	// provide some constructor
	digpro.QuickPanic(
		// register object
		c.Supply("a"),
		c.Supply(1),
		c.Supply(2, dig.Name("c")),
		c.Supply(true),
		// register a struct
		c.Struct(Foo{
			ignore: 3,
		}),
	)
	// extract object from container
	foo, err := c.Extract(Foo{})
	if err != nil {
		digpro.QuickPanic(err)
	}
	fmt.Printf("%#v", foo)
	// Output: digpro_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
}
