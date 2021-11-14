package digpro

import (
	"reflect"

	"github.com/rectcircle/digpro/internal"
	"github.com/rectcircle/digpro/internal/digcopy"
	"go.uber.org/dig"
)

// ResolveCyclic
// TODO add doc
func ResolveCyclic() dig.ProvideOption {
	return resolveCyclicProvideOption{}
}

type resolveCyclicOriginProvideInfoProvideOption struct {
	dig.ProvideOption
	provideInfo *internal.ProvideInfosWrapper
}

// resolveCyclicProvideMiddleware record all provide inputs
func resolveCyclicProvideMiddleware(pc *provideContext) error {
	var (
		resolveCyclic = false
		provideInfo   *internal.ProvideInfosWrapper
		_opts         = make([]dig.ProvideOption, 0, len(pc.opts))
	)
	for _, opt := range pc.opts {
		if o, ok := opt.(resolveCyclicOriginProvideInfoProvideOption); ok {
			provideInfo = o.provideInfo
			resolveCyclic = true
		} else {
			_opts = append(_opts, opt)
		}
	}
	pc.opts = _opts
	err := pc.next()
	if err != nil {
		return err
	}
	if provideInfo == nil {
		provideInfo = &pc.c.provideInfos[len(pc.c.provideInfos)-1]
	}
	propertyInfo := internal.PropertyInfo{
		ResolveCyclic: resolveCyclic,
		Inputs:        provideInfo.ExportedInputs(),
		Injected:      false,
		Error:         err,
	}
	for _, output := range provideInfo.ExportedOutputs() {
		pc.c.propertyInjects[output] = &propertyInfo
	}
	return nil
}

func (c *ContainerWrapper) getLocationByOutput(input internal.ProvideOutput) *digcopy.Func {
	containerValue := reflect.ValueOf(&c.Container).Elem()

	providersValue := internal.EnsureValueExported(containerValue.FieldByName("providers")) // map[dig.key][]*dig.node

	keyType := providersValue.Type().Key()

	key := reflect.New(keyType).Elem()
	internal.EnsureValueExported(key.FieldByName("t")).Set(reflect.ValueOf(input.Type))
	internal.EnsureValueExported(key.FieldByName("name")).Set(reflect.ValueOf(input.Name))
	internal.EnsureValueExported(key.FieldByName("group")).Set(reflect.ValueOf(input.Group))

	node := providersValue.MapIndex(key)
	if !node.IsValid() {
		return nil
	}
	if node.Len() != 1 {
		// dead code
		return nil
	}
	node = node.Index(0).Elem()
	originLocation := internal.EnsureValueExported(node.FieldByName("location")).Elem()
	location := &digcopy.Func{}
	location.Name = originLocation.FieldByName("Name").Interface().(string)
	location.Package = originLocation.FieldByName("Package").Interface().(string)
	location.File = originLocation.FieldByName("File").Interface().(string)
	location.Line = originLocation.FieldByName("Line").Interface().(int)
	return location
}

// doPropertyInject for arg do property inject
func (c *ContainerWrapper) doPropertyInject(arg reflect.Value, structField *reflect.StructField) error {
	// arg is In type for everyone field to call doPropertyInject
	if dig.IsIn(arg.Type()) {
		for i := 0; i < arg.NumField(); i++ {
			sf := arg.Type().Field(i)
			if dig.IsIn(sf.Type) {
				continue
			}
			err := c.doPropertyInject(arg.Field(i), &sf)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// make key of arg
	key := internal.ProvideOutput{
		Type: arg.Type(),
	}
	if structField != nil {
		if group := structField.Tag.Get(internal.DigGroupTag); group != "" {
			key.Group = group
		} else if name := structField.Tag.Get(internal.DigNameTag); name != "" {
			key.Name = name
		}
	}
	// get property inject info
	propertyInject := c.propertyInjects[key]
	if propertyInject == nil || propertyInject.Injected {
		return nil
	}
	if propertyInject.Error != nil {
		return propertyInject.Error
	}
	// do property inject
	propertyInject.Injected = true
	for _, input := range propertyInject.Inputs {
		// recursion
		inputValue := reflect.New(input.Type)
		if input.Type.Kind() != reflect.Interface {
			inputValue = inputValue.Elem()
		}
		// fmt.Println(input.Type)
		value, err := internal.ExtractWithLocationForPC(c.Invoke, 0, inputValue.Interface(), ExtractByName(input.Name), ExtractByGroup(input.Group))
		// fmt.Println("-", input.Type)
		if err != nil {
			if input.Optional && internal.IsDigErrMissingDependencies(err) {
				continue
			} else {
				propertyInject.Injected = false
				propertyInject.Error = internal.WrapResolveCyclicError(err, c.getLocationByOutput(key), &key)
				return propertyInject.Error
			}
		}
		// only do property inject when enable resolve cyclic option
		if propertyInject.ResolveCyclic {
			argStruct := internal.EnsureValueExported(underlyingValue(arg))
			if argStruct.Kind() != reflect.Struct {
				return nil
			}

			for i := 0; i < argStruct.NumField(); i++ {
				fieldValue := internal.EnsureValueExported(argStruct.Field(i))
				fieldType := argStruct.Type().Field(i)
				if fieldType.Type == input.Type &&
					fieldType.Tag.Get(internal.DigNameTag) == input.Name &&
					fieldType.Tag.Get(internal.DigGroupTag) == input.Group {
					fieldValue.Set(reflect.ValueOf(value))
					break
				}
			}
		}
	}

	return nil
}
