package digpro

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/rectcircle/digpro/internal"
	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

func TestNew(t *testing.T) {
	c := New()
	errs := []error{
		c.Supply("a"),
		c.Supply(1),
		c.Supply(true),
		c.Struct(new(Bar)),
	}
	for _, err := range errs {
		if err != nil {
			t.Errorf("c.Value() / c.Struct() error = %v, wantErr %v", err, false)
		}
	}
	err := c.Invoke(func(bar *Bar) {
		if bar.A != "a" {
			t.Errorf("bar.A got = %s, want = %s", bar.A, "a")
		}
		if bar.B != 1 {
			t.Errorf("bar.A got = %d, want = %d", bar.B, 1)
		}
		if bar.private != true {
			t.Errorf("bar.A got = %t, want = %t", bar.private, true)
		}
	})
	if err != nil {
		t.Errorf("c.Invoke() error = %v, wantErr %v", err, false)
	}
}

func TestContainerWrapper_Unwrap(t *testing.T) {
	c := New()
	if got := c.Unwrap(); !reflect.DeepEqual(got, &c.Container) {
		t.Errorf("ContainerWrapper.Unwrap() = %v, want %v", got, &c.Container)
		return
	}
}

func TestContainerWrapper_Visualize(t *testing.T) {
	c := New()
	if err := c.Visualize(bytes.NewBuffer(nil)); (err != nil) != false {
		t.Errorf("ContainerWrapper.Visualize() error = %v, wantErr %v", err, false)
		return
	}
}

func Test_provideMiddleware(t *testing.T) {
	type args struct {
		provideMiddlewares []provideMiddleware
		providerSet        tests.Provider
		extractOptions     []ExtractOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "no use middlewares",
			args: args{
				provideMiddlewares: []provideMiddleware{},
				providerSet:        tests.ProviderSet(tests.ProviderOne(Supply(1))),
			},
			wantErr: false,
			want:    1,
		},
		{
			name: "one middlewares call next",
			args: args{
				provideMiddlewares: []provideMiddleware{func(pc *provideContext) error {
					return pc.next()
				}},
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
			},
			wantErr: false,
			want:    1,
		},
		{
			name: "one middlewares not call next",
			args: args{
				provideMiddlewares: []provideMiddleware{func(pc *provideContext) error {
					return nil
				}},
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
			},
			wantErr: false,
			want:    1,
		},
		{
			name: "one middlewares error",
			args: args{
				provideMiddlewares: []provideMiddleware{func(pc *provideContext) error {
					return errors.New("middlewares return error")
				}},
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
			},
			wantErr: true,
		},
		{
			name: "two middlewares change constructor",
			args: args{
				provideMiddlewares: []provideMiddleware{
					func(pc *provideContext) error {
						return nil
					},
					func(pc *provideContext) error {
						pc.constructor = Supply(2)
						return nil
					},
				},
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
			},
			wantErr: false,
			want:    2,
		},
		{
			name: "two middlewares change constructor and add options",
			args: args{
				provideMiddlewares: []provideMiddleware{
					func(pc *provideContext) error {
						pc.opts = append(pc.opts, dig.Name("a"))
						return nil
					},
					func(pc *provideContext) error {
						pc.constructor = Supply(2)
						return nil
					},
				},
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
				extractOptions: []ExtractOption{
					ExtractByName("a"),
				},
			},
			wantErr: false,
			want:    2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContainerWrapper{
				Container:   *dig.New(),
				middlewares: tt.args.provideMiddlewares,
			}
			err := tt.args.providerSet.Apply(c.Provide)
			if (err != nil) != tt.wantErr {
				t.Errorf("c.Provide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got, err := c.Extract(tt.want, tt.args.extractOptions...); err != nil {
				t.Errorf("c.Extract() error = %v", err)
				return
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("c.Extract() got = %#v, want = %#v", got, tt.want)
			}
		})
	}
}

func TestContainerWrapper_provideInfos(t *testing.T) {
	info := dig.ProvideInfo{}
	type args struct {
		providerSet tests.Provider
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    [][]internal.ProvideOutput
	}{
		{
			name: "one only type",
			args: args{
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
			},
			wantErr: false,
			want: [][]internal.ProvideOutput{
				{
					{
						Type:  reflect.TypeOf(1),
						Name:  "",
						Group: "",
					},
				},
			},
		},
		{
			name: "one only type with FillProvideInfo",
			args: args{
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1), dig.FillProvideInfo(&info))),
			},
			wantErr: false,
			want: [][]internal.ProvideOutput{
				{
					{
						Type:  reflect.TypeOf(1),
						Name:  "",
						Group: "",
					},
				},
			},
		},
		{
			name: "one type with name",
			args: args{
				providerSet: tests.ProviderSet(tests.ProviderOne(Supply(1), dig.Name("a"))),
			},
			wantErr: false,
			want: [][]internal.ProvideOutput{
				{
					{
						Type:  reflect.TypeOf(1),
						Name:  "a",
						Group: "",
					},
				},
			},
		},
		{
			name: "two provide",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(func() (int, bool, error) {
						return 1, true, nil
					}, dig.Name("a")),
					tests.ProviderOne(Supply("string"), dig.Name("b")),
				),
			},
			wantErr: false,
			want: [][]internal.ProvideOutput{
				{
					{
						Type:  reflect.TypeOf(1),
						Name:  "a",
						Group: "",
					},
					{
						Type:  reflect.TypeOf(true),
						Name:  "a",
						Group: "",
					},
				},
				{
					{
						Type:  reflect.TypeOf("string"),
						Name:  "b",
						Group: "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := tt.args.providerSet.Apply(c.Provide)
			if (err != nil) != tt.wantErr {
				t.Errorf("c.Provide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(c.provideInfos) != len(tt.want) {
				t.Errorf("c.ProvideInfos len = %v, len(want) = %v", len(c.provideInfos), len(tt.want))
				return
			}
			for i, outputs := range tt.want {
				gotOutputs := c.provideInfos[i].ExportedOutputs()
				if len(gotOutputs) != len(outputs) {
					t.Errorf("c.ProvideInfos[%d].ExportedOutputs() len = %v, len(want[%d]) = %v", i, len(gotOutputs), i, len(outputs))
					return
				}
				for j, output := range outputs {
					gotOutput := gotOutputs[j]
					if !reflect.DeepEqual(output, gotOutput) {
						t.Errorf("c.ProvideInfos[%d].ExportedOutputs()[%d] len = %v, len(want[%d][%d]) = %v", i, j, output, i, j, gotOutput)
						return
					}
				}
			}
		})
	}
}

func TestContainerWrapper_Invoke(t *testing.T) {
	type args struct {
		function interface{}
		opts     []dig.InvokeOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "function is nil",
			args: args{
				function: nil,

				opts: []dig.InvokeOption{},
			},
			wantErr: true,
		},
		{
			name: "function is not function",
			args: args{
				function: 1,
				opts:     []dig.InvokeOption{},
			},
			wantErr: true,
		},
		{
			name: "function return error",
			args: args{
				function: func(a string) error {
					return errors.New("error")
				},
				opts: []dig.InvokeOption{},
			},
			wantErr: true,
		},
		{
			name: "function return not error",
			args: args{
				function: func(a string) string {
					return a
				},
				opts: []dig.InvokeOption{},
			},
			wantErr: false,
		},
		{
			name: "invoke missing type error",
			args: args{
				function: func(a bool) {},
				opts:     []dig.InvokeOption{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			_ = c.Supply(1) // please handle error in production
			_ = c.Supply("a")
			_ = c.Struct(new(D1), ResolveCyclic()) // enable resolve cyclic dependency
			_ = c.Struct(new(D2))
			if err := c.Invoke(tt.args.function, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("ContainerWrapper.Provide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
