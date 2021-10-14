package internal

import (
	"reflect"

	"go.uber.org/dig"
)

var DigInField = reflect.TypeOf(struct{ dig.In }{}).Field(0)
var ErrorType = reflect.TypeOf(new(error)).Elem()

var DigProvideOptionsType reflect.Type // dig.provideOptions

func initDigProvideOptionsType() reflect.Type {
	typ := reflect.TypeOf(new(dig.ProvideOption)).Elem()
	method, _ := typ.MethodByName("applyProvideOption")
	methodIn0 := method.Type.In(0).Elem()
	return methodIn0
}

func init() {
	DigProvideOptionsType = initDigProvideOptionsType()
}
