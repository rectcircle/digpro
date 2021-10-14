package internal

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

func TestTryFixDigErr(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	commonErr := errors.New("abc")
	type args struct {
		err error
		pc  uintptr
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantErrType      string
		wantErrContain   string
		wantErrRootCause error
	}{
		{
			name: "nil",
			args: args{
				err: nil,
				pc:  0,
			},
			wantErr: false,
		},
		{
			name: "common error",
			args: args{
				err: commonErr,
			},
			wantErr:     true,
			wantErrType: reflect.TypeOf(commonErr).String(),
		},
		{
			name: "common error with pc",
			args: args{
				err: commonErr,
				pc:  pc,
			},
			wantErr:     true,
			wantErrType: reflect.TypeOf(commonErr).String(),
		},
		{
			name: "dig.errProvide return void",
			args: args{
				err: func() error {
					c := dig.New()
					return c.Provide(func() {})
				}(),
				pc: pc,
			},
			wantErrType: "dig.errProvide",
			wantErr:     true,
		},
		{
			name: "dig.errProvide with pc",
			args: args{
				err: func() error {
					c := dig.New()
					c.Provide(func() int { return 1 })
					return c.Provide(func() int { return 1 })
				}(),
				pc: pc,
			},
			wantErr:        true,
			wantErrType:    "dig.errProvide",
			wantErrContain: tests.GetSelfSourceCodeFilePath(),
		},
		{
			name: "dig.errArgumentsFailed with pc",
			args: args{
				err: func() error {
					c := dig.New()
					c.Provide(func() (int, error) { return 0, commonErr })
					return c.Invoke(func(a int) {})
				}(),
				pc: pc,
			},
			wantErr:          true,
			wantErrType:      "dig.errArgumentsFailed",
			wantErrContain:   tests.GetSelfSourceCodeFilePath(),
			wantErrRootCause: commonErr,
		},
		{
			name: "dig.errMissingDependencies with pc",
			args: args{
				err: func() error {
					c := dig.New()
					return c.Invoke(func(a int) {})
				}(),
				pc: pc,
			},
			wantErr:        true,
			wantErrType:    "dig.errMissingDependencies",
			wantErrContain: tests.GetSelfSourceCodeFilePath(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TryFixDigErr(tt.args.err, tt.args.pc)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryConvertDigErr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.wantErrContain != "" {
					if !strings.Contains(err.Error(), tt.wantErrContain) {
						t.Errorf("TryConvertDigErr() error want contain %s, got %s", tt.wantErrContain, err.Error())
						return
					}
				}
				if tt.wantErrType != "" {
					if reflect.TypeOf(err).String() != tt.wantErrType {
						t.Errorf("TryConvertDigErr() error want type %s, got %s", tt.wantErrType, reflect.TypeOf(err).String())
						return
					}
				}
				if tt.wantErrRootCause != nil {
					if !errors.Is(dig.RootCause(err), tt.wantErrRootCause) {
						t.Errorf("dig.RootCause(err) error want type %v, got %v", tt.wantErrRootCause, dig.RootCause(err))
						return
					}
				}
			}
		})
	}
}

func TestWrapErrorWithLocationForPC(t *testing.T) {
	type args struct {
		callSkip int
		f        func(pc uintptr) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not error",
			args: args{
				callSkip: 0,
				f: func(pc uintptr) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "want error",
			args: args{
				callSkip: 0,
				f: func(pc uintptr) error {
					return errors.New("abc")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapErrorWithLocationForPC(tt.args.callSkip, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("WrapErrorWithLocationForPC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
