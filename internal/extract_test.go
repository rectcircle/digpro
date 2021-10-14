package internal

import (
	"reflect"
	"strings"
	"testing"

	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

func TestExtractOptionFunc_ApplyExtractOption(t *testing.T) {
	type args struct {
		opts *ExtractOptions
	}
	tests := []struct {
		name string
		f    ExtractOptionFunc
		args args
		want *ExtractOptions
	}{
		{
			name: "",
			f: func(eo *ExtractOptions) {
				eo.Name = "abc"
			},
			args: args{&ExtractOptions{}},
			want: &ExtractOptions{Name: "abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.ApplyExtractOption(tt.args.opts)
			if !reflect.DeepEqual(tt.args.opts, tt.want) {
				t.Errorf("ApplyExtractOption() args want %#v, got %#v", tt.want, tt.args.opts)
				return
			}
		})
	}
}

var a int = 1
var ap *int = &a
var interfaceA interface{} = a
var interfaceAP *interface{} = &interfaceA

func TestMakeExtractFunc(t *testing.T) {
	type extractArgs struct {
		provider tests.Provider
		ptr      interface{}
		opts     []ExtractOption
	}
	var extractTests = []struct {
		name           string
		args           extractArgs
		want           interface{}
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "success as and named type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.Name("a"), dig.As(new(interface{}))),
				ptr:      new(interface{}),
				opts:     []ExtractOption{ExtractOptionFunc(func(eo *ExtractOptions) { eo.Name = "a" })},
			},
			want:    &interfaceA,
			wantErr: false,
		},
		{
			name: "success as type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.As(new(interface{}))),
				ptr:      new(interface{}),
			},
			want:    &interfaceA,
			wantErr: false,
		},
		{
			name: "success group type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.Group("g")),
				ptr:      new([]int),
				opts:     []ExtractOption{ExtractOptionFunc(func(eo *ExtractOptions) { eo.Group = "g" })},
			},
			want:    &[]int{1},
			wantErr: false,
		},
		{
			name: "success named type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.Name("a")),
				ptr:      new(int),
				opts:     []ExtractOption{ExtractOptionFunc(func(eo *ExtractOptions) { eo.Name = "a" })},
			},
			want:    &a,
			wantErr: false,
		},
		{
			name: "success interface pointer",
			args: extractArgs{
				provider: tests.ProviderOne(func() *interface{} { return interfaceAP }),
				ptr:      new(*interface{}),
			},
			want:    &interfaceAP,
			wantErr: false,
		},
		{
			name: "success interface",
			args: extractArgs{
				provider: tests.ProviderOne(func() interface{} { return interfaceA }),
				ptr:      new(interface{}),
			},
			want:    &interfaceA,
			wantErr: false,
		},
		{
			name: "success type pointer",
			args: extractArgs{
				provider: tests.ProviderOne(func() *int { return &a }),
				ptr:      new(*int),
			},
			want:    &ap,
			wantErr: false,
		},
		{
			name: "success type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }),
				ptr:      new(int),
			},
			want:    &a,
			wantErr: false,
		},
		{
			name: "error error type",
			args: extractArgs{
				ptr: new(error),
			},
			want:           new(error),
			wantErr:        true,
			wantErrContain: "can't extract an error",
		},
		{
			name: "error not pointer",
			args: extractArgs{
				ptr: "abc",
			},
			want:           nil,
			wantErr:        true,
			wantErrContain: "ptr must is pointer",
		},
		{
			name:           "error nil",
			args:           extractArgs{},
			want:           nil,
			wantErr:        true,
			wantErrContain: "can't extract an untyped nil",
		},
	}
	for _, tt := range extractTests {
		t.Run(tt.name, func(t *testing.T) {
			c := dig.New()
			if tt.args.provider != nil {
				err := tt.args.provider.Apply(c.Provide)
				if err != nil {
					t.Errorf("provider.Apply err = %s", err)
					return
				}
			}
			extractFunc := MakeExtractFunc(tt.args.ptr, tt.args.opts...)
			err := c.Invoke(extractFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("c.Invoke(extractFunc) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.wantErrContain != "" {
					if !strings.Contains(err.Error(), tt.wantErrContain) {
						t.Errorf("c.Invoke(extractFunc) error want contain %s, got %s", tt.wantErrContain, err.Error())
						return
					}
				}
				return
			}
			if !reflect.DeepEqual(tt.args.ptr, tt.want) {
				t.Errorf("ptr = %#v, want %#v", tt.args.ptr, tt.want)
			}
		})
	}
}

func TestExtractWithLocationForPC(t *testing.T) {
	type extractArgs struct {
		provider tests.Provider
		typ      interface{}
		opts     []ExtractOption
	}
	var extractTests = []struct {
		name           string
		args           extractArgs
		want           interface{}
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "success as and named type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.Name("a"), dig.As(new(interface{}))),
				typ:      new(interface{}),
				opts:     []ExtractOption{ExtractOptionFunc(func(eo *ExtractOptions) { eo.Name = "a" })},
			},
			want:    interfaceA,
			wantErr: false,
		},
		{
			name: "success as type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.As(new(interface{}))),
				typ:      new(interface{}),
			},
			want:    interfaceA,
			wantErr: false,
		},
		{
			name: "success group type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.Group("g")),
				typ:      []int{},
				opts:     []ExtractOption{ExtractOptionFunc(func(eo *ExtractOptions) { eo.Group = "g" })},
			},
			want:    []int{1},
			wantErr: false,
		},
		{
			name: "success named type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }, dig.Name("a")),
				typ:      int(1),
				opts:     []ExtractOption{ExtractOptionFunc(func(eo *ExtractOptions) { eo.Name = "a" })},
			},
			want:    a,
			wantErr: false,
		},
		{
			name: "success interface pointer",
			args: extractArgs{
				provider: tests.ProviderOne(func() *interface{} { return interfaceAP }),
				typ:      new(*interface{}),
			},
			want:    interfaceAP,
			wantErr: false,
		},
		{
			name: "success interface",
			args: extractArgs{
				provider: tests.ProviderOne(func() interface{} { return interfaceA }),
				typ:      new(interface{}),
			},
			want:    interfaceA,
			wantErr: false,
		},
		{
			name: "success type pointer",
			args: extractArgs{
				provider: tests.ProviderOne(func() *int { return &a }),
				typ:      new(int),
			},
			want:    ap,
			wantErr: false,
		},
		{
			name: "success type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }),
				typ:      1,
			},
			want:    a,
			wantErr: false,
		},
		{
			name: "error error type",
			args: extractArgs{
				typ: new(error),
			},
			want:           new(error),
			wantErr:        true,
			wantErrContain: "can't extract an error",
		},
		{
			name:           "error nil",
			args:           extractArgs{},
			want:           nil,
			wantErr:        true,
			wantErrContain: "can't extract an untyped nil",
		},
		{
			name: "error fix errMissingDependencies caller",
			args: extractArgs{
				typ: "",
			},
			wantErr:        true,
			wantErrContain: tests.GetSelfSourceCodeFilePath(),
		},
	}

	for _, tt := range extractTests {
		t.Run(tt.name, func(t *testing.T) {
			c := dig.New()
			if tt.args.provider != nil {
				err := tt.args.provider.Apply(c.Provide)
				if err != nil {
					t.Errorf("provider.Apply err = %s", err)
					return
				}
			}
			got, err := ExtractWithLocationForPC(c.Invoke, 2, tt.args.typ, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractWithLocationForPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.wantErrContain != "" {
					if !strings.Contains(err.Error(), tt.wantErrContain) {
						t.Errorf("ExtractWithLocationForPC() error want contain %s, got %s", tt.wantErrContain, err.Error())
						return
					}
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ptr = %#v, want %#v", tt.args.typ, tt.want)
			}
		})
	}
}
