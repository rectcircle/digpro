package digpro

import (
	"io"

	"go.uber.org/dig"
)

// ContainerWrapper is a dig.Container wrapper, for add some method
type ContainerWrapper struct {
	dig.Container
}

// New constructs a dig.Container wrapper and export some metholds.
//
// For example.
//
//   c = digpro.New
//   // dig.Container methold
//   c.Provide(...)
//   c.Invoke(...)
//   // digpro exported methold
//   c.Value(...)
//   c.Struct(...)
//
func New(opts ...dig.Option) *ContainerWrapper {
	return &ContainerWrapper{
		Container: *dig.New(opts...),
	}
}

// Unwrap *ContainerWrapper to obtain *dig.Container
func (c *ContainerWrapper) Unwrap() *dig.Container {
	return &c.Container
}

// Visualize for write dot graph to io.Writer
func (c *ContainerWrapper) Visualize(w io.Writer, opts ...dig.VisualizeOption) error {
	return dig.Visualize(c.Unwrap(), w, opts...)
}
