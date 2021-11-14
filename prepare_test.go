package digpro

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type EmptyInner struct {
}

type NonEmptyInner struct {
	I1 string
	I2 string
}

type Foo struct {
	EmptyInner
	A int
	NonEmptyInner
	S struct {
		S1 int
		S2 int
	}
	Arr     []int
	B       string
	private bool //lint:ignore U1000 for test
}

type Bar struct {
	A       string
	B       int
	private bool
}

type Biz struct {
	A int      `name:"a"`
	B int      `digpro:"ignore"`
	C []string `group:"c"`
}

var mockDataMapping = map[reflect.Kind]reflect.Value{
	reflect.Bool:    reflect.ValueOf(true),
	reflect.Int:     reflect.ValueOf(1),
	reflect.Int8:    reflect.ValueOf(2),
	reflect.Int16:   reflect.ValueOf(3),
	reflect.Int32:   reflect.ValueOf(4),
	reflect.Int64:   reflect.ValueOf(5),
	reflect.Uint:    reflect.ValueOf(6),
	reflect.Uint8:   reflect.ValueOf(7),
	reflect.Uint16:  reflect.ValueOf(8),
	reflect.Uint32:  reflect.ValueOf(9),
	reflect.Uint64:  reflect.ValueOf(10),
	reflect.Uintptr: reflect.ValueOf(11),
	reflect.Float32: reflect.ValueOf(12),
	reflect.Float64: reflect.ValueOf(13),
	reflect.String:  reflect.ValueOf("string"),
}

func mockParameterObject(typ reflect.Type) reflect.Value {
	value := reflect.New(typ).Elem()
	for i := 0; i < value.NumField(); i++ {
		f := value.Field(i)
		if v, ok := mockDataMapping[f.Type().Kind()]; ok {
			f.Set(v)
		}
	}
	return value
}

func assertStructOrStructPtr(t *testing.T, structOrStructPtr interface{}) (pass bool) {
	pass = true
	_, structPtrValue, err := structPtrValueOf(structOrStructPtr)
	structValue := structPtrValue.Elem()
	if err != nil {
		t.Error(err)
		pass = false
		return
	}
	for i := 0; i < structValue.NumField(); i++ {
		structFieldValue := structValue.Field(i)
		if !structFieldValue.CanSet() {
			structFieldValue = reflect.NewAt(structFieldValue.Type(), unsafe.Pointer(structFieldValue.UnsafeAddr())).Elem()
		}
		if mockValue, ok := mockDataMapping[structFieldValue.Type().Kind()]; ok {
			// fmt.Printf("field %v got = %v, want %v\n", structFieldValue, mockValue.Interface(), structFieldValue.Interface())
			if !reflect.DeepEqual(structFieldValue.Interface(), mockValue.Interface()) {
				pass = false
				t.Errorf("field %v got = %v, want %v", structFieldValue, mockValue.Interface(), structFieldValue.Interface())
			}
		}
	}
	return pass
}

type D1 struct {
	D2    *D2
	Value int
}

func (d1 *D1) String() string {
	return fmt.Sprintf("D1: {D2: {D1: ..., Value: '%s'}, Value: %d}", d1.D2.Value, d1.Value)
}

type D2 struct {
	D1    *D1
	Value string
}

func (d2 *D2) String() string {
	return fmt.Sprintf("D2: {D1: {D2: ..., Value: %d}, Value: '%s'}", d2.D1.Value, d2.Value)
}

type D3 struct {
	D2    *D2
	Value bool
}

func (d3 D3) String() string {
	return fmt.Sprintf("D3 {%s, Value: %t}", d3.D2.String(), d3.Value)
}

type D4 struct {
	D2     *D2     `optional:"true"`
	Value  float64 `optional:"true"`
	Ignore int     `digpro:"ignore"`
}

func (d4 *D4) String() string {
	d2String := "null"
	if d4.D2 != nil {
		d2String = d4.D2.String()
	}
	return fmt.Sprintf("D4 {%s, Value: %f, Ignore: %d}", d2String, d4.Value, d4.Ignore)
}

type I1 interface{ String() string }
type I2 interface{ String() string }

type DI1 struct {
	I2    I2
	Value int
}

func (d1 *DI1) String() string {
	return fmt.Sprintf("DI1: {I2: {I1: ..., Value: '%s'}, Value: %d}", d1.I2.(*DI2).Value, d1.Value)
}

type DI2 struct {
	I1    I1
	Value string
}

func (d2 *DI2) String() string {
	return fmt.Sprintf("DI2: {I1: {I2: ..., Value: %d}, Value: '%s'}", d2.I1.(*DI1).Value, d2.Value)
}

type PrepareFunc func(c *ContainerWrapper) error
