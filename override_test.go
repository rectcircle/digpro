package digpro

import (
	"reflect"
	"testing"

	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

func TestOverride(t *testing.T) {
	type args struct {
		beforeProviderSet tests.Provider
		beforeExtract     interface{}
		providerSet       tests.Provider
		extractOptions    []ExtractOption
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		forceAssertWant bool
		want            interface{}
	}{
		{
			name: "error override and group",
			args: args{
				providerSet:    tests.ProviderSet(tests.ProviderOne(Supply(1), dig.Group("a"), Override())),
				extractOptions: []ExtractOption{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error not override conflict",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(Supply(1)),
					tests.ProviderOne(Supply(2)),
				),
				extractOptions: []ExtractOption{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error no provider to override was found 1",
			args: args{
				providerSet:    tests.ProviderSet(tests.ProviderOne(Supply(1), Override())),
				extractOptions: []ExtractOption{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error no provider to override was found 2",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(Supply(1)),
					tests.ProviderOne(func() (int, string) {
						return 2, "a"
					}, Override()),
				),
				extractOptions: []ExtractOption{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "success two return",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(func() (int, string) {
						return 1, "a"
					}),
					tests.ProviderOne(func() (int, string) {
						return 2, "b"
					}, Override()),
				),
				extractOptions: []ExtractOption{},
			},
			wantErr: false,
			want:    2,
		},
		{
			name: "error outputs has more than two registered provider",
			args: args{
				beforeProviderSet: tests.ProviderSet(tests.ProviderOne(Supply(1))),
				beforeExtract:     1,
				providerSet:       tests.ProviderSet(tests.ProviderOne(Supply(1), Override())),
				extractOptions:    []ExtractOption{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "success override many a time",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(Supply(1)),
					tests.ProviderOne(Supply(2), Override()),
					tests.ProviderOne(Supply(3), Override()),
				),
				extractOptions: []ExtractOption{},
			},
			wantErr: false,
			want:    3,
		},
		{
			name: "success has other provider",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(Supply(true)),
					tests.ProviderOne(Supply(1)),
					tests.ProviderOne(Supply(2), Override()),
					tests.ProviderOne(Supply(3), Override()),
				),
				extractOptions: []ExtractOption{},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "error recover",
			args: args{
				providerSet: tests.ProviderSet(
					tests.ProviderOne(func(a int) bool {
						return true
					}),
					tests.ProviderOne(func() int {
						return 1
					}),
					tests.ProviderOne(func(b bool) int {
						return 2
					}, Override()),
				),
				extractOptions: []ExtractOption{},
			},
			wantErr:         true,
			forceAssertWant: true,
			want:            1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			if tt.args.beforeProviderSet != nil {
				err := tt.args.beforeProviderSet.Apply(c.Provide)
				if err != nil {
					t.Errorf("before c.Provide() error = %v", err)
					return
				}
			}
			if tt.args.beforeExtract != nil {
				_, err := c.Extract(tt.args.beforeExtract)
				if err != nil {
					t.Errorf("before c.Extract() error = %v", err)
					return
				}
			}
			err := tt.args.providerSet.Apply(c.Provide)
			if (err != nil) != tt.wantErr {
				t.Errorf("c.Provide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !tt.forceAssertWant {
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
