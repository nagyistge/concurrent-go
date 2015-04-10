package concurrent

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

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

func (this *blockingReadWriteCloser) Write(p []byte) (int, error) {
	if this.closed.Value() {
		return 0, errors.New("concurrent: already closed")
	}
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	this.cv.L.Lock()
	defer this.cv.L.Unlock()
	n, err := this.buffer.Write(p)
	this.cv.Signal()
	return n, err
}

func (this *blockingReadWriteCloser) Read(p []byte) (int, error) {
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	this.cv.L.Lock()
	for this.buffer.Len() == 0 {
		if this.closed.Value() {
			return 0, io.EOF
		}
		this.cv.Wait()
	}
	defer this.cv.L.Unlock()
	return this.buffer.Read(p)
}

func (this *blockingReadWriteCloser) Close() error {
	this.cv.L.Lock()
	defer this.cv.L.Unlock()
	if !this.closed.CompareAndSwap(false, true) {
		return errors.New("concurrent: already called close")
	}
	return nil
}
