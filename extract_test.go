package digpro

import (
	"reflect"
	"strings"
	"testing"

	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

type testExtractArgs struct {
	prepare tests.Provider
	typ     interface{}
	opts    []ExtractOption
}

var testExtractData = []struct {
	name           string
	args           testExtractArgs
	want           interface{}
	wantErr        bool
	wantErrContain string
}{
	{
		name: "error missing dependencies",
		args: testExtractArgs{
			prepare: tests.ProviderSet(),
			typ:     1,
			opts:    []ExtractOption{},
		},
		wantErr:        true,
		wantErrContain: tests.GetSelfSourceCodeFilePath(),
	},
	{
		name: "success name",
		args: testExtractArgs{
			prepare: tests.ProviderOne(Supply(1), dig.Name("a")),
			typ:     1,
			opts:    []ExtractOption{ExtractByName("a")},
		},
		want: 1,
	},
	{
		name: "success group",
		args: testExtractArgs{
			prepare: tests.ProviderOne(Supply(1), dig.Group("g")),
			typ:     []int{},
			opts:    []ExtractOption{ExtractByGroup("g")},
		},
		want: []int{1},
	},
}

func TestExtract(t *testing.T) {
	for _, tt := range testExtractData {
		t.Run(tt.name, func(t *testing.T) {
			c := dig.New()
			err := tt.args.prepare.Apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error = %v", err)
				return
			}
			got, err := Extract(c, tt.args.typ, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				// fmt.Println(err)
				if !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("Extract() error = %v, want contain = %s", err, tt.wantErrContain)
				}
				return
			}
			if !reflect.DeepEqual(got, interface{}(tt.want)) {
				t.Errorf("Extract() = %v, want %v", got, interface{}(tt.want))
			}
		})
	}
}

func TestContainerWrapper_Extract(t *testing.T) {
	for _, tt := range testExtractData {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := tt.args.prepare.Apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error = %v", err)
				return
			}
			got, err := c.Extract(tt.args.typ, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerWrapper.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("ContainerWrapper.Extract() error = %v, want contain = %s", err, tt.wantErrContain)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContainerWrapper.Extract() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestMakeExtractFunc(t *testing.T) {
	a := 1
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
			name: "success type",
			args: extractArgs{
				provider: tests.ProviderOne(func() int { return 1 }),
				ptr:      new(int),
			},
			want:    &a,
			wantErr: false,
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
