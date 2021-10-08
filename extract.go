package digpro

import (
	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

type ExtractOption internal.ExtractOption
type extractOptionFunc func(*internal.ExtractOptions)

func (f extractOptionFunc) ApplyExtractOption(opts *internal.ExtractOptions) { f(opts) }

// ExtractByName, for example
//   c := digpro.New()
//   _ = c.Supply(1, dig.Name("int"))  // please handle error in production
//   i, _ := c.Extract(int(0), digpro.ExtractByName("int"))
//   fmt.Println(i.(int) == 1) // true
func ExtractByName(name string) ExtractOption {
	return extractOptionFunc(func(opts *internal.ExtractOptions) {
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
	return extractOptionFunc(func(opts *internal.ExtractOptions) {
		opts.Group = name
	})
}

func toInternalExtractOption(opts []ExtractOption) []internal.ExtractOption {
	result := make([]internal.ExtractOption, 0, len(opts))
	for _, opt := range opts {
		result = append(result, internal.ExtractOption(opt))
	}
	return result
}

// Extract a value from dig.Container by type of value, for example
//   c := dig.New()
//   _ = c.Provide(func() int {return 1})  // please handle error in production
//   i, _ := digpro.Extract(c, int(0))
//   fmt.Println(i.(int) == 1)
//   // Output: true
func Extract(c *dig.Container, typInterface interface{}, opts ...ExtractOption) (interface{}, error) {
	return internal.ExtractWithLocationForPC(c, 2, typInterface, toInternalExtractOption(opts)...)
}

// Extract a value from dig.Container by type of value, for example
//   c := digpro.New()
//   _ = c.Supply(1)  // please handle error in production
//   i, _ := c.Extract(int(0))
//   fmt.Println(i.(int) == 1)
//   // Output: true
func (c *ContainerWrapper) Extract(typInterface interface{}, opts ...ExtractOption) (interface{}, error) {
	return internal.ExtractWithLocationForPC(&c.Container, 2, typInterface, toInternalExtractOption(opts)...)
}
