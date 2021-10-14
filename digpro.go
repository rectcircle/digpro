package digpro

import (
	"io"

	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

// ContainerWrapper is a dig.Container wrapper, for add some method
type ContainerWrapper struct {
	dig.Container
	middlewares  []provideMiddleware
	provideInfos []internal.ProvideInfoWrapper
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
		middlewares: []provideMiddleware{
			overrideProvideMiddleware,
		},
	}
}

// Unwrap *ContainerWrapper to obtain *dig.Container.
//
// WARNING: the methold only for debug, please not use in production
func (c *ContainerWrapper) Unwrap() *dig.Container {
	return &c.Container
}

// Visualize for write dot graph to io.Writer
func (c *ContainerWrapper) Visualize(w io.Writer, opts ...dig.VisualizeOption) error {
	return dig.Visualize(c.Unwrap(), w, opts...)
}

// Provide teaches the container how to build values of one or more types and expresses their dependencies.
// more see: https://pkg.go.dev/go.uber.org/dig#Container.Provide.
//
// digpro.ContainerWrapper.Provide() support digpro.Override() options, but dig.Container.Provide() not support
func (c *ContainerWrapper) Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	return newProvideContext(c, constructor, opts).next()
}

type provideContext struct {
	c           *ContainerWrapper
	constructor interface{}
	opts        []dig.ProvideOption
	index       int
}

func newProvideContext(c *ContainerWrapper, constructor interface{}, opts []dig.ProvideOption) *provideContext {
	return &provideContext{
		c:           c,
		constructor: constructor,
		opts:        opts,
		index:       -1,
	}
}

func (pc *provideContext) next() error {
	pc.index += 1
	nowIndex := pc.index
	if pc.index < len(pc.c.middlewares) {
		err := pc.c.middlewares[pc.index](pc)
		if err != nil {
			return err
		}
		if nowIndex == pc.index {
			return pc.next()
		}
		return err
	}
	return pc.doProvide()
}

func (pc *provideContext) doProvide() error {
	internalOpts := internal.ApplyProvideOptions(pc.opts...)
	info := internal.ProvideInfoWrapper{}
	if internalOpts.Info == nil {
		pc.opts = append(pc.opts, dig.FillProvideInfo(&info.ProvideInfo))
	}
	err := pc.c.Container.Provide(pc.constructor, pc.opts...)
	if err != nil {
		return err
	}
	if internalOpts.Info != nil {
		info.ProvideInfo = *internalOpts.Info
	}
	pc.c.provideInfos = append(pc.c.provideInfos, info)
	return nil
}

type provideMiddleware func(pc *provideContext) error
