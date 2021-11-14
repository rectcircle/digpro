package internal

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/dig"
)

type ExtractOptions struct {
	Name  string
	Group string
}

type ExtractOption interface {
	ApplyExtractOption(*ExtractOptions)
}

type ExtractOptionFunc func(*ExtractOptions)

func (f ExtractOptionFunc) ApplyExtractOption(opts *ExtractOptions) { f(opts) }

func MakeExtractFunc(ptr interface{}, opts ...ExtractOption) interface{} {
	if ptr == nil {
		return fmt.Errorf("[MakeExtractFunc] can't extract an untyped nil")
	}
	ptrValue := reflect.ValueOf(ptr)
	ptrType := ptrValue.Type()
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("[MakeExtractFunc] ptr must is pointer")
	}
	if ptrValue.IsNil() {
		return fmt.Errorf("[MakeExtractFunc] can't extract an untyped nil") // dead code ?
	}
	if ptrType.Elem() == ErrorType {
		return fmt.Errorf("[MakeExtractFunc] can't extract an error")
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
		Type: ptrType.Elem(),
		Tag:  reflect.StructTag(strings.Join(tags, " ")),
	}}
	argsType := reflect.StructOf(argFields)

	argsTypes := []reflect.Type{argsType}
	returnTypes := []reflect.Type{}

	ft := reflect.FuncOf(argsTypes, returnTypes, false)
	fv := reflect.MakeFunc(ft, func(args []reflect.Value) []reflect.Value {
		ptrValue.Elem().Set(args[0].Field(1))
		return nil
	})
	return fv.Interface()
}

func getPtrFinalKind(t reflect.Type) reflect.Kind {
	if k := t.Kind(); k != reflect.Ptr {
		return k
	}
	return getPtrFinalKind(t.Elem())
}

func ExtractWithLocationForPC(Invoke func(function interface{}, opts ...dig.InvokeOption) error, callSkip int, typ interface{}, opts ...ExtractOption) (interface{}, error) {
	if typ == nil {
		return nil, fmt.Errorf("can't extract an untyped nil")
	}
	typPtrInterface := reflect.New(reflect.TypeOf(typ))
	// pointer of interface will do once addressing operation
	// that means Extract(*interfaceA)) will return -> interfaceA
	if reflect.TypeOf(typ).Kind() == reflect.Ptr && getPtrFinalKind(reflect.TypeOf(typ)) == reflect.Interface {
		typPtrInterface = reflect.New(reflect.ValueOf(typ).Elem().Type())
	}
	f := MakeExtractFunc(typPtrInterface.Interface(), opts...)
	if err, ok := f.(error); ok {
		return nil, err
	}
	var err error
	if callSkip <= 0 {
		err = Invoke(f)
	} else {
		err = WrapErrorWithLocationForPC(callSkip, func(uintptr) error { return Invoke(f) })
	}
	if err != nil {
		return nil, err
	}
	return typPtrInterface.Elem().Interface(), nil
}
