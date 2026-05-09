package transport

import (
	"errors"
	"sync"
)

var ErrConnClosing = errors.New("connection is closing")

type serverHandlerTransport struct {
	closeOnce sync.Once
	closedCh  chan struct{}
	writes    chan func()
}

// do — BUG (pre-PR1115): selects send-on-writes / closedCh; if writes is closed
// concurrently with the send, panic + race.
func (ht *serverHandlerTransport) do(fn func()) error {
	select {
	case ht.writes <- fn: // line 31 BUG send
		return nil
	case <-ht.closedCh:
		return ErrConnClosing
	}
}
