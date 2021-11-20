package digpro

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

// ContainerWrapper is a dig.Container wrapper, for add some method
type ContainerWrapper struct {
	dig.Container
	middlewares              []provideMiddleware
	provideInfos             []internal.ProvideInfosWrapper
	existResolveCyclicOption bool
	propertyInjects          map[internal.ProvideOutput]*internal.PropertyInfo
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
			resolveCyclicProvideMiddleware,
			overrideProvideMiddleware,
		},
		propertyInjects: make(map[internal.ProvideOutput]*internal.PropertyInfo),
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

func (c *ContainerWrapper) Invoke(function interface{}, opts ...dig.InvokeOption) error {

	_opts, digproOpts := filterInvokeOptionAndGetDigproOptions(opts, locationFixOptionType)
	opts = _opts
	callSkip := digproOpts.locationFixCallSkip
	if callSkip == 0 {
		callSkip = 2
	}

	// pruning
	if !c.existResolveCyclicOption {
		return c.Container.Invoke(function, opts...)
	}

	// check type, see https://github.com/uber-go/dig/blob/v1.13.0/dig.go#L561
	ftype := reflect.TypeOf(function)
	if ftype == nil {
		return errors.New("can't invoke an untyped nil")
	}
	if ftype.Kind() != reflect.Func {
		return fmt.Errorf("can't invoke non-function %v (type %v)", function, ftype)
	}
	var (
		fargs                 []reflect.Value // function arguments
		doPropertyInjectError error           // error for doPropertyInject
	)
	// build a function that arguments as same as function and return void.
	// call this function by dig Invoke
	inTypes := make([]reflect.Type, 0, ftype.NumIn())
	for i := 0; i < ftype.NumIn(); i++ {
		inTypes = append(inTypes, ftype.In(i))
	}
	err := internal.WrapErrorWithLocationForPC(callSkip, func(pc uintptr) error {
		return c.Container.Invoke(reflect.MakeFunc(reflect.FuncOf(inTypes, nil, false), func(args []reflect.Value) (results []reflect.Value) {
			// for everyone arg call doPropertyInject
			for _, arg := range args {
				doPropertyInjectError = c.doPropertyInject(arg, nil)
				if doPropertyInjectError != nil {
					return
				}
			}
			fargs = args
			return
		}).Interface())
	})
	if err != nil {
		return err
	}
	// check error
	if doPropertyInjectError != nil {
		return doPropertyInjectError
	}

	// see: https://github.com/uber-go/dig/blob/v1.13.0/dig.go#L594
	// c.invokerFn(reflect.ValueOf(function), args)
	digContainerValue := reflect.ValueOf(&c.Container).Elem()
	invokerFnField := internal.EnsureValueExported(digContainerValue.FieldByName("invokerFn"))
	invokerFnFieldReturns := invokerFnField.Call([]reflect.Value{reflect.ValueOf(reflect.ValueOf(function)), reflect.ValueOf(fargs)})
	returned := invokerFnFieldReturns[0].Interface().([]reflect.Value)
	if len(returned) == 0 {
		return nil
	}
	if last := returned[len(returned)-1]; last.Type().Implements(internal.ErrorType) {
		if err, _ := last.Interface().(error); err != nil {
			return err
		}
	}
	return nil
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
	info := internal.ProvideInfosWrapper{}
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
