package digpro

import (
	"reflect"

	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

// Supply a value into container.
// for example
//   c := dig.New()
//   digpro.QuickPanic(
//   	// register object
//   	c.Provide(digpro.Supply("a")),
//   	// equals to
//   	// c.Provide(func() string {return "a"}),
//   )
//   foo, err := digpro.Extract(c, string(""))
//   if err != nil {
//   	digpro.QuickPanic(err)
//   }
//   fmt.Println(foo)
//   // Output: a
func Supply(value interface{}) interface{} {
	switch v := value.(type) {
	case nil:
		return func() interface{} { return v }
	case error:
		return func() error { return v }
	}

	typ := reflect.TypeOf(value)
	returnTypes := []reflect.Type{typ}
	returnValues := []reflect.Value{reflect.ValueOf(value)}

	ft := reflect.FuncOf([]reflect.Type{}, returnTypes, false)
	fv := reflect.MakeFunc(ft, func([]reflect.Value) []reflect.Value {
		return returnValues
	})
	return fv.Interface()
}

// Supply a value into container.
// for example
//   c := digpro.New()
//   digpro.QuickPanic(
//   	// register object
//   	c.Supply("a"),
//   	// equals to
//   	// c.Provide(func() string {return "a"}),
//   )
//   foo, err := c.Extract(string(""))
//   if err != nil {
//   	digpro.QuickPanic(err)
//   }
//   fmt.Println(foo)
//   // Output: a
func (c *ContainerWrapper) Supply(value interface{}, opts ...dig.ProvideOption) error {
	filteredOpts, digproOptsResult := filterAndGetDigproProvideOptions(opts, locationFixOptionType)
	callSkip := digproOptsResult.locationFixCallSkip
	if callSkip == 0 {
		callSkip = 3
	}
	return internal.ProvideWithLocationForPC(c.Provide, callSkip, Supply(value), filteredOpts...)
}
