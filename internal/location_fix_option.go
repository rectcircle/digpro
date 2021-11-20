package internal

import "go.uber.org/dig"

type LocationFixOption struct {
	dig.ProvideOption
	CallSkip int
}
