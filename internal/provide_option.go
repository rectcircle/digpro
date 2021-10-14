package internal

import (
	"reflect"

	"go.uber.org/dig"
)

type ProvideOptions struct {
	Name  string
	Group string
	Info  *dig.ProvideInfo
	As    []interface{}
}

func ApplyProvideOptions(opts ...dig.ProvideOption) ProvideOptions {
	DigProvideOptionsPtrValue := reflect.New(DigProvideOptionsType)
	for _, opt := range opts {
		optValue := reflect.ValueOf(opt)
		optValue.Call([]reflect.Value{DigProvideOptionsPtrValue})
	}
	DigProvideOptionValue := DigProvideOptionsPtrValue.Elem()
	provideOptions := ProvideOptions{
		Name:  DigProvideOptionValue.FieldByName("Name").Interface().(string),
		Group: DigProvideOptionValue.FieldByName("Group").Interface().(string),
		Info:  DigProvideOptionValue.FieldByName("Info").Interface().(*dig.ProvideInfo),
		As:    DigProvideOptionValue.FieldByName("As").Interface().([]interface{}),
	}
	return provideOptions
}
