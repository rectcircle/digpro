# :hammer_and_pick: digpro

[![MIT][license-img]][license] [![GoDoc][doc-img]][doc] [![GitHub release][release-img]][release] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][report-card-img]][report-card]

[English]((README.md)) | 简体中文

## 介绍

digpro 是一个增强版的 [uber-go/dig][dig-github]，继承了 [uber-go/dig][dig-github] 的全部能力，并添加了如下特性：

* 渐进式的使用 digpro
* 值 Provider
* 属性依赖注入
* 从容器里提取对象
* Override 已注册的 Provider
* 添加全局容器
* 导出一些函数
  * `QuickPanic` 函数
  * `Visualize` 函数
  * `Unwrap` 函数
* 循环引用

## 安装

```bash
go get github.com/rectcircle/digpro
```

## 指南

### [dig][dig-github] 简介

[uber-go/dig][dig-github] 是一个 uber 开源的 Go 语言，轻量依赖注入库。该库通过 Go 反射能力，在运行时提供：

* 基于构造函数的依赖注入能力
* 对象和接口绑定
* 命名对象和组对象
* 参数对象、结果对象和可选依赖

更多细节参见：[go docs][dig-go-docs]

### 为什么需要 digpro 项目

[dig][dig-github] 提供了非常轻量级的运行依赖注入，代码优质且稳定。但是其缺少如下比较有用的能力：

* 属性依赖注入，只需提供一个结构体类型，依赖注入库就可以构造一个结构体对象并将依赖注入到该对象中，并放入容器中。该能力可以给依赖注入使用者节省大量的时间，避免编写大量的样板式的构造函数。
* 值 Provider，可以将用户提供的构造好的对象直接放到容器中。
* 提取对象，将容器内构造出的对象提取出来，以便使用。

### 渐进使用

#### 低级 API

可以直接使用 `c := dig.New()` 构造的容器。通过如下方式使用

```go
c := dig.New()
c.Provide(digpro.Supply(...))
c.Provide(digpro.Struct(...))
digpro.Extract(c, ...)
```

该方式可以零成本的在已经使用 [dig][dig-github] 的项目中引入 digpro 所提供的所有能力。但是由于原理限制，该方式报错没有正确的代码文件和行号的信息，只能显示如下信息

```
... missing dependencies for function "reflect".makeFuncStub (/usr/local/Cellar/go/1.17.1/libexec/src/reflect/asm_amd64.s:30) ...
```

因此，如果从零开始使用 digpro 建议使用高级 API 或全局容器。

#### 高级 API

通过 `c := digpro.New()` 构造一个 `dig.Container` 的 wrapper 对象（`digpro.ContainerWrapper`）。通过如下方式使用

```go
c := digpro.New()
c.Supply(...)
c.Struct(...)
c.Extract(...)
```

可以看出来该方式比低级 API 更加优雅简洁。

### 全局容器

通过导入 `digglobal` 包，可以使用全局 `digpro.ContainerWrapper` 对象，如果追求开发极致的开发效率，可以使用该用法

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

注意：对于全局容器，Provider 类型的函数（`Provide`、`Struct`、`Supply`）将不再返回错误，直接 `Panic`

### 值类型依赖注入

可以将用户提供的构造好的对象直接放到容器中

```go
func Supply(value interface{}) interface{}
func (c *ContainerWrapper) Supply(value interface{}, opts ...dig.ProvideOption) error
```

* value 参数支持任意非 `error` 类型
* 如果 value 为无类型的 nil，则在容器中 `type = interface{}` 且 `value = nil` 的对象

Example

```go
// High Level API
c.Supply("a")
// Lower Level API
c.Provide(digpro.Supply("a"))
```

等价于

```go
c.Provide(func() string {return "a"})
```

### 属性依赖注入

提供一个结构体类型，依赖注入库就可以构造一个结构体对象并将依赖注入到该对象中，并放入容器中

```go
func Struct(structOrStructPtr interface{}) interface{}
func (c *ContainerWrapper) Struct(structOrStructPtr interface{}, opts ...dig.ProvideOption) error
```

* structOrStructPtr 必须为 struct 类型，或者 struct 指针类型

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

等价于

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

### 提取对象

将容器内构造出的对象提取出来，以便使用。

```go
func MakeExtractFunc(ptr interface{}, opts ...ExtractOption) interface{}
func Extract(c *dig.Container, typ interface{}, opts ...ExtractOption) (interface{}, error)
func (c *ContainerWrapper) Extract(typ interface{}, opts ...ExtractOption) (interface{}, error)
```

对于两个 `Extract` 函数，如果想提取一个非接口类型，typ 和返回值的类型关系为 `reflect.TypeOf(result) == reflect.TypeOf(typ)`，即

```go
func(int) -> int    // func(int(0)) -> int
func(*int) -> *int  // func(new(int)) -> *int
```

对于两个 `Extract` 函数，如果想提取一个接口类型，typ 和返回值的类型关系为 `reflect.TypeOf(result) == reflect.TypeOf(typ).Elem()`，即

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

等价于

```go
var i interface{}
err := c.Invoke(func(_i int){
  i = _i
})
```

### Override

> :warning: 仅支持高级 API

在使用依赖注入时，会按照生产环境的配置来注册 Provider，并在 main 函数中启动服务。当要在其他环境启动服务时，（比如测试），一般只需要，将少量的某几个生产环境的 Provider 替换为专有的 Provider（比如 db mock）。

为了更优雅支持如上场景，添加 Override 能力。

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

为了提前暴露问题，如果容器里不存在相同 Provider，使用  `digpro.Override()` 将返回错误 `no provider to override was found`

### 循环引用

> :warning: 仅支持高级 API `Struct` 方法

在某些情况下，可能会出现几个结构体之间存在循环引用的情况。在 Go 语言中表现为如下两种情况以及这两种情况的混合情况：

* 指针类型：两个结构体相互包含指向对方的指针
* 接口类型：两个结构体分别实现了两个接口，这两个结构体包含了指向对方实现接口的引用

为了解决两种循环引用的场景，digpro 提供了给 `*digpro.WrapContainer.Struct` 方法，添加了一个选项 `digpro.ResolveCyclic()`。当出现循环引用时，
只需在引用环的任意一个结构体的 `Struct` 方法调用处，添加 `digpro.ResolveCyclic()` 选项，即可自动解决循环引用的问题。

Example

```go
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
```

### 其他

#### QuickPanic

```go
func QuickPanic(errs ...error)
```

如果 errs 有任何一个不是 nil，将直接 panic

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

将 dot graph 写入 io.Writer

#### Unwrap

```go
func (c *ContainerWrapper) Unwrap() *dig.Container
```

从 `*digpro.ContainerWrapper` 中获取 `*dig.Container`

## 最佳实践

### 配置文件及配置项

> [Playground](https://play.studygolang.com/p/ZsShZgajrgH)

假设一个项目需要从配置文件配置读取到一个结构体中，很多组件依赖该配置结构体的某些属性，使用 digpro 的做法如下

```go
package main

import (
	"fmt"
	"github.com/rectcircle/digpro"
	"go.uber.org/dig"
)

type Config struct {
	dig.Out // it important
	DB      struct {
	        dig.Out        // it important
	        DSN     string `name:"config.db.dsn"` // must decalre name
	}
}

type DB struct {
	DB_DSN string `name:"config.db.dsn"` // this name must be same as config.db.dsn of Config struct
}


func main() {
	c := digpro.New()
	digpro.QuickPanic(
		c.Provide(func() Config {
			// ... read config from file
			return Config{
				DB: struct {
					dig.Out
					DSN string `name:"config.db.dsn"`
				}{
					DSN: "this is db dsn",
				},
			}
		}),
		c.Struct(new(DB)),
		c.Invoke(func(db *DB) {
			fmt.Println(db.DB_DSN)
		}),
	)
	// Output: this is db dsn
}
```

[dig-github]: https://github.com/uber-go/dig
[dig-go-docs]: https://pkg.go.dev/go.uber.org/dig#example-package-Minimal

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
