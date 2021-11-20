package internal

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"go.uber.org/dig"
)

type ProvideInfosWrapper struct {
	dig.ProvideInfo
	exportedOutputs []ProvideOutput
	exportedInputs  []ProvideInput
}

func EnsureValueExported(value reflect.Value) reflect.Value {
	if !value.CanSet() {
		return reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem()
	} else {
		return value
	}
}

func (piw *ProvideInfosWrapper) ExportedOutputs() []ProvideOutput {
	if piw.exportedOutputs != nil {
		return piw.exportedOutputs
	}
	result := make([]ProvideOutput, 0, len(piw.Outputs))
	for _, output := range piw.Outputs {
		outputValue := reflect.ValueOf(output).Elem()
		result = append(result, ProvideOutput{
			Type:  EnsureValueExported(outputValue.FieldByName("t")).Interface().(reflect.Type),
			Name:  EnsureValueExported(outputValue.FieldByName("name")).Interface().(string),
			Group: EnsureValueExported(outputValue.FieldByName("group")).Interface().(string),
		})
	}
	piw.exportedOutputs = result
	return result
}

func (piw *ProvideInfosWrapper) ExportedInputs() []ProvideInput {
	if piw.exportedInputs != nil {
		return piw.exportedInputs
	}
	result := make([]ProvideInput, 0, len(piw.Inputs))
	for _, input := range piw.Inputs {
		inputValue := reflect.ValueOf(input).Elem()
		result = append(result, ProvideInput{
			Type:     EnsureValueExported(inputValue.FieldByName("t")).Interface().(reflect.Type),
			Optional: EnsureValueExported(inputValue.FieldByName("optional")).Interface().(bool),
			Name:     EnsureValueExported(inputValue.FieldByName("name")).Interface().(string),
			Group:    EnsureValueExported(inputValue.FieldByName("group")).Interface().(string),
		})
	}
	piw.exportedInputs = result
	return result
}

type ProvideOutput struct {
	Type        reflect.Type
	Name, Group string
}

func (po *ProvideOutput) String() string {
	if po.Name == "" && po.Group == "" {
		return po.Type.String()
	} else {
		s := []string{}
		if po.Name != "" {
			s = append(s, fmt.Sprintf("name=\"%s\"", po.Name))
		}
		if po.Group != "" {
			s = append(s, fmt.Sprintf("group=\"%s\"", po.Group))
		}
		return fmt.Sprintf("%s[%s]", po.Type.String(), strings.Join(s, ","))
	}
}

func EqualsProvideOutputs(a []ProvideOutput, b []ProvideOutput) bool {
	if len(a) != len(b) {
		return false
	}
	am := make(map[ProvideOutput]struct{}, len(a))
	bm := make(map[ProvideOutput]struct{}, len(a))
	for _, o := range a {
		am[o] = struct{}{}
	}
	for _, o := range b {
		bm[o] = struct{}{}
	}
	return reflect.DeepEqual(am, bm)
}

type ProvideInput struct {
	Type        reflect.Type
	Optional    bool
	Name, Group string
}

func (po *ProvideInput) String() string {
	if po.Name == "" && po.Group == "" {
		return po.Type.String()
	} else {
		s := []string{}
		if po.Name != "" {
			s = append(s, fmt.Sprintf("name=\"%s\"", po.Name))
		}
		if po.Group != "" {
			s = append(s, fmt.Sprintf("group=\"%s\"", po.Group))
		}
		if po.Optional {
			s = append(s, "optional")
		}
		return fmt.Sprintf("%s[%s]", po.Type.String(), strings.Join(s, ","))
	}
}
