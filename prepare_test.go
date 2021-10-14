package digpro

import (
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
