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
	return TryFixDigErrByFunc(err, digcopy.InspectFuncPC(pc))
}

func IsDigErrMissingDependencies(err error) bool {
	if err == nil {
		return false
	}
	return reflect.TypeOf(err).String() == "dig.errMissingDependencies"
}

func TryFixDigErrByFunc(err error, location *digcopy.Func) error {

	if err == nil {
		return nil
	}
	if location == nil {
		return err
	}

	if _, ok := err.(digcopy.ErrParamSingleFailed); ok {
		err = digcopy.ErrArgumentsFailed{
			Func:   location,
			Reason: err,
		}
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

		errFuncValue.FieldByName("Name").SetString(location.Name)
		errFuncValue.FieldByName("Package").SetString(location.Package)
		errFuncValue.FieldByName("File").SetString(location.File)
		errFuncValue.FieldByName("Line").SetInt(int64(location.Line))
		return err
	default:
		return err
	}
}

func WrapResolveCyclicError(err error, location *digcopy.Func, output *ProvideOutput) error {
	if err == nil {
		return nil
	}
	return digcopy.ErrParamSingleFailed{
		Key:    digcopy.Key{T: output.Type, Name: output.Name, Group: output.Group},
		Reason: TryFixDigErrByFunc(err, location),
	}
}

func WrapErrorWithLocationForPC(callSkip int, f func(pc uintptr) error) error {
	pc, _, _, ok := runtime.Caller(callSkip)
	if !ok {
		pc = 0
	}
	return TryFixDigErr(f(pc), pc)
}
