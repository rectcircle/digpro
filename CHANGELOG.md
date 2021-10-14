# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0][1.1.0] - 2021-10-14

### Add

* `digpro.Override()` API for override a registered provider
* `digpro.MakeExtractFunc(ptr)` API for make Invoke function to extract a value and assign to *ptr from dig.Container

### Fixed

* `dig.RootCause()` may be not work
* `Extract()` can't extract interface type value
* README spell

### Security

* improve test coverage

### Changed

## [1.0.0] - 2021-10-08

### Add

* Progressive use digpro
* Value Provider
* Property dependency injection
* Extract object from the container
* Add global container
* Export some function
  * `QuickPanic` function
  * `Visualize` function
  * `Unwrap` function

[1.0.0]: https://github.com/rectcircle/digpro/releases/tag/v1.0.0
[1.1.0]: https://github.com/rectcircle/digpro/compare/v1.0.0...v1.1.0
