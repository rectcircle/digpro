# :hammer_and_pick: digpro

[![MIT][license-img]][license] [![GoDoc][doc-img]][doc] [![GitHub release][release-img]][release] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][report-card-img]][report-card]

[English]((README.md)) | 简体中文

## 介绍

digpro 是一个增强版的 [uber-go/dig][dig-github]，继承了 [uber-go/dig][dig-github] 的全部能力，并添加了如下特性：

* 渐进式的使用 digpro
* 值 Provider
* 属性依赖注入
* 从容器里提取对象
* 添加全局容器
* 导出一些函数
  * `QuickPanic` 函数
  * `Visualize` 函数
  * `Unwrap` 函数

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
func Extract(c *dig.Container, typInterface interface{}, opts ...ExtractOption) (interface{}, error)
func (c *ContainerWrapper) Extract(typInterface interface{}, opts ...ExtractOption) (interface{}, error)
```

* structOrStructPtr 必须为 struct 类型，或者 struct 指针类型

Example

```go
// High Level API
i, err := c.Extract(int(0))
// Lower Level API
i, err := digpro.Extract(c, int(0))
```

等价于

```go
var i interface{}
err := c.Invoke(func(_i int){
  i = _i
})
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

## TODO

* [ ] `digpro.ForceOverride` 覆盖 Provider （支持测试）
* [ ] 循环引用问题（修改高级API中 `Struct` 、 `Invoke` 和 `Extract` 实现为 `() => new(struct)`，并在提取的时候进行拼装（不考虑非指针的情况）

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
