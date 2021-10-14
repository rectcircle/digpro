package digpro_test

import (
	"fmt"

	"github.com/rectcircle/digpro"
	"go.uber.org/dig"
)

func ExampleExtract() {
	c := dig.New()
	_ = c.Provide(func() int { return 1 }) // please handle error in production
	i, _ := digpro.Extract(c, int(0))
	fmt.Println(i.(int) == 1)
	// Output: true
}

func ExampleContainerWrapper_Extract() {
	c := digpro.New()
	_ = c.Supply(1) // please handle error in production
	i, _ := c.Extract(int(0))
	fmt.Println(i.(int) == 1)
	// Output: true
}

func ExampleMakeExtractFunc() {
	c := dig.New()
	_ = c.Provide(func() int { return 1 }) // please handle error in production
	i := new(int)
	_ = c.Invoke(digpro.MakeExtractFunc(i))
	fmt.Println(*i == 1)
	// Output: true
}
