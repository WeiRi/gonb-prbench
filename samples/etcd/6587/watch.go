package clientv3

type initReq struct {
	rev int64
}

type watcherStream struct {
	respc   chan int64
	initReq initReq
	donec   chan struct{}
}

func newWatcherStream() *watcherStream {
	return &watcherStream{
		respc: make(chan int64, 16),
		donec: make(chan struct{}),
	}
}

// serveSubstream — BUG (pre-PR6587): defer writes ws.initReq.rev (line 41).
func (w *watcherStream) serveSubstream() {
	defer func() {
		w.initReq.rev = -1 // BUG line 41 write inside defer
		close(w.donec)
	}()
	for r := range w.respc { // line 47
		w.initReq.rev = r
	}
}

// resume — BUG: reads w.initReq.rev concurrently (line 59).
func (w *watcherStream) resume() int64 {
	return w.initReq.rev // BUG line 59
}
