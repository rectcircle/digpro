package digglobal

import (
	"io"

	"github.com/rectcircle/digpro"
	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

var g = digpro.New()

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// Provide see https://pkg.go.dev/go.uber.org/dig#Container.Provide
//
// Note: if has error will panic
func Provide(constructor interface{}, opts ...dig.ProvideOption) {
	digpro.QuickPanic(g.Provide(constructor, opts...))
}

// Invoke see https://pkg.go.dev/go.uber.org/dig#Container.Invoke
func Invoke(function interface{}, opts ...dig.InvokeOption) error {
	return g.Invoke(function, opts...)
}

// String see https://pkg.go.dev/go.uber.org/dig#Container.String
func String() string {
	return g.String()
}

// Supply see https://pkg.go.dev/github.com/rectcircle/digpro#ContainerWrapper.Supply
//
// Note: if has error will panic
func Supply(value interface{}, opts ...dig.ProvideOption) {
	panicIfError(internal.ProvideWithLocationForPC(&g.Container, 2, digpro.Supply(value), opts...))
}

// Struct see https://pkg.go.dev/github.com/rectcircle/digpro#ContainerWrapper.Struct
//
// Note: if has error will panic
func Struct(structOrStructPtr interface{}, opts ...dig.ProvideOption) {
	panicIfError(internal.ProvideWithLocationForPC(&g.Container, 2, digpro.Struct(structOrStructPtr), opts...))
}

// Extract see https://pkg.go.dev/github.com/rectcircle/digpro#ContainerWrapper.Extract
func Extract(typInterface interface{}, opts ...digpro.ExtractOption) (interface{}, error) {
	return internal.ExtractWithLocationForPC(&g.Container, 2, typInterface, toInternalExtractOption(opts)...)
}

// Unwrap see https://pkg.go.dev/github.com/rectcircle/digpro#ContainerWrapper.Unwrap
func Unwrap() *dig.Container {
	return g.Unwrap()
}

// Visualize see https://pkg.go.dev/github.com/rectcircle/digpro#ContainerWrapper.Visualize
func Visualize(w io.Writer, opts ...dig.VisualizeOption) error {
	return g.Visualize(w, opts...)
}

// toInternalExtractOption
// TODO Duplicate code block
func toInternalExtractOption(opts []digpro.ExtractOption) []internal.ExtractOption {
	result := make([]internal.ExtractOption, 0, len(opts))
	for _, opt := range opts {
		result = append(result, internal.ExtractOption(opt))
	}
	return result
}
