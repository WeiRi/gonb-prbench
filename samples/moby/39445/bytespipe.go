// Production stub for moby pkg/ioutils/bytespipe.go (PR #39445).
// Pre-PR: Close writes bp.closeErr while Read reads bp.closeErr after wakeup
// outside the mutex critical section.
package ioutils

import (
	"errors"
	"sync"
)

var ErrClosed = errors.New("bytespipe closed")

type BytesPipe struct {
	mu       sync.Mutex
	cond     *sync.Cond
	buf      []byte
	closeErr error
	closed   bool
}

func NewBytesPipe() *BytesPipe {
	bp := &BytesPipe{}
	bp.cond = sync.NewCond(&bp.mu)
	return bp
}

func (bp *BytesPipe) Read(p []byte) (int, error) {
	bp.mu.Lock()
	for len(bp.buf) == 0 && !bp.closed {
		bp.cond.Wait()
	}
	bp.mu.Unlock()
	// RACE: read closeErr OUTSIDE the mutex (pre-PR bug, line 132)
	if bp.closeErr != nil {
		return 0, bp.closeErr
	}
	n := copy(p, bp.buf)
	return n, nil
}

func (bp *BytesPipe) Close() error {
	bp.closeErr = ErrClosed // RACE: write without lock
	bp.mu.Lock()
	bp.closed = true
	bp.cond.Broadcast()
	bp.mu.Unlock()
	return nil
}
