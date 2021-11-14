package digpro

import "go.uber.org/dig"

type overrideProvideOption struct {
	dig.ProvideOption
}

type resolveCyclicProvideOption struct {
	dig.ProvideOption
}

func isOriginProvideOption(opt dig.ProvideOption) bool {
	if _, ok := opt.(resolveCyclicProvideOption); ok {
		return false
	} else if _, ok := opt.(overrideProvideOption); ok {
		return false
	}
	return true
}
