package digpro

import (
	"fmt"
	"reflect"

	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

func wrapError(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] %s", prefix, err.Error())
}

// Struct make a struct constructor.
//
// support all dig tags and `digpro:"ignore"`
//
//   struct {
//   	A string   `name:"a"`
//   	B []string `group:"b"`
//   	C bool     `optional:"true"`
//   	D string   `digpro:"ignore"`  // ignore this field
//   }
//
// for example
//   type Foo struct {
//   	A       string
//   	B       int
//   	C       int `name:"c"`
//   	private bool
//   	ignore  int `digpro:"ignore"`
//   }
//   // new a *dig.Container
//   c := dig.New()
//   // provide some constructor
//   digpro.QuickPanic(
//   	// register object
//   	c.Provide(digpro.Supply("a")),
//   	c.Provide(digpro.Supply(1)),
//   	c.Provide(digpro.Supply(true)),
//   	c.Provide(digpro.Supply(2), dig.Name("c")),
//   	// register a struct
//   	c.Provide(digpro.Struct(Foo{
//   		ignore: 3,
//   	})),
//   	// equals to
//   	// c.Provide(func(in struct {
//   	// 	A       string
//   	// 	B       int
//   	// 	C       int `name:"c"`
//   	// 	Private bool
//   	// }) Foo {
//   	// 	return Foo{
//   	// 		A:       in.A,
//   	// 		B:       in.B,
//   	// 		C:       in.C,
//   	// 		private: in.Private,
//   	// 		ignore:  3,
//   	// 	}
//   	// }),
//   )
//   // extract object from container
//   foo, err := digpro.Extract(c, Foo{})
//   if err != nil {
//   	digpro.QuickPanic(err)
//   }
//   fmt.Printf("%#v", foo)
//   // Output: digpro_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
func Struct(structOrStructPtr interface{}) interface{} {
	parameterObjectType, fieldMapping, err := makeParameterObjectType(structOrStructPtr)
	if err != nil {
		return wrapError("Struct", err)
	}

	parameterTypes := []reflect.Type{parameterObjectType}
	structOrStructPtrType := reflect.TypeOf(structOrStructPtr)
	returnTypes := []reflect.Type{structOrStructPtrType, internal.ErrorType}

	ft := reflect.FuncOf(parameterTypes, returnTypes, false)
	fv := reflect.MakeFunc(ft, func(p []reflect.Value) []reflect.Value {
		// copy from parameter to injectedObject and return
		injectedObject, err := copyFromParameterObject(structOrStructPtr, p[0], fieldMapping)
		errValue := reflect.ValueOf(wrapError("Struct", err))
		var injectedObjectValue reflect.Value
		if err == nil {
			// handle result to value
			injectedObjectValue = reflect.ValueOf(injectedObject)
			errValue = reflect.New(internal.ErrorType).Elem()
		} else {
			injectedObjectValue = reflect.New(structOrStructPtrType).Elem()
		}

		return []reflect.Value{injectedObjectValue, errValue}
	})
	return fv.Interface()
}

// Struct make a struct constructor.
//
// support all dig tags and `digpro:"ignore"`
//
//   struct {
//   	A string   `name:"a"`
//   	B []string `group:"b"`
//   	C bool     `optional:"true"`
//   	D string   `digpro:"ignore"`  // ignore this field
//   }
//
// for example
//   type Foo struct {
//   	A       string
//   	B       int
//   	C       int `name:"c"`
//   	private bool
//   	ignore  int `digpro:"ignore"`
//   }
//   // new a *dig.Container
//   c := digpro.New()
//   // provide some constructor
//   digpro.QuickPanic(
//   	// register object
//   	c.Supply("a"),
//   	c.Supply(1),
//   	c.Supply(true),
//   	c.Supply(2, dig.Name("c")),
//   	// register a struct
//   	c.Struct(Foo{
//   		ignore: 3,
//   	}),
//   	// equals to
//   	// c.Provide(func(in struct {
//   	// 	A       string
//   	// 	B       int
//   	// 	C       int `name:"c"`
//   	// 	Private bool
//   	// }) Foo {
//   	// 	return Foo{
//   	// 		A:       in.A,
//   	// 		B:       in.B,
//   	// 		C:       in.C,
//   	// 		private: in.Private,
//   	// 		ignore:  3,
//   	// 	}
//   	// }),
//   )
//   // extract object from container
//   foo, err := c.Extract(Foo{})
//   if err != nil {
//   	digpro.QuickPanic(err)
//   }
//   fmt.Printf("%#v", foo)
//   // Output: digpro_test.Foo{A:"a", B:1, C:2, private:true, ignore:3}
func (c *ContainerWrapper) Struct(structOrStructPtr interface{}, opts ...dig.ProvideOption) error {
	return internal.ProvideWithLocationForPC(c.Provide, 3, Struct(structOrStructPtr), opts...)
}
