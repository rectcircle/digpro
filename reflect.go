package digpro

import (
	"fmt"
	"reflect"
	"unicode"
	"unsafe"

	"github.com/rectcircle/digpro/internal"
)

func ensureStructFieldExported(f reflect.StructField) reflect.StructField {
	if f.PkgPath == "" { // f.IsExported() 1.17 add
		return f
	}
	name := f.Name
	// first letter upper
	for i, v := range name {
		name = string(unicode.ToUpper(v)) + name[i+1:]
		break
	}
	f.Name = name
	f.PkgPath = ""
	return f
}

func structTypeOf(structOrStructPtr interface{}) (isPtr bool, structTyp reflect.Type, err error) {
	structTyp = reflect.TypeOf(structOrStructPtr)

	if structOrStructPtr == nil || (structTyp.Kind() == reflect.Ptr && reflect.ValueOf(structOrStructPtr).IsNil()) {
		return false, nil, fmt.Errorf("structOrStructPtr want struct or non nil struct pointer, but got %#v", structOrStructPtr)
	}

	switch structTyp.Kind() {
	case reflect.Ptr:
		isPtr = true
		structTyp = structTyp.Elem()
		if structTyp.Kind() != reflect.Struct {
			return false, nil, fmt.Errorf("structOrStructPtr want struct or non nil struct pointer, but got %s", reflect.TypeOf(structOrStructPtr))
		}
	case reflect.Struct:
		isPtr = false
	default:
		return false, nil, fmt.Errorf("structOrStructPtr want struct or struct pointer, but got %s", reflect.TypeOf(structOrStructPtr))
	}
	return
}

func structPtrValueOf(structOrStructPtr interface{}) (isPtr bool, structPtrValue reflect.Value, err error) {
	isPtr, _, err = structTypeOf(structOrStructPtr)
	if err != nil {
		return
	}
	structPtrValue = reflect.ValueOf(structOrStructPtr)
	if structPtrValue.Type().Kind() == reflect.Struct {
		// new a struct ptr, and copy ele to this ptr
		originStructValue := structPtrValue
		structPtrValue = reflect.New(originStructValue.Type())
		structPtrValue.Elem().Set(originStructValue)
	}
	return
}

func makeParameterObjectType(structOrStructPtr interface{}, resolveCyclic bool) (reflect.Type, map[string]int, error) {
	_, structTyp, err := structTypeOf(structOrStructPtr)
	if err != nil {
		return nil, nil, err
	}

	// append digInField to parameterObjectFields first element
	parameterObjectFields := []reflect.StructField{internal.DigInField}

	// map[parameterObjectFieldName]valueFieldIndex
	fieldMapping := map[string]int{}

	if !resolveCyclic {
		for i := 0; i < structTyp.NumField(); i++ {
			originField := structTyp.Field(i)
			f := ensureStructFieldExported(originField)
			if f.Tag.Get("digpro") == "ignore" {
				continue
			}
			if existIdx, ok := fieldMapping[f.Name]; ok {
				existField := structTyp.Field(existIdx)
				return nil, nil, fmt.Errorf("field conflict, %s:%s and %s:%s", originField.Name, originField.Type, existField.Name, existField.Type)
			}
			parameterObjectFields = append(parameterObjectFields, f)
			fieldMapping[f.Name] = i
		}
	}
	return reflect.StructOf(parameterObjectFields), fieldMapping, nil
}

func copyFromParameterObject(structOrStructPtr interface{}, parameterObjectValue reflect.Value, fieldMapping map[string]int) (interface{}, error) {
	isPtr, structPtrValue, err := structPtrValueOf(structOrStructPtr)
	if err != nil {
		return nil, err
	}
	addressableStructValue := structPtrValue.Elem()
	parameterObjectTyp := parameterObjectValue.Type()
	// first element is dig.In, ignore it

	for i := 1; i < parameterObjectValue.NumField(); i++ {
		parameterObjectFieldValue := parameterObjectValue.Field(i)
		structFieldName := parameterObjectTyp.Field(i).Name
		structFieldIndex := fieldMapping[structFieldName]
		structFieldValue := addressableStructValue.Field(structFieldIndex)
		if !structFieldValue.CanSet() {
			structFieldValue = reflect.NewAt(structFieldValue.Type(), unsafe.Pointer(structFieldValue.UnsafeAddr())).Elem()
		}
		structFieldValue.Set(parameterObjectFieldValue)
	}
	if isPtr {
		return structPtrValue.Interface(), nil
	}
	return addressableStructValue.Interface(), nil
}

func underlyingValue(value reflect.Value) reflect.Value {
	if k := value.Kind(); k != reflect.Ptr && k != reflect.Interface {
		return value
	}
	return underlyingValue(value.Elem())
}

// func isReferenceType(t reflect.Type) bool {
// 	switch t.Kind() {
// 	case reflect.Ptr | reflect.Interface:
// 		return true
// 	default:
// 		return false
// 	}
// }
