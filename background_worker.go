package concurrent

import "sync"

type BackgroundWorker interface {
	Do(f func() (interface{}, error))
	WaitToFinish() ([]interface{}, error)
}

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

func (this *backgroundWorker) Do(f func() (interface{}, error)) {
	this.destroyable.Do(func() (interface{}, error) {
		this.wg.Add(1)
		go func() {
			value, err := f()
			this.valueErrors <- &valueError{value, err}
			this.wg.Done()
		}()
		return nil, nil
	})
}

func (this *backgroundWorker) WaitToFinish() ([]interface{}, error) {
	values := make([]interface{}, 0)
	errs := make([]error, 0)
	if err := this.destroyable.Destroy(); err != nil {
		errs = append(errs, err)
	}
	go func() {
		for {
			valueError, ok := <-this.valueErrors
			if !ok {
				break
			}
			values = append(values, valueError.value)
			if valueError.err != nil {
				errs = append(errs, valueError.err)
			}
		}
	}()
	this.wg.Wait()
	close(this.valueErrors)
	if len(errs) == 0 {
		return values, nil
	}
	return values, NewCombinedError(errs)
}
