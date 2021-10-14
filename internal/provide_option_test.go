package internal

import (
	"reflect"
	"testing"

	"go.uber.org/dig"
)

func TestApplyProvideOptions(t *testing.T) {
	info := dig.ProvideInfo{}
	type args struct {
		opts []dig.ProvideOption
	}
	tests := []struct {
		name string
		args args
		want ProvideOptions
	}{
		{
			name: "type",
			args: args{
				opts: []dig.ProvideOption{
					dig.Name("a"), dig.Name("aa"),
					dig.As(new(interface{})), dig.As(new(interface{})),
					dig.Group("b"), dig.Group("bb"),
					dig.FillProvideInfo(&info)},
			},
			want: ProvideOptions{
				Name:  "aa",
				Group: "bb",
				Info:  &info,
				As:    []interface{}{new(interface{}), new(interface{})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyProvideOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyProvideOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
