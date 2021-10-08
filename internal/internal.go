package internal

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/rectcircle/digpro/internal/digcopy"
	"go.uber.org/dig"
)

func TryConvertDigErr(err error, pc uintptr) error {
	if err == nil {
		return nil
	}

	errType := reflect.TypeOf(err)

	switch errType.String() {
	case "dig.errProvide", "dig.errConstructorFailed", "dig.errArgumentsFailed", "dig.errMissingDependencies":
	default:
		return err
	}

	errValue := reflect.ValueOf(err)
	errValueReasonField := errValue.FieldByName("Reason")
	reason := errValueReasonField.Interface().(error)

	switch errType.String() {
	case "dig.errProvide":
		return digcopy.ErrProvide{
			Func:   digcopy.InspectFuncPC(pc),
			Reason: reason,
		}
	case "dig.errConstructorFailed":
		return digcopy.ErrConstructorFailed{
			Func:   digcopy.InspectFuncPC(pc),
			Reason: reason,
		}
	case "dig.errArgumentsFailed":
		return digcopy.ErrArgumentsFailed{
			Func:   digcopy.InspectFuncPC(pc),
			Reason: reason,
		}
	case "dig.errMissingDependencies":
		return digcopy.ErrMissingDependencies{
			Func:   digcopy.InspectFuncPC(pc),
			Reason: reason,
		}
	default:
		return err
	}
}

func ProvideWithLocationForPC(c *dig.Container, callSkip int, constructor interface{}, opts ...dig.ProvideOption) error {
	pc, _, _, ok := runtime.Caller(callSkip)
	if ok {
		return TryConvertDigErr(c.Provide(constructor, append([]dig.ProvideOption{dig.LocationForPC(pc)}, opts...)...), pc)
	} else {
		return c.Provide(constructor, opts...)
	}
}

type ExtractOptions struct {
	Name  string
	Group string
}

type ExtractOption interface {
	ApplyExtractOption(*ExtractOptions)
}

func ExtractWithLocationForPC(c *dig.Container, callSkip int, typInterface interface{}, opts ...ExtractOption) (interface{}, error) {
	switch v := typInterface.(type) {
	case nil:
		return nil, fmt.Errorf("typInterface want value or not nil pointer, but got %v", v)
	case error:
		return nil, fmt.Errorf("typInterface want value or not nil pointer, but got %v", v)
	}

	var options ExtractOptions
	for _, o := range opts {
		o.ApplyExtractOption(&options)
	}
	tags := []string{}
	if options.Name != "" {
		tags = append(tags, fmt.Sprintf(`name:"%s"`, options.Name))
	}
	if options.Group != "" {
		tags = append(tags, fmt.Sprintf(`group:"%s"`, options.Group))
	}

	argFields := []reflect.StructField{DigInField, {
		Name: "Value",
		Type: reflect.TypeOf(typInterface),
		Tag:  reflect.StructTag(strings.Join(tags, " ")),
	}}
	argsType := reflect.StructOf(argFields)

	var returnInterface interface{}

	argsTypes := []reflect.Type{argsType}
	returnTypes := []reflect.Type{}

	ft := reflect.FuncOf(argsTypes, returnTypes, false)
	fv := reflect.MakeFunc(ft, func(args []reflect.Value) []reflect.Value {
		returnInterface = args[0].Field(1).Interface()
		return nil
	})
	f := fv.Interface()

	pc, _, _, ok := runtime.Caller(callSkip)
	if ok {
		return returnInterface, TryConvertDigErr(c.Invoke(f), pc)
	} else {
		return returnInterface, c.Invoke(f)
	}
}
