package concurrent

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	// ErrAlreadyDestroyed is the error returned when a function is called
	// on a Destroyable that has already been destroyed.
	ErrAlreadyDestroyed = errors.New("concurrent: already destroyed")
)

// Destroyable is a wrapper for any object that allows an atomic destroy operation
// to be performed, and will monitor if other functions are called after the object
// has already been destroyed. Destroyables can also have children Destroyables
// that will be automatically recursively destroyed.
type Destroyable interface {
	Destroy() error
	Do(func() (interface{}, error)) (interface{}, error)
	AddChild(Destroyable) error
}

// NewDestroyable creates a new Destroyable.
func NewDestroyable(destroyCallback func() error) Destroyable {
	return newDestroyable(destroyCallback)
}

type destroyable struct {
	destroyCallback func() error
	cv              *sync.Cond
	destroyed       VolatileBool
	numOperations   int32
	children        []Destroyable
}

func newDestroyable(destroyCallback func() error) *destroyable {
	return &destroyable{
		destroyCallback,
		sync.NewCond(&sync.Mutex{}),
		NewVolatileBool(false),
		0,
		make([]Destroyable, 0),
	}
}

func (d *destroyable) Destroy() error {
	d.cv.L.Lock()
	if !d.destroyed.CompareAndSwap(false, true) {
		d.cv.L.Unlock()
		return ErrAlreadyDestroyed
	}
	for atomic.LoadInt32(&d.numOperations) > 0 {
		d.cv.Wait()
	}
	defer d.cv.L.Unlock()
	for _, child := range d.children {
		// children can destroy themselves, ignore error
		_ = child.Destroy()
	}
	if d.destroyCallback != nil {
		return d.destroyCallback()
	}
	return nil
}

func (d *destroyable) Do(f func() (interface{}, error)) (interface{}, error) {
	d.cv.L.Lock()
	if d.destroyed.Value() {
		d.cv.L.Unlock()
		return nil, ErrAlreadyDestroyed
	}
	atomic.AddInt32(&d.numOperations, 1)
	d.cv.L.Unlock()
	value, err := f()
	atomic.AddInt32(&d.numOperations, -1)
	d.cv.Signal()
	return value, err
}

func (d *destroyable) AddChild(destroyable Destroyable) error {
	d.cv.L.Lock()
	if d.destroyed.Value() {
		d.cv.L.Unlock()
		return ErrAlreadyDestroyed
	}
	defer d.cv.L.Unlock()
	d.children = append(d.children, destroyable)
	return nil
}
