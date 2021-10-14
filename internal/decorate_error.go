package internal

import (
	"reflect"
	"runtime"

	"github.com/rectcircle/digpro/internal/digcopy"
)

func TryFixDigErr(err error, pc uintptr) error {
	if err == nil {
		return nil
	}
	if pc == 0 {
		return err
	}

	errType := reflect.TypeOf(err)

	switch errType.String() {
	case "dig.errProvide", "dig.errConstructorFailed" /* not top level error */, "dig.errArgumentsFailed", "dig.errMissingDependencies":
		errValue := reflect.ValueOf(err)
		errFuncValuePtr := errValue.FieldByName("Func")
		if errFuncValuePtr.IsNil() {
			return err
		}
		errFuncValue := errFuncValuePtr.Elem()
		funcInfo := digcopy.InspectFuncPC(pc)
		errFuncValue.FieldByName("Name").SetString(funcInfo.Name)
		errFuncValue.FieldByName("Package").SetString(funcInfo.Package)
		errFuncValue.FieldByName("File").SetString(funcInfo.File)
		errFuncValue.FieldByName("Line").SetInt(int64(funcInfo.Line))
		return err
	default:
		return err
	}
}

func WrapErrorWithLocationForPC(callSkip int, f func(pc uintptr) error) error {
	pc, _, _, ok := runtime.Caller(callSkip)
	if !ok {
		pc = 0
	}
	return TryFixDigErr(f(pc), pc)
}
