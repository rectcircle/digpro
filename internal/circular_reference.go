package internal

type PropertyInfo struct {
	ResolveCyclic bool
	Inputs        []ProvideInput
	Injected      bool
	Error         error
}
