package concurrent

import (
	"fmt"
	"sync"
)

// BackgroundWorker does work in the background.
type BackgroundWorker interface {
	Do(f func() (interface{}, error)) error
	Close() ([]interface{}, error)
}

// NewBackgroundWorker creates a new BackgroundWorker.
func NewBackgroundWorker() BackgroundWorker {
	return newBackgroundWorker()
}

type valueError struct {
	value interface{}
	err   error
}

type backgroundWorker struct {
	wg          *sync.WaitGroup
	valueErrors chan *valueError
	destroyable Destroyable
}

func newBackgroundWorker() *backgroundWorker {
	return &backgroundWorker{&sync.WaitGroup{}, make(chan *valueError), NewDestroyable(nil)}
}

func (b *backgroundWorker) Do(f func() (interface{}, error)) error {
	_, err := b.destroyable.Do(func() (interface{}, error) {
		b.wg.Add(1)
		go func() {
			value, err := f()
			b.valueErrors <- &valueError{value, err}
			b.wg.Done()
		}()
		return nil, nil
	})
	return err
}

func (b *backgroundWorker) Close() ([]interface{}, error) {
	var values []interface{}
	var errs []error
	if err := b.destroyable.Destroy(); err != nil {
		errs = append(errs, err)
	}
	go func() {
		for {
			valueError, ok := <-b.valueErrors
			if !ok {
				break
			}
			values = append(values, valueError.value)
			if valueError.err != nil {
				errs = append(errs, valueError.err)
			}
		}
	}()
	b.wg.Wait()
	close(b.valueErrors)
	if len(errs) == 0 {
		return values, nil
	}
	return values, fmt.Errorf("%v", errs)
}
