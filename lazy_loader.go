package concurrent

import (
	"sync"
	"sync/atomic"
)

type LazyLoader interface {
	Load() (interface{}, error)
}

func NewLazyLoader(f func() (interface{}, error)) LazyLoader {
	return newLazyLoader(f)
}

type lazyLoader struct {
	once  *sync.Once
	f     func() (interface{}, error)
	value *atomic.Value
	err   *atomic.Value
}

func newLazyLoader(f func() (interface{}, error)) *lazyLoader {
	return &lazyLoader{&sync.Once{}, f, &atomic.Value{}, &atomic.Value{}}
}

func (this *lazyLoader) Load() (interface{}, error) {
	this.once.Do(func() {
		value, err := this.f()
		if value != nil {
			this.value.Store(value)
		}
		if err != nil {
			this.err.Store(err)
		}
	})
	value := this.value.Load()
	err := this.err.Load()
	if err != nil {
		return nil, err.(error)
	}
	return value, nil
}
