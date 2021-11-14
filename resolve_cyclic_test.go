package digpro

import (
	"testing"

	"github.com/rectcircle/digpro/internal/digcopy"
	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

func TestResolveCyclic(t *testing.T) {

	type args struct {
		prepare PrepareFunc
	}
	tests := []struct {
		name           string
		args           args
		prepareWantErr bool
		assert         func(t *testing.T, c *ContainerWrapper)
	}{
		{
			name: "circular dependency error lower api",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					return tests.ProviderSet(
						tests.ProviderOne(Supply(1)),
						tests.ProviderOne(Supply("a")),
						tests.ProviderOne(Struct(new(D1))),
						tests.ProviderOne(Struct(new(D2))),
					).Apply(c.Provide)
				},
			},
			prepareWantErr: true,
		},
		{
			name: "circular dependency error high api",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply(1),
						c.Supply("a"),
						c.Struct(new(D1)),
						c.Struct(new(D2)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: true,
		},
		{
			name: "ResolveCyclic first get *D1",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply(1),
						c.Supply("a"),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &D1{Value: 1}
				d2 := &D2{Value: "a"}
				d1.D2 = d2
				d2.D1 = d1
				for i := 0; i < 2; i++ {
					targetD1, err := c.Extract(new(D1))
					if err != nil {
						t.Errorf("[%d] Extract *D1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(*D1).String() {
						t.Errorf("[%d] Expected *D1 %s, got %s", i, d1.String(), targetD1.(*D1).String())
						return
					}
					targetD2, err := c.Extract(new(D2))
					if err != nil {
						t.Errorf("[%d] Extract *D2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(*D2).String() {
						t.Errorf("[%d] Expected *D2 %s, got %s", i, d2.String(), targetD2.(*D2).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic first get *D2",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply(1),
						c.Supply("a"),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &D1{Value: 1}
				d2 := &D2{Value: "a"}
				d1.D2 = d2
				d2.D1 = d1
				for i := 0; i < 2; i++ {
					targetD2, err := c.Extract(new(D2))
					if err != nil {
						t.Errorf("[%d] Extract *D2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(*D2).String() {
						t.Errorf("[%d] Expected *D2 %s, got %s", i, d2.String(), targetD2.(*D2).String())
						return
					}
					targetD1, err := c.Extract(new(D1))
					if err != nil {
						t.Errorf("[%d] Extract *D1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(*D1).String() {
						t.Errorf("[%d] Expected *D1 %s, got %s", i, d1.String(), targetD1.(*D1).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic first get D3",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply(1),
						c.Supply("a"),
						c.Supply(true),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
						c.Struct(D3{}),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &D1{Value: 1}
				d2 := &D2{Value: "a"}
				d3 := D3{Value: true}
				d1.D2 = d2
				d2.D1 = d1
				d3.D2 = d2
				for i := 0; i < 2; i++ {
					targetD3, err := c.Extract(D3{})
					if err != nil {
						t.Errorf("[%d] Extract D3 error: %v", i, err)
						return
					}
					if d3.String() != targetD3.(D3).String() {
						t.Errorf("[%d] Expected D3 %s, got %s", i, d3.String(), targetD3.(D3).String())
						return
					}
					targetD2, err := c.Extract(new(D2))
					if err != nil {
						t.Errorf("[%d] Extract *D2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(*D2).String() {
						t.Errorf("[%d] Expected *D2 %s, got %s", i, d2.String(), targetD2.(*D2).String())
						return
					}
					targetD1, err := c.Extract(new(D1))
					if err != nil {
						t.Errorf("[%d] Extract *D1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(*D1).String() {
						t.Errorf("[%d] Expected *D1 %s, got %s", i, d1.String(), targetD1.(*D1).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic error miss dependencies one level",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply("a"),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				_, err := c.Extract(new(D1))
				if err == nil {
					t.Errorf("want error but got %v", err)
					return
				}
				if _, ok := err.(digcopy.ErrArgumentsFailed); !ok {
					t.Errorf("want ErrArgumentsFailed but got %+v", err)
				}
			},
		},
		{
			name: "ResolveCyclic error miss dependencies two level",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply("a"),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				_, err := c.Extract(new(D2))
				if err == nil {
					t.Errorf("want error but got %v", err)
					return
				}
				if _, ok := err.(digcopy.ErrArgumentsFailed); !ok {
					t.Errorf("want ErrArgumentsFailed but got %+v", err)
				}
			},
		},
		{
			name: "ResolveCyclic with optional ignore 1 - has dependency",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply("a"),
						c.Supply(true),
						c.Supply(float64(1.1)),
						c.Supply(1),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
						c.Struct(new(D4)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &D1{Value: 1}
				d2 := &D2{Value: "a"}
				d4 := &D4{Value: 1.1, Ignore: 0}
				d1.D2 = d2
				d2.D1 = d1
				d4.D2 = d2
				for i := 0; i < 2; i++ {
					targetD4, err := c.Extract(new(D4))
					if err != nil {
						t.Errorf("[%d] Extract *D4 error: %v", i, err)
						return
					}
					if d4.String() != targetD4.(*D4).String() {
						t.Errorf("[%d] Expected D3 %s, got %s", i, d4.String(), targetD4.(*D4).String())
						return
					}
					targetD2, err := c.Extract(new(D2))
					if err != nil {
						t.Errorf("[%d] Extract *D2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(*D2).String() {
						t.Errorf("[%d] Expected *D2 %s, got %s", i, d2.String(), targetD2.(*D2).String())
						return
					}
					targetD1, err := c.Extract(new(D1))
					if err != nil {
						t.Errorf("[%d] Extract *D1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(*D1).String() {
						t.Errorf("[%d] Expected *D1 %s, got %s", i, d1.String(), targetD1.(*D1).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic with optional ignore 2 - self use ResolveCyclic",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply("a"),
						c.Supply(true),
						c.Supply(float64(1.1)),
						c.Supply(1),
						c.Struct(new(D1), ResolveCyclic()),
						c.Struct(new(D2)),
						c.Struct(new(D4), ResolveCyclic()),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &D1{Value: 1}
				d2 := &D2{Value: "a"}
				d4 := &D4{Value: 1.1, Ignore: 0}
				d1.D2 = d2
				d2.D1 = d1
				d4.D2 = d2
				for i := 0; i < 2; i++ {
					targetD4, err := c.Extract(new(D4))
					if err != nil {
						t.Errorf("[%d] Extract *D4 error: %v", i, err)
						return
					}
					if d4.String() != targetD4.(*D4).String() {
						t.Errorf("[%d] Expected D3 %s, got %s", i, d4.String(), targetD4.(*D4).String())
						return
					}
					targetD2, err := c.Extract(new(D2))
					if err != nil {
						t.Errorf("[%d] Extract *D2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(*D2).String() {
						t.Errorf("[%d] Expected *D2 %s, got %s", i, d2.String(), targetD2.(*D2).String())
						return
					}
					targetD1, err := c.Extract(new(D1))
					if err != nil {
						t.Errorf("[%d] Extract *D1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(*D1).String() {
						t.Errorf("[%d] Expected *D1 %s, got %s", i, d1.String(), targetD1.(*D1).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic with optional ignore 3 - no dependency",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Struct(new(D4)),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d4 := &D4{}
				for i := 0; i < 2; i++ {
					targetD4, err := c.Extract(new(D4))
					if err != nil {
						t.Errorf("[%d] Extract *D4 error: %v", i, err)
						return
					}
					if d4.String() != targetD4.(*D4).String() {
						t.Errorf("[%d] Expected D3 %s, got %s", i, d4.String(), targetD4.(*D4).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic with as",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply(1),
						c.Supply("a"),
						c.Struct(new(DI1), ResolveCyclic(), dig.As(new(I1))),
						c.Struct(new(DI2), dig.As(new(I2))),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &DI1{Value: 1}
				d2 := &DI2{Value: "a"}
				d1.I2 = d2
				d2.I1 = d1
				for i := 0; i < 2; i++ {
					targetD1, err := c.Extract(new(I1))
					if err != nil {
						t.Errorf("[%d] Extract *DI1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(I1).String() {
						t.Errorf("[%d] Expected *DI1 %s, got %s", i, d1.String(), targetD1.(I1).String())
						return
					}
					targetD2, err := c.Extract(new(I2))
					if err != nil {
						t.Errorf("[%d] Extract *DI2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(I2).String() {
						t.Errorf("[%d] Expected *DI2 %s, got %s", i, d2.String(), targetD2.(I2).String())
						return
					}
				}
			},
		},
		{
			name: "ResolveCyclic with as and override",
			args: args{
				prepare: func(c *ContainerWrapper) error {
					errs := []error{
						c.Supply(1),
						c.Supply("a"),
						c.Struct(new(DI1), ResolveCyclic(), dig.As(new(I1))),
						c.Struct(new(DI2), dig.As(new(I2))),
						c.Supply(2, Override()),
						c.Supply("b", Override()),
						c.Struct(new(DI1), ResolveCyclic(), dig.As(new(I1)), Override()),
					}
					for _, err := range errs {
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			prepareWantErr: false,
			assert: func(t *testing.T, c *ContainerWrapper) {
				d1 := &DI1{Value: 2}
				d2 := &DI2{Value: "b"}
				d1.I2 = d2
				d2.I1 = d1
				for i := 0; i < 2; i++ {
					targetD1, err := c.Extract(new(I1))
					if err != nil {
						t.Errorf("[%d] Extract *DI1 error: %v", i, err)
						return
					}
					if d1.String() != targetD1.(I1).String() {
						t.Errorf("[%d] Expected *DI1 %s, got %s", i, d1.String(), targetD1.(I1).String())
						return
					}
					targetD2, err := c.Extract(new(I2))
					if err != nil {
						t.Errorf("[%d] Extract *DI2 error: %v", i, err)
						return
					}
					if d2.String() != targetD2.(I2).String() {
						t.Errorf("[%d] Expected *DI2 %s, got %s", i, d2.String(), targetD2.(I2).String())
						return
					}
				}
			},
		},
		// TODO 添加 Example
		// TODO 添加 docs
		// TODO 添加文档
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := tt.args.prepare(c)
			if (err != nil) != tt.prepareWantErr {
				t.Errorf("prepare() error = %v, wantErr %v", err, tt.prepareWantErr)
				return
			}
			if err != nil {
				return
			}
			tt.assert(t, c)
		})
	}
}
