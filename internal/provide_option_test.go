package internal

import (
	"reflect"
	"testing"

	"go.uber.org/dig"
)

func TestApplyProvideOptions(t *testing.T) {
	i := new(interface{})
	info := dig.ProvideInfo{}
	type args struct {
		opts []dig.ProvideOption
	}
	tests := []struct {
		name string
		args args
		want *ProvideOptions
	}{
		{
			name: "type",
			args: args{
				opts: []dig.ProvideOption{
					dig.Name("a"), dig.Name("aa"),
					dig.As(i), dig.As(i),
					dig.Group("b"), dig.Group("bb"),
					dig.FillProvideInfo(&info)},
			},
			want: &ProvideOptions{
				Name:  "aa",
				Group: "bb",
				Info:  &info,
				As:    []interface{}{i, i},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyProvideOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyProvideOptions() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
