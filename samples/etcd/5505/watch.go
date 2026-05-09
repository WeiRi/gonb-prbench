package v3rpc

import "sync"

type ctrlMsg struct{ id int }

type serverWatchStream struct {
	wg         sync.WaitGroup
	ctrlStream chan ctrlMsg
	closec     chan struct{}
	closeOnce  sync.Once
}

func newServerWatchStream() *serverWatchStream {
	return &serverWatchStream{
		ctrlStream: make(chan ctrlMsg, 8),
		closec:     make(chan struct{}),
	}
}

// recvLoop — BUG (pre-PR5505): sends on ctrlStream that close() may have closed.
func (sws *serverWatchStream) recvLoop() {
	defer sws.wg.Done()
	for {
		select {
		case sws.ctrlStream <- ctrlMsg{id: 1}: // line 45 send
			return
		case <-sws.closec:
			return
		}
	}
}

func (sws *serverWatchStream) drain() {
	defer sws.wg.Done()
	for {
		select {
		case <-sws.ctrlStream:
		case <-sws.closec:
			return
		}
	}
}

// close — BUG: close(ctrlStream) races with recvLoop send (line 63).
func (sws *serverWatchStream) close() {
	sws.closeOnce.Do(func() {
		close(sws.ctrlStream) // line 63 BUG
		close(sws.closec)
	})
	sws.wg.Wait()
}
