package digpro

import (
	"reflect"
	"strings"
	"testing"

	"go.uber.org/dig"
)

// TODO test error 错误输出问题

var i = 1

type testExtractArgs struct {
	prepare      *_providerSet
	typInterface interface{}
	opts         []ExtractOption
}

var testExtractData = []struct {
	name           string
	args           testExtractArgs
	want           interface{}
	wantErr        bool
	wantErrContain string
}{
	{
		name: "int(1)",
		args: testExtractArgs{
			prepare: providerSet(
				provide(Supply(1)),
			),
			typInterface: 0,
			opts:         []ExtractOption{},
		},
		want:    1,
		wantErr: false,
	},
	{
		name: "*int(1)",
		args: testExtractArgs{
			prepare: providerSet(
				provide(Supply(&i)),
			),
			typInterface: &i,
			opts:         []ExtractOption{},
		},
		want:    &i,
		wantErr: false,
	},
	{
		name: "name",
		args: testExtractArgs{
			prepare: providerSet(
				provide(Supply(i), dig.Name("i")),
			),
			typInterface: i,
			opts:         []ExtractOption{ExtractByName("i")},
		},
		want:    i,
		wantErr: false,
	},
	{
		name: "group",
		args: testExtractArgs{
			prepare: providerSet(
				provide(Supply(1), dig.Group("i")),
				provide(Supply(1), dig.Group("i")),
			),
			typInterface: []int{},
			opts:         []ExtractOption{ExtractByGroup("i")},
		},
		want:    []int{1, 1},
		wantErr: false,
	},
	{
		name: "error missing dependencies 1",
		args: testExtractArgs{
			prepare:      providerSet(),
			typInterface: 1,
			opts:         []ExtractOption{},
		},
		wantErr:        true,
		wantErrContain: getSelfSourceCodeFilePath(),
	},
	{
		name: "error missing dependencies 2",
		args: testExtractArgs{
			prepare:      providerSet(provide(Struct(Biz{}))),
			typInterface: Biz{},
			opts:         []ExtractOption{},
		},
		wantErr:        true,
		wantErrContain: getSelfSourceCodeFilePath(),
	},
}

func TestExtract(t *testing.T) {
	for _, tt := range testExtractData {
		t.Run(tt.name, func(t *testing.T) {
			c := dig.New()
			err := tt.args.prepare.apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error = %v", err)
				return
			}
			got, err := Extract(c, tt.args.typInterface, tt.args.opts...)
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
			err := tt.args.prepare.apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error = %v", err)
				return
			}
			got, err := c.Extract(tt.args.typInterface, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerWrapper.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				// fmt.Println(err)
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
