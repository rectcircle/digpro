package digpro_test

import (
	"fmt"

	"github.com/rectcircle/digpro"
	"go.uber.org/dig"
)

type D1 struct {
	D2    *D2
	Value int
}

func (d1 *D1) String() string {
	return fmt.Sprintf("D1: {D2: {D1: ..., Value: '%s'}, Value: %d}", d1.D2.Value, d1.Value)
}

type D2 struct {
	D1    *D1
	Value string
}

func (d2 *D2) String() string {
	return fmt.Sprintf("D2: {D1: {D2: ..., Value: %d}, Value: '%s'}", d2.D1.Value, d2.Value)
}

func resolvePointerTypeCyclicDependency() {
	c := digpro.New()
	_ = c.Supply(1) // please handle error in production
	_ = c.Supply("a")
	_ = c.Struct(new(D1), digpro.ResolveCyclic()) // enable resolve cyclic dependency
	_ = c.Struct(new(D2))
	d1, _ := c.Extract(new(D1))
	d2, _ := c.Extract(new(D2))
	fmt.Println(d1.(*D1).String())
	fmt.Println(d2.(*D2).String())
}

type I1 interface{ String1() string }
type I2 interface{ String2() string }

type DI1 struct {
	I2    I2
	Value int
}

func (d1 *DI1) String1() string {
	return fmt.Sprintf("DI1: {I2: {I1: ..., Value: '%s'}, Value: %d}", d1.I2.(*DI2).Value, d1.Value)
}

type DI2 struct {
	I1    I1
	Value string
}

func (d2 *DI2) String2() string {
	return fmt.Sprintf("DI2: {I1: {I2: ..., Value: %d}, Value: '%s'}", d2.I1.(*DI1).Value, d2.Value)
}

func resolveInterfaceTypeCyclicDependency() {
	c := digpro.New()
	_ = c.Supply(1) // please handle error in production
	_ = c.Supply("a")
	_ = c.Struct(new(DI1), dig.As(new(I1)))
	_ = c.Struct(new(DI2), dig.As(new(I2)), digpro.ResolveCyclic()) // enable resolve cyclic dependency
	i1, _ := c.Extract(new(I1))
	i2, _ := c.Extract(new(I2))
	fmt.Println(i1.(I1).String1())
	fmt.Println(i2.(I2).String2())
}

func ExampleResolveCyclic() {
	resolvePointerTypeCyclicDependency()
	resolveInterfaceTypeCyclicDependency()
	// Output:
	// D1: {D2: {D1: ..., Value: 'a'}, Value: 1}
	// D2: {D1: {D2: ..., Value: 1}, Value: 'a'}
	// DI1: {I2: {I1: ..., Value: 'a'}, Value: 1}
	// DI2: {I1: {I2: ..., Value: 1}, Value: 'a'}
}
