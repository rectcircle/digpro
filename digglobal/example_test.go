package digglobal_test

import (
	"fmt"
	"os"
	"reflect"

	"github.com/rectcircle/digpro"
	"github.com/rectcircle/digpro/digglobal"
	"go.uber.org/dig"
)

type Foo struct {
	A       string
	B       int
	C       int  `name:"c"`
	private bool //lint:ignore U1000 for test
	ignore  int  `digpro:"ignore"`
}

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

func init() {
	// register object
	digglobal.Supply("a")
	digglobal.Supply(1)
	digglobal.Supply(2, dig.Name("c"))
	digglobal.Supply(true)

	// use dig API
	digglobal.Provide(func(in struct {
		dig.In
		B int
		C int `name:"c"`
	}) int {
		return in.B + in.C
	}, dig.Name("d"))

	// register a struct
	digglobal.Struct(Foo{ignore: 3})

	// register two cyclic dependency strcuts
	digglobal.Struct(new(D1))
	digglobal.Struct(new(D2), digpro.ResolveCyclic())
}

func Example() {

	// if this is a test, can override and replace provider
	digglobal.Supply("aaa", digpro.Override())

	// use dig API
	digglobal.Invoke(func(a string) {
		fmt.Println("### a be override to \"aaa\" from \"a\" ###")
		fmt.Println(a)
	})

	// use digpro API
	foo, err := digglobal.Extract(Foo{})
	if err != nil {
		digpro.QuickPanic(err)
	}
	c, err := digglobal.Extract(int(0), digpro.ExtractByName("c"))
	if err != nil {
		digpro.QuickPanic(err)
	}
	d1, err := digglobal.Extract(new(D1))
	if err != nil {
		digpro.QuickPanic(err)
	}

	fmt.Println("### foo ###")
	fmt.Printf("%#v\n", foo)
	fmt.Println("### int[name=\"c\"] ###")
	fmt.Printf("%#v\n", c)
	fmt.Println("### d1 ###")
	fmt.Printf("%#s\n", d1)
	fmt.Println("### type of digglobal.Unwrap() ###")
	fmt.Println(reflect.TypeOf(digglobal.Unwrap()))
	fmt.Println("### inspect node and value <see stderr> ###")
	os.Stderr.WriteString(digglobal.String())
	fmt.Println("### inspect dot graph <see stderr> ###")
	digglobal.Visualize(os.Stderr)
	os.Stderr.WriteString("\n")
	// Output:
	// ### a be override to "aaa" from "a" ###
	// aaa
	// ### foo ###
	// digglobal_test.Foo{A:"aaa", B:1, C:2, private:true, ignore:3}
	// ### int[name="c"] ###
	// 2
	// ### d1 ###
	// D1: {D2: {D1: ..., Value: 'aaa'}, Value: 1}
	// ### type of digglobal.Unwrap() ###
	// *dig.Container
	// ### inspect node and value <see stderr> ###
	// ### inspect dot graph <see stderr> ###
}
