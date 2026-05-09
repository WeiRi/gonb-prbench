package pocgrpc2695

type serverHandlerTransport struct {
	writes   chan func()
	closedCh chan struct{}
}

func newServerHandlerTransport() *serverHandlerTransport {
	return &serverHandlerTransport{
		writes:   make(chan func(), 8),
		closedCh: make(chan struct{}),
	}
}

// do — BUG (pre-PR2695): outer-default + inner blocking send; close(writes) elsewhere races.
func (ht *serverHandlerTransport) do(fn func()) error {
	select {
	case <-ht.closedCh:
		return ErrClosing
	default:
		select {
		case ht.writes <- fn: // line 29 BUG send-on-writes
			return nil
		case <-ht.closedCh:
			return ErrClosing
		}
	}
}

func (ht *serverHandlerTransport) WriteStatus() {
	close(ht.writes)
	ht.Close()
}

func (ht *serverHandlerTransport) Close() {
	select {
	case <-ht.closedCh:
	default:
		close(ht.closedCh)
	}
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

type errClosing string

func (e errClosing) Error() string { return string(e) }

var ErrClosing = errClosing("connection is closing")
