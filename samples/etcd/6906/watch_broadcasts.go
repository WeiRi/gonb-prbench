package grpcproxy

import "sync"

type watcher struct{ id int }

type watchBroadcast struct {
	mu sync.Mutex
}

type watchBroadcasts struct {
	mu       sync.Mutex
	bcasts   map[*watchBroadcast]struct{}
	watchers map[*watcher]*watchBroadcast
	updatec  chan *watchBroadcast
	donec    chan struct{}
}

// empty — BUG (pre-PR6906): reads len(wbs.bcasts) without lock.
func (wbs *watchBroadcasts) empty() bool {
	return len(wbs.bcasts) == 0 // BUG: unlocked map len
}
