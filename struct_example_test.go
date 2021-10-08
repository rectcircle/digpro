package digpro_test

import (
	"fmt"

	"github.com/rectcircle/digpro"
	"go.uber.org/dig"
)

func ExampleStruct() {
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
		// equals to
		// c.Provide(func(in struct {
		// 	A       string
		// 	B       int
		// 	C       int `name:"c"`
		// 	Private bool
		// }) Foo {
		// 	return Foo{
		// 		A:       in.A,
		// 		B:       in.B,
		// 		C:       in.C,
		// 		private: in.Private,
		// 		ignore:  3,
		// 	}
		// }),
	)
	// extract object from container
	foo, err := digpro.Extract(c, Foo{})
	if err != nil {
		digpro.QuickPanic(err)
	}
	fmt.Printf("%#v", foo)
	// Output: digpro_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
}

func ExampleContainerWrapper_Struct() {
	type Foo struct {
		A       string
		B       int
		C       int  `name:"c"`
		private bool //lint:ignore U1000 for test
		ignore  int  `digpro:"ignore"`
	}
	// new a *dig.Container
	c := digpro.New()
	// provide some constructor
	digpro.QuickPanic(
		// register object
		c.Supply("a"),
		c.Supply(1),
		c.Supply(true),
		c.Supply(2, dig.Name("c")),
		// register a struct
		c.Struct(Foo{
			ignore: 3,
		}),
		// equals to
		// c.Provide(func(in struct {
		// 	A       string
		// 	B       int
		// 	C       int `name:"c"`
		// 	Private bool
		// }) Foo {
		// 	return Foo{
		// 		A:       in.A,
		// 		B:       in.B,
		// 		C:       in.C,
		// 		private: in.Private,
		// 		ignore:  3,
		// 	}
		// }),
	)
	// extract object from container
	foo, err := c.Extract(Foo{})
	if err != nil {
		digpro.QuickPanic(err)
	}
	fmt.Printf("%#v", foo)
	// Output: digpro_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
}
