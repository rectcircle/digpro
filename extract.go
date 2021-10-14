package digpro

import (
	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

type ExtractOption = internal.ExtractOption

// ExtractByName, for example
//   c := digpro.New()
//   _ = c.Supply(1, dig.Name("int"))  // please handle error in production
//   i, _ := c.Extract(int(0), digpro.ExtractByName("int"))
//   fmt.Println(i.(int) == 1) // true
func ExtractByName(name string) ExtractOption {
	return internal.ExtractOptionFunc(func(opts *internal.ExtractOptions) {
		opts.Name = name
	})
}

// ExtractByGroup, for example
//   c := digpro.New()
//   _ = c.Supply(1, dig.Group("ints"))  // please handle error in production
//   _ = c.Supply(1, dig.Group("ints"))
//   is, _ := c.Extract(int(0), digpro.ExtractByGroup("ints"))
//   fmt.Println(reflect.DeepEqual(is.([]int), []int{1, 1})) // true
func ExtractByGroup(name string) ExtractOption {
	return internal.ExtractOptionFunc(func(opts *internal.ExtractOptions) {
		opts.Group = name
	})
}

// MakeExtractFunc make Invoke function to extract a value and assign to *ptr from dig.Container, for example
//   c := dig.New()
//   _ = c.Provide(func() int { return 1 }) // please handle error in production
//   i := new(int)
//   _ = c.Invoke(digpro.MakeExtractFunc(i))
//   fmt.Println(*i == 1)
//   // Output: true
func MakeExtractFunc(ptr interface{}, opts ...ExtractOption) interface{} {
	return internal.MakeExtractFunc(ptr, opts...)
}

// Extract a value from dig.Container by type of value.
//
// if want to extract a non-interface, reflect.TypeOf(result) == reflect.TypeOf(typ). look like
//   func(int) -> int    // func(int(0)) -> int
//   func(*int) -> *int  // func(new(int)) -> *int
//
// if want to extract a interface, reflect.TypeOf(result) == reflect.TypeOf(typ).Elem(). look like
//   type A interface { ... }
//   func(A) -> error   // func(A(nil)) -> error
//   func(*A) -> A      // func(new(A)) -> A
//   func(**A) -> *A    // func(new(*A)) -> *A
//
// for example
//   c := dig.New()
//   _ = c.Provide(func() int {return 1})  // please handle error in production
//   i, _ := digpro.Extract(c, int(0))
//   fmt.Println(i.(int) == 1)
//   // Output: true
func Extract(c *dig.Container, typ interface{}, opts ...ExtractOption) (interface{}, error) {
	return internal.ExtractWithLocationForPC(c.Invoke, 3, typ, opts...)
}

// Extract a value from dig.Container by type of value.
//
// if want to extract a non-interface, reflect.TypeOf(result) == reflect.TypeOf(typ). look like
//   func(int) -> int    // func(int(0)) -> int
//   func(*int) -> *int  // func(new(int)) -> *int
//
// if want to extract a interface, reflect.TypeOf(result) == reflect.TypeOf(typ).Elem(). look like
//   type A interface { ... }
//   func(A) -> error   // func(A(nil)) -> error
//   func(*A) -> A      // func(new(A)) -> A
//   func(**A) -> *A    // func(new(*A)) -> *A
//
//for example
//   c := digpro.New()
//   _ = c.Supply(1)  // please handle error in production
//   i, _ := c.Extract(int(0))
//   fmt.Println(i.(int) == 1)
//   // Output: true
func (c *ContainerWrapper) Extract(typ interface{}, opts ...ExtractOption) (interface{}, error) {
	return internal.ExtractWithLocationForPC(c.Invoke, 3, typ, opts...)
}
