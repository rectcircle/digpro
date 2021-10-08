package digpro_test

import (
	"fmt"

	"github.com/rectcircle/digpro"
	"go.uber.org/dig"
)

func ExampleSupply() {
	c := dig.New()
	digpro.QuickPanic(
		// register object
		c.Provide(digpro.Supply("a")),
		// equals to
		// c.Provide(func() string {return "a"}),
	)
	foo, err := digpro.Extract(c, string(""))
	if err != nil {
		digpro.QuickPanic(err)
	}
	fmt.Println(foo)
	// Output: a
}

func ExampleContainerWrapper_Supply() {
	c := digpro.New()
	digpro.QuickPanic(
		// register object
		c.Supply("a"),
		// equals to
		// c.Provide(func() string {return "a"}),
	)
	foo, err := c.Extract(string(""))
	if err != nil {
		digpro.QuickPanic(err)
	}
	fmt.Println(foo)
	// Output: a
}
