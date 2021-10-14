package tests

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/dig"
)

type Provider interface {
	Apply(f ProviderApplyFunc) error
}

type provider struct {
	Constructor interface{}
	Opts        []dig.ProvideOption
}

func (p *provider) Apply(f ProviderApplyFunc) error {
	if p == nil {
		return nil
	}
	return f(p.Constructor, p.Opts...)
}

func ProviderOne(constructor interface{}, opts ...dig.ProvideOption) Provider {
	return &provider{
		Constructor: constructor,
		Opts:        opts,
	}
}

type ProviderApplyFunc func(constructor interface{}, opts ...dig.ProvideOption) error

type providerSet struct {
	providers []Provider
}

func ProviderSet(providers ...Provider) Provider {
	return &providerSet{
		providers: providers,
	}
}

func (ps *providerSet) Apply(f ProviderApplyFunc) error {
	if ps == nil {
		return nil
	}
	errs := []string{}
	for i, p := range ps.providers {
		err := p.Apply(f)
		if err != nil {
			errs = append(errs, fmt.Sprintf("[%d] %s", i, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "\n"))
}

func GetSelfSourceCodeFilePath() string {
	_, fpath, _, ok := runtime.Caller(1)
	if !ok {
		err := errors.New("failed to get filename")
		panic(err)
	}
	// return path.Base(fpath)
	return fpath
}
