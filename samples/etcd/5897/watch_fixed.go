package v3rpc

import "sync"

type WatchID int64

type serverWatchStream struct {
	mu       sync.Mutex
	mu       sync.Mutex
	progress map[WatchID]bool
	prevKV   map[WatchID]bool
}

func newServerWatchStream() *serverWatchStream {
	return &serverWatchStream{
		progress: make(map[WatchID]bool),
		prevKV:   make(map[WatchID]bool),
	}
}

// recvLoop — BUG (pre-PR5897): writes progress/prevKV maps without sws.mu.
func (sws *serverWatchStream) recvLoop(id WatchID, p, pk bool) {
	sws.mu.Lock()
	defer sws.mu.Unlock()
	sws.progress[id] = p // BUG line 46
	sws.prevKV[id] = pk  // BUG line 49
}

// sendLoop — BUG: reads progress/prevKV maps without sws.mu.
func (sws *serverWatchStream) sendLoop(id WatchID) bool {
	sws.mu.Lock()
	defer sws.mu.Unlock()
	_ = sws.progress[id] // BUG line 56
	return sws.prevKV[id]
}
