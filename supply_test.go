package digpro

import (
	"errors"
	"testing"

	"go.uber.org/dig"
)

type testSupplyArgs struct {
	prepare *_providerSet
	value   interface{}
}

var testSupplyData = []struct {
	name    string
	args    testSupplyArgs
	wantErr bool
}{
	{
		name: "error duplicate",
		args: testSupplyArgs{
			prepare: providerSet(
				provide(Supply(1)),
			),
			value: 1,
		},
		wantErr: true,
	},
	{
		name: "error give nil",
		args: testSupplyArgs{
			value: nil,
		},
		wantErr: false,
	},
	{
		name: "error give error",
		args: testSupplyArgs{
			value: errors.New("abc"),
		},
		wantErr: true,
	},
}

func TestSupply(t *testing.T) {
	for _, tt := range testSupplyData {
		t.Run(tt.name, func(t *testing.T) {
			c := dig.New()
			err := tt.args.prepare.apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error: %s", err)
				return
			}
			err = c.Provide(Supply(tt.args.value))
			if (err != nil) != tt.wantErr {
				t.Errorf("provider.apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
		})
	}
}

func TestContainerWrapper_Supply(t *testing.T) {
	for _, tt := range testSupplyData {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := tt.args.prepare.apply(c.Provide)
			if err != nil {
				t.Errorf("prepare error: %s", err)
				return
			}
			err = c.Supply(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("provider.apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
		})
	}
}
