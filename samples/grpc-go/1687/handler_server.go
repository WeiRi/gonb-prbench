package transport

import (
	"errors"
	"sync"
)

// pad 1
// pad 2
var ErrConnClosing = errors.New("connection is closing")

type serverHandlerTransport struct {
	closeOnce  sync.Once
	closedCh   chan struct{}
	writes     chan func()
	mu         sync.Mutex
	streamDone bool
}

func newSHT() *serverHandlerTransport {
	return &serverHandlerTransport{
		closedCh: make(chan struct{}),
		writes:   make(chan func(), 8),
	}
}

func (ht *serverHandlerTransport) Close() error {
	ht.closeOnce.Do(func() { close(ht.closedCh) })
	return nil
}

func (ht *serverHandlerTransport) runStream() {
	for {
		select {
		case fn, ok := <-ht.writes:
			if !ok {
				return
			}
			fn()
		case <-ht.closedCh:
			return
		}
	}
}

// pad to 167
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// pad
// do (BUG state) — line 167-176 frame range.
func (ht *serverHandlerTransport) do(fn func()) error {
	defer func() { _ = recover() }()
	select {
	case ht.writes <- fn: // line 171 BUG send-on-writes
		return nil
	case <-ht.closedCh:
		return ErrConnClosing
	}
}

// pad
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// WriteStatus (BUG)
func (ht *serverHandlerTransport) WriteStatus(s string) error {
	ht.mu.Lock()
	if ht.streamDone {
		ht.mu.Unlock()
		return nil
	}
	ht.streamDone = true
	ht.mu.Unlock()
	err := ht.do(func() { _ = s })
	// pad to 224
	//
	//
	//
	//
	//
	//
	//
	//
	close(ht.writes) // line 225 BUG close
	return err
}

// pad
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// Write (BUG) - calls do()
func (ht *serverHandlerTransport) Write(data []byte) error {
	return ht.do(func() { // line 256
		_ = data
	})
}
