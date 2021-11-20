package digpro

import (
	"reflect"

	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

type overrideProvideOption struct {
	dig.ProvideOption
}

type resolveCyclicProvideOption struct {
	dig.ProvideOption
}

type digproProvideOptions struct {
	enableOverride      bool
	enableResolveCyclic bool
	locationFixCallSkip int
}

var overrideProvideOptionType = reflect.TypeOf(overrideProvideOption{})
var resolveCyclicProvideOptionType = reflect.TypeOf(resolveCyclicProvideOption{})
var locationFixOptionType = reflect.TypeOf(internal.LocationFixOption{})

var digproProvideOptionTypeEnum = []reflect.Type{
	overrideProvideOptionType,
	resolveCyclicProvideOptionType,
	locationFixOptionType,
}

func filterAndGetDigproProvideOptions(opts []dig.ProvideOption, excludes ...reflect.Type) ([]dig.ProvideOption, digproProvideOptions) {
	filteredOpts := make([]dig.ProvideOption, 0, len(opts))
	result := digproProvideOptions{}
	for _, opt := range opts {
		isExclude := false
		for _, exclude := range excludes {
			if reflect.TypeOf(opt) == exclude {
				isExclude = true
				break
			}
		}
		if !isExclude {
			filteredOpts = append(filteredOpts, opt)
		}
		if _, ok := opt.(overrideProvideOption); ok {
			result.enableOverride = true
		} else if _, ok := opt.(resolveCyclicProvideOption); ok {
			result.enableResolveCyclic = true
		} else if lfo, ok := opt.(internal.LocationFixOption); ok {
			result.locationFixCallSkip = lfo.CallSkip
		}
	}
	return filteredOpts, result
}
