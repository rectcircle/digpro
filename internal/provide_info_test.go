package internal

import (
	"reflect"
	"testing"

	"go.uber.org/dig"
)

func TestEqualsProvideOutputs(t *testing.T) {
	type args struct {
		a []ProvideOutput
		b []ProvideOutput
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "0",
			args: args{},
			want: true,
		},
		{
			name: "one equal",
			args: args{
				a: []ProvideOutput{{
					Type:  reflect.TypeOf(1),
					Name:  "",
					Group: "",
				}},
				b: []ProvideOutput{{
					Type:  reflect.TypeOf(1),
					Name:  "",
					Group: "",
				}},
			},
			want: true,
		},
		{
			name: "one not equal",
			args: args{
				a: []ProvideOutput{{
					Type:  reflect.TypeOf(1),
					Name:  "",
					Group: "",
				}},
				b: []ProvideOutput{{
					Type:  reflect.TypeOf(""),
					Name:  "",
					Group: "",
				}},
			},
			want: false,
		},
		{
			name: "len not equal",
			args: args{
				a: []ProvideOutput{
					{
						Type:  reflect.TypeOf(1),
						Name:  "",
						Group: "",
					},
					{
						Type:  reflect.TypeOf(""),
						Name:  "",
						Group: "",
					},
				},
				b: []ProvideOutput{{
					Type:  reflect.TypeOf(""),
					Name:  "",
					Group: "",
				}},
			},
			want: false,
		},
		{
			name: "order different",
			args: args{
				a: []ProvideOutput{
					{
						Type:  reflect.TypeOf(1),
						Name:  "int",
						Group: "",
					},
					{
						Type:  reflect.TypeOf(""),
						Name:  "",
						Group: "",
					},
				},
				b: []ProvideOutput{
					{
						Type:  reflect.TypeOf(""),
						Name:  "",
						Group: "",
					},
					{
						Type:  reflect.TypeOf(1),
						Name:  "int",
						Group: "",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualsProvideOutputs(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("EqualsProvideOutputs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvideOutput_String(t *testing.T) {
	type fields struct {
		Type  reflect.Type
		Name  string
		Group string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "int",
			fields: fields{
				Type:  reflect.TypeOf(1),
				Name:  "",
				Group: "",
			},
			want: "int",
		},
		{
			name: "int name",
			fields: fields{
				Type:  reflect.TypeOf(1),
				Name:  "a",
				Group: "",
			},
			want: "int[name=\"a\"]",
		},
		{
			name: "int group",
			fields: fields{
				Type:  reflect.TypeOf(1),
				Name:  "",
				Group: "g",
			},
			want: "int[group=\"g\"]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			po := &ProvideOutput{
				Type:  tt.fields.Type,
				Name:  tt.fields.Name,
				Group: tt.fields.Group,
			}
			if got := po.String(); got != tt.want {
				t.Errorf("ProvideOutput.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvideInfoWrapper_ExportedOutputs(t *testing.T) {
	type args struct {
		constructor interface{}
		opts        []dig.ProvideOption
	}
	tests := []struct {
		name string
		arg  args
		want []ProvideOutput
	}{
		{
			name: "one type",
			arg: args{
				constructor: func() int {
					return 1
				},
				opts: []dig.ProvideOption{},
			},
			want: []ProvideOutput{
				{
					Type:  reflect.TypeOf(1),
					Name:  "",
					Group: "",
				},
			},
		},
		{
			name: "one named type",
			arg: args{
				constructor: func() int {
					return 1
				},
				opts: []dig.ProvideOption{dig.Name("a")},
			},
			want: []ProvideOutput{
				{
					Type:  reflect.TypeOf(1),
					Name:  "a",
					Group: "",
				},
			},
		},
		{
			name: "one group type",
			arg: args{
				constructor: func() int {
					return 1
				},
				opts: []dig.ProvideOption{dig.Group("g")},
			},
			want: []ProvideOutput{
				{
					Type:  reflect.TypeOf(1),
					Name:  "",
					Group: "g",
				},
			},
		},
		{
			name: "one named type and as interface{}",
			arg: args{
				constructor: func() int {
					return 1
				},
				opts: []dig.ProvideOption{dig.Name("a"), dig.As(new(interface{}))},
			},
			want: []ProvideOutput{
				{
					Type:  reflect.TypeOf(new(interface{})).Elem(),
					Name:  "a",
					Group: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			piw := ProvideInfosWrapper{}
			c := dig.New()
			err := c.Provide(tt.arg.constructor, append(tt.arg.opts, dig.FillProvideInfo(&piw.ProvideInfo))...)
			if err != nil {
				t.Errorf("c.Provide() got err = %v", err)
				return
			}
			if got := piw.ExportedOutputs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProvideInfoWrapper.ExportedOutputs() = %v, want %v", got, tt.want)
			}
			// use cache
			if got := piw.ExportedOutputs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProvideInfoWrapper.ExportedOutputs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnsureValueExported(t *testing.T) {
	a := &struct{ a string }{
		a: "a",
	}
	fieldAValue := reflect.ValueOf(a).Elem().FieldByName("a")
	exportedAValue := EnsureValueExported(fieldAValue)

	if fieldAValue == exportedAValue {
		t.Errorf("two fieldAValue and exportedAValue want equals, got not equals")
		return
	}

	if exportedAValue != EnsureValueExported(exportedAValue) {
		t.Errorf("two EnsureValueExported will equal")
		return
	}
	if exportedAValue.Interface().(string) != "a" {
		t.Errorf("a.a want %s, got %s", "a", exportedAValue.Interface().(string))
		return
	}
	exportedAValue.SetString("b")
	if a.a != "b" {
		t.Errorf("a.a want %s, got %s", "b", a.a)
		return
	}

	b := &struct{ B string }{
		B: "a",
	}
	fieldBValue := reflect.ValueOf(b).Elem().FieldByName("B")
	exportedBValue := EnsureValueExported(fieldBValue)
	if fieldBValue != exportedBValue {
		t.Errorf("two fieldBValue and exportedBValue want not equals, got equals")
		return
	}

}
