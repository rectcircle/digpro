package digpro

import (
	"errors"
	"testing"
)

func TestQuickPanic(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name      string
		args      args
		wantPanic interface{}
	}{
		{
			name: "one nil",
			args: args{
				errs: []error{nil},
			},
		},
		{
			name: "some nil",
			args: args{
				errs: []error{nil, nil},
			},
		},
		{
			name: "one error",
			args: args{
				errs: []error{errors.New("abc")},
			},
			wantPanic: "[0]: abc",
		},
		{
			name: "all error",
			args: args{
				errs: []error{errors.New("abc"), errors.New("abc"), errors.New("abc")},
			},
			wantPanic: "[0]: abc\n[1]: abc\n[2]: abc",
		},
		{
			name: "some error",
			args: args{
				errs: []error{errors.New("abc"), nil, errors.New("abc")},
			},
			wantPanic: "[0]: abc\n[2]: abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if err != tt.wantPanic {
					t.Errorf("QuickPanic got = %v, want panic = %v", err, tt.wantPanic)
				}
			}()
			QuickPanic(tt.args.errs...)
		})
	}
}
