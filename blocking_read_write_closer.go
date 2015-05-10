package concurrent

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// NewBlockingReadWriteCloser creates an io.ReadWriteCloser that blocks
// on Read until there is data or until closed.
func NewBlockingReadWriteCloser() io.ReadWriteCloser {
	return newBlockingReadWriteCloser()
}

type blockingReadWriteCloser struct {
	buffer *bytes.Buffer
	cv     *sync.Cond
	closed VolatileBool
}

func newBlockingReadWriteCloser() *blockingReadWriteCloser {
	return &blockingReadWriteCloser{
		bytes.NewBuffer(nil),
		sync.NewCond(&sync.Mutex{}),
		NewVolatileBool(false),
	}
}

func (b *blockingReadWriteCloser) Write(p []byte) (int, error) {
	if b.closed.Value() {
		return 0, errors.New("concurrent: already closed")
	}
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	b.cv.L.Lock()
	defer b.cv.L.Unlock()
	n, err := b.buffer.Write(p)
	b.cv.Signal()
	return n, err
}

func (b *blockingReadWriteCloser) Read(p []byte) (int, error) {
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	b.cv.L.Lock()
	for b.buffer.Len() == 0 {
		if b.closed.Value() {
			return 0, io.EOF
		}
		b.cv.Wait()
	}
	defer b.cv.L.Unlock()
	return b.buffer.Read(p)
}

func (b *blockingReadWriteCloser) Close() error {
	b.cv.L.Lock()
	defer b.cv.L.Unlock()
	if !b.closed.CompareAndSwap(false, true) {
		return errors.New("concurrent: already called close")
	}
	b.cv.Broadcast()
	return nil
}
