package internal

import (
	"reflect"

	"go.uber.org/dig"
)

var DigInField = reflect.TypeOf(struct{ dig.In }{}).Field(0)
var ErrorType = reflect.TypeOf(new(error)).Elem()
