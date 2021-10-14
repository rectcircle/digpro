package internal

import (
	"strings"
	"testing"

	"github.com/rectcircle/digpro/internal/tests"
	"go.uber.org/dig"
)

func TestProvideWithLocationForPC(t *testing.T) {
	type args struct {
		c           *dig.Container
		callSkip    int
		constructor interface{}
		opts        []dig.ProvideOption
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "error with right caller info",
			args: args{
				c:        dig.New(),
				callSkip: 2,
				constructor: func() {
				},
				opts: []dig.ProvideOption{},
			},
			wantErr:        true,
			wantErrContain: tests.GetSelfSourceCodeFilePath(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProvideWithLocationForPC(tt.args.c.Provide, tt.args.callSkip, tt.args.constructor, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProvideWithLocationForPC() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if tt.wantErrContain != "" {
					if !strings.Contains(err.Error(), tt.wantErrContain) {
						t.Errorf("ProvideWithLocationForPC() error want contain %s, got %s", tt.wantErrContain, err.Error())
						return
					}
				}
			}
		})
	}
}
