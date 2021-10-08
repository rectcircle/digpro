package digpro

import (
	"reflect"
	"strings"
	"testing"

	"go.uber.org/dig"
)

var fooNilPtr *Foo

type testStructArgs struct {
	prepare           *_providerSet
	structOrStructPtr interface{}
	opts              []dig.ProvideOption
}

var testStructData = []struct {
	name                  string
	args                  testStructArgs
	wantErr               bool
	wantErrContain        string
	want                  interface{}
	wantOpts              []ExtractOption
	wantExtractErr        bool
	wantExtractErrContain string
}{
	{
		name: "error nil",
		args: testStructArgs{
			structOrStructPtr: nil,
		},
		wantErr:        true,
		wantErrContain: "<nil>",
	},
	{
		name: "error nil struct",
		args: testStructArgs{
			structOrStructPtr: fooNilPtr,
		},
		wantErr:        true,
		wantErrContain: "(nil)",
	},
	{
		name: "error string",
		args: testStructArgs{
			structOrStructPtr: "string",
		},
		wantErr:        true,
		wantErrContain: "but got string",
	},
	{
		name: "error int ptr",
		args: testStructArgs{
			structOrStructPtr: new(int),
		},
		wantErr:        true,
		wantErrContain: "but got *int",
	},
	{
		name: "error field conflict",
		args: testStructArgs{
			structOrStructPtr: struct {
				A string
				a string
			}{},
		},
		wantErr:        true,
		wantErrContain: "field conflict",
	},
	{
		name: "error missing dependencies",
		args: testStructArgs{
			structOrStructPtr: Biz{},
		},
		wantErr:               false,
		want:                  Biz{},
		wantExtractErr:        true,
		wantExtractErrContain: "missing dependencies",
	},
	{
		name: "success struct ptr",
		args: testStructArgs{
			prepare: providerSet(
				provide(Supply("a")),
				provide(Supply(1)),
				provide(Supply(true)),
			),
			structOrStructPtr: new(Bar),
		},
		wantErr: false,
		want: &Bar{
			A:       "a",
			B:       1,
			private: true,
		},
	},
	{
		name: "success struct",
		args: testStructArgs{
			prepare: providerSet(
				provide(Supply("a")),
				provide(Supply(1)),
				provide(Supply(true)),
			),
			structOrStructPtr: Bar{},
		},
		wantErr: false,
		want: Bar{
			A:       "a",
			B:       1,
			private: true,
		},
	},
	{
		name: "tag",
		args: testStructArgs{
			prepare: providerSet(
				provide(Supply(1), dig.Name("a")),
				provide(Supply("c"), dig.Group("c")),
				provide(Supply("c"), dig.Group("c")),
			),
			structOrStructPtr: Biz{
				A: 0,
				B: 3,
				C: []string{},
			},
			opts: []dig.ProvideOption{},
		},
		wantErr: false,
		want: Biz{
			A: 1,
			B: 3,
			C: []string{"c", "c"},
		},
	},
}

func TestStruct(t *testing.T) {
	for _, tt := range testStructData {
		t.Run(tt.name, func(t *testing.T) {
			c := dig.New()

			type Foo struct {
				A       string
				B       int
				private bool
			}
			c.Provide(func(in struct {
				dig.In
				A       string
				B       int
				Private bool
			}) Foo {
				return Foo{
					A:       in.A,
					B:       in.B,
					private: in.Private,
				}
			})

			err := tt.args.prepare.apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error = %v", err)
				return
			}
			err = c.Provide(Struct(tt.args.structOrStructPtr), tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("c.Provide(Struct(structOrStructPtr), opts...) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				// fmt.Println(err)
				if !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("c.Provide(Struct(structOrStructPtr), opts...) error = %v, want contain = %s", err, tt.wantErrContain)
				}
				return
			}
			got, err := Extract(c, tt.want, tt.wantOpts...)
			if (err != nil) != tt.wantExtractErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantExtractErr)
				return
			}
			if err != nil {
				// fmt.Println(err)
				if !strings.Contains(err.Error(), tt.wantExtractErrContain) {
					t.Errorf("Extract() error = %v, want contain = %s", err, tt.wantExtractErrContain)
				}
				return
			}
			if !reflect.DeepEqual(got, interface{}(tt.want)) {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainerWrapper_Struct(t *testing.T) {
	for _, tt := range testStructData {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := tt.args.prepare.apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error = %v", err)
				return
			}
			err = c.Struct(tt.args.structOrStructPtr, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerWrapper.Struct error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			got, err := c.Extract(tt.want, tt.wantOpts...)
			if (err != nil) != tt.wantExtractErr {
				t.Errorf("ContainerWrapper.Extract() error = %v, wantErr %v", err, tt.wantExtractErr)
				return
			}
			if err != nil {
				if !strings.Contains(err.Error(), tt.wantExtractErrContain) {
					t.Errorf("ContainerWrapper.Extract() error = %v, want contain = %s", err, tt.wantExtractErrContain)
				}
				// if contain /usr/local/Cellar/go/1.17.1/libexec/src/reflect/asm_amd64.s:30, failed. because of c.Struct bug
				if strings.Contains(err.Error(), "makeFuncStub") {
					t.Errorf("ContainerWrapper.Extract() error = %v, contain = %s", err, "makeFuncStub")
				}
				return
			}
			if !reflect.DeepEqual(got, interface{}(tt.want)) {
				t.Errorf("ContainerWrapper.Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}
