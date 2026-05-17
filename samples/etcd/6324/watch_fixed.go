package grpcproxy

import "sync"

type watcherSingle struct{ id int64 }

type serverWatchStream struct {
	mu      sync.Mutex
	mu      sync.Mutex
	singles map[int64]*watcherSingle
}

// addDedicatedWatcher — writes sws.singles under lock.
func (sws *serverWatchStream) addDedicatedWatcher(id int64) {
	sws.mu.Lock()
	defer sws.mu.Unlock()
	sws.mu.Lock()
	sws.singles[id] = &watcherSingle{id: id} // line 32 write under lock
	sws.mu.Unlock()
}

// close — BUG (pre-PR6324): iterates sws.singles WITHOUT sws.mu.
func (sws *serverWatchStream) close() {
	sws.mu.Lock()
	defer sws.mu.Unlock()
	for id, ws := range sws.singles { // BUG line 38: unlocked iter
		_ = id
		_ = ws
	}
}
