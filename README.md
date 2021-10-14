# :hammer_and_pick: digpro

[![MIT][license-img]][license] [![GoDoc][doc-img]][doc] [![GitHub release][release-img]][release] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][report-card-img]][report-card]

English | [简体中文](README_zh-CN.md)

## Introduction

digpro is a enhanced [uber-go/dig][dig-github], inherit all feature of [uber-go/dig][dig-github] and add the following features:

* Progressive use digpro
* Value Provider
* Property dependency injection
* Extract object from the container
* Override a registered Provider
* Add global container
* Export some function
  * `QuickPanic` function
  * `Visualize` function
  * `Unwrap` function

## Installation

```bash
go get github.com/rectcircle/digpro
```

## Document

https://pkg.go.dev/github.com/rectcircle/digpro

## Guide

### [dig][dig-github] Introduction

[uber-go/dig][dig-github] is a lightweight dependency injection library that supported by uber for go language. The library driven by reflection and has the following features:

* Dependency injection based on constructor
* Object and interface binding
* Named objects and group objects
* Parameter object, result object and optional dependencies

More see: [go docs][dig-go-docs]

### Why need digpro

[dig][dig-github] provides a very lightweight runtime dependency injection, and the code is high-quality and stable. But it lacks the following more useful features:

* Property dependency injection, by simply providing a structure type, the dependency injection library can construct a structure object and inject the dependency into that object and into the container. This feature can save dependency injection users a lot of time and avoid writing a lot of sample constructors.
* Value Provider, which can take a constructed object provided by the user and put it directly into a container.
* Extract Object, which extracts the constructed object from the container for use.

### Progressive use

#### Lower level API

Containers constructed with `c := dig.New()` can be used directly. use it like

```go
c := dig.New()
c.Provide(digpro.Supply(...))
c.Provide(digpro.Struct(...))
digpro.Extract(c, ...)
```

This approach introduces all the capabilities provided by digpro at zero cost to projects already using [dig][dig-github]. However, due to a limitation of the principle, this method reports errors without the correct code file and line number information, and only displays the following information

```
... missing dependencies for function "reflect".makeFuncStub (/usr/local/Cellar/go/1.17.1/libexec/src/reflect/asm_amd64.s:30) ...
```

Therefore, if starting from scratch with digpro it is recommended to use the High level API or the global container.

#### High level API

Construct a wrapper object for `dig.Container` (`digpro.ContainerWrapper`) with `c := digpro.New()`. use it like

```go
c := digpro.New()
c.Supply(...)
c.Struct(...)
c.Extract(...)
```

This approach is much more elegant and simple than the low-level API.

### Global container

By importing the `digglobal` package, you can use the global `digpro.ContainerWrapper` object, which can be used if you are looking for the ultimate efficiency

```go
import "github.com/rectcircle/digpro/digglobal"

// dao/module.go
func init() {
  digglobal.Struct(new(XxxDAO))
  digglobal.Struct(new(XxxDAO))
  digglobal.Struct(new(XxxDAO))
}

// service/module.go
func init() {
  digglobal.Struct(new(XxxService))
  digglobal.Struct(new(XxxService))
  digglobal.Struct(new(XxxService))
}

// controller/module.go
func init() {
  digglobal.Struct(new(XxxService))
  digglobal.Struct(new(XxxService))
  digglobal.Struct(new(XxxService))
}

// main.go
func main() {
  digglobal.Supply(Config{})
  digglobal.Supply(DB{})
  digglobal.Provide(NewServer)
  server, err := digglobal.Extract(new(Server))
  digpro.QuickPanic(err)
  server.Run()
}
```

Note: For global containers, functions of type Provider (`Provide`, `Struct`, `Supply`) will no longer return an error, directly `Panic`

### Value Provider

It can take a constructed object provided by the user and put it directly into a container.

```go
func Supply(value interface{}) interface{}
func (c *ContainerWrapper) Supply(value interface{}, opts ...dig.ProvideOption) error
```

* The value parameter supports any non-`error` type
* If value is untyped nil, then the object in the container with `type = interface{}` and `value = nil`

Example

```go
// High Level API
c.Supply("a")
// Lower Level API
c.Provide(digpro.Supply("a"))
```

Equals to

```go
c.Provide(func() string {return "a"})
```

### Property dependency injection

By simply providing a structure type, the dependency injection library can construct a structure object and inject the dependency into that object and into the container.

```go
func Struct(structOrStructPtr interface{}) interface{}
func (c *ContainerWrapper) Struct(structOrStructPtr interface{}, opts ...dig.ProvideOption) error
```

* structOrStructPtr must be of struct type, or struct pointer type

Example

```go
type Foo struct {
	A       string
	B       int
	C       int `name:"c"`
	private bool
	ignore  int `digpro:"ignore"`
}
// High Level API
c.Struct(Foo{
  ignore: 3,
})
// Lower Level API
c.Provide(digpro.Struct(Foo{
  ignore: 3,
}))
```

Equals to

```go
c.Provide(func(in struct {
  A       string
  B       int
  C       int `name:"c"`
  Private bool
}) Foo {
  return Foo{
    A:       in.A,
    B:       in.B,
    C:       in.C,
    private: in.Private,
    ignore:  3,
  }
}),
```

### Extract object

Extracts the object constructed inside the container for use.

```go
func MakeExtractFunc(ptr interface{}, opts ...ExtractOption) interface{}
func Extract(c *dig.Container, typ interface{}, opts ...ExtractOption) (interface{}, error)
func (c *ContainerWrapper) Extract(typ interface{}, opts ...ExtractOption) (interface{}, error)
```

For two `Extract` functions, if want to extract a non-interface, `reflect.TypeOf(result) == reflect.TypeOf(typ)`

```go
func(int) -> int    // func(int(0)) -> int
func(*int) -> *int  // func(new(int)) -> *int
```

For two `Extract` functions, if want to extract a interface, `reflect.TypeOf(result) == reflect.TypeOf(typ).Elem()`

```go
type A interface { ... }
func(A) -> error   // func(A(nil)) -> error
func(*A) -> A      // func(new(A)) -> A
func(**A) -> *A    // func(new(*A)) -> *A
```

Example

```go
// High Level API
i, err := c.Extract(int(0)) 
// Lower Level API (1)
i, err := digpro.Extract(c, int(0))
// Lower Level API (2)
var i int
err := digpro.Invoke(digpro.MakeExtractFunc(&i))
```

Equals to

```go
var i interface{}
err := c.Invoke(func(_i int){
  i = _i
})
```

## Override

> :warning: Only support High Level API

When using dependency injection, the Provider is registered according to the configuration of the production environment, and the service is started in the main function. When starting a service in another environment (e.g. testing), it is generally only necessary to replace a small number of the production environment's Providers with proprietary ones (e.g. db mock).

To more elegantly support scenarios such as the above, add the Override capability.

Example

```go
c := digpro.New()
_ = c.Supply(1) // please handle error in production
_ = c.Supply(1, digpro.Override())
// _ = c.Supply("a", digpro.Override())  // has error
i, _ := c.Extract(0)
fmt.Println(i.(int) == 1)
// Output: true
```

To expose the problem in advance, using `digpro.Override()` will return the error `no provider to override was found` if the same Provider does not exist in the container

### Others

#### QuickPanic

```go
func QuickPanic(errs ...error)
```

If any of the errs is not nil, it will panic directly

Example

```go
c := digpro.New()
digpro.QuickPanic(
	c.Supply(1),
	c.Supply(1),
)
// panic: [1]: cannot provide function "xxx".Xxx (xxx.go:n): cannot provide int from [0]: already provided by "xxx".Xxx (xxx.go:m)
```

#### Visualize

```go
func (c *ContainerWrapper) Visualize(w io.Writer, opts ...dig.VisualizeOption) error
```

Write dot graph to io.Writer.

#### Unwrap

```go
func (c *ContainerWrapper) Unwrap() *dig.Container
```

From `*digpro.ContainerWrapper` obtain `*dig.Container`

## TODO

* [ ] Circular reference problem (change high level api `Struct` 、 `Invoke` and `Extract` Implementation as  `() => new(struct)`, and assembled at the time of extraction)(disregarding the non-pointer case)

[dig-github]: https://github.com/uber-go/dig
[dig-go-docs]: https://pkg.go.dev/go.uber.org/dig

[license-img]: https://img.shields.io/github/license/rectcircle/digpro
[license]: https://github.com/rectcircle/digpro/blob/master/LICENSE

[doc-img]: http://img.shields.io/badge/GoDoc-Reference-blue.svg
[doc]: https://godoc.org/github.com/rectcircle/digpro

[release-img]: https://img.shields.io/github/release/rectcircle/digpro.svg
[release]: https://github.com/rectcircle/digpro/releases

[ci-img]: https://github.com/rectcircle/digpro/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/rectcircle/digpro/actions/workflows/go.yml

[cov-img]: https://codecov.io/gh/rectcircle/digpro/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/rectcircle/digpro/branch/master

[report-card-img]: https://goreportcard.com/badge/github.com/rectcircle/digpro
[report-card]: https://goreportcard.com/report/github.com/rectcircle/digpro
