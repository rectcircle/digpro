package digpro_test

import (
	"fmt"

	"github.com/rectcircle/digpro"
)

func ExampleOverride() {
	c := digpro.New()
	_ = c.Supply(1) // please handle error in production
	_ = c.Supply(1, digpro.Override())
	// _ = c.Supply("a", digpro.Override())  // has error
	i, _ := c.Extract(0)
	fmt.Println(i.(int) == 1)
	// Output: true
}
