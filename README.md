[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/go-concurrent)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/go-concurrent/blob/master/LICENSE)

Some concurrent utilities for Go.

## Installation
```bash
go get -u github.com/peter-edge/go-concurrent
```

## Import
```go
import (
    "github.com/peter-edge/go-concurrent"
)
```

## Usage

```go
var (
	ErrAlreadyDestroyed = errors.New("concurrent: already destroyed")
)
```

#### func  NewBlockingReadWriteCloser

```go
func NewBlockingReadWriteCloser() io.ReadWriteCloser
```

#### type BackgroundWorker

```go
type BackgroundWorker interface {
	Do(f func() (interface{}, error))
	WaitToFinish() ([]interface{}, error)
}
```


#### func  NewBackgroundWorker

```go
func NewBackgroundWorker() BackgroundWorker
```

#### type CombinedError

```go
type CombinedError interface {
	error
	Errors() []error
}
```


#### func  NewCombinedError

```go
func NewCombinedError(errs []error) CombinedError
```

#### type Destroyable

```go
type Destroyable interface {
	Destroy() error
	Do(func() (interface{}, error)) (interface{}, error)
	AddChild(Destroyable) error
}
```


#### func  NewDestroyable

```go
func NewDestroyable(destroyCallback func() error) Destroyable
```

#### type LazyLoader

```go
type LazyLoader interface {
	Load() (interface{}, error)
}
```


#### func  NewLazyLoader

```go
func NewLazyLoader(f func() (interface{}, error)) LazyLoader
```

#### type VolatileBool

```go
type VolatileBool interface {
	Value() bool
	// return old value == new value
	CompareAndSwap(oldBool bool, newBool bool) bool
}
```

TODO(pedge): is this even needed? need to understand go memory model better

#### func  NewVolatileBool

```go
func NewVolatileBool(initialBool bool) VolatileBool
```
