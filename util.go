package digpro

import (
	"fmt"
	"strings"
)

// QuickPanic if anyone of errs is not nil, will panic.
// for example
//   c := digpro.New()
//   digpro.QuickPanic(
//   	c.Supply(1),
//   	c.Supply(1),
//   )
//   // panic: [1]: cannot provide function "xxx".Xxx (xxx.go:n): cannot provide int from [0]: already provided by "xxx".Xxx (xxx.go:m)
func QuickPanic(errs ...error) {
	msgs := []string{}
	for i, err := range errs {
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("[%d]: %s", i, err.Error()))
		}
	}
	if len(msgs) != 0 {
		panic(strings.Join(msgs, "\n"))
	}
}
