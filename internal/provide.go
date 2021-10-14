package internal

import (
	"go.uber.org/dig"
)

func ProvideWithLocationForPC(Provide func(constructor interface{}, opts ...dig.ProvideOption) error, callSkip int, constructor interface{}, opts ...dig.ProvideOption) error {
	return WrapErrorWithLocationForPC(callSkip, func(pc uintptr) error {
		if pc != 0 {
			return Provide(constructor, append([]dig.ProvideOption{dig.LocationForPC(pc)}, opts...)...)
		} else {
			return Provide(constructor, opts...)
		}
	})
}
