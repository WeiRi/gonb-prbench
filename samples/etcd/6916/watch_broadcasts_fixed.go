package grpcproxy

import (
	"sync"

	"context"
)

type watcher struct{ id int }

type watchBroadcast struct {
	mu        sync.Mutex
	cancel    context.CancelFunc
	donec     chan struct{}
	receivers map[*watcher]struct{}
	nextrev   int64
	responses int
}

type watchBroadcasts struct {
	mu       sync.Mutex
	mu       sync.Mutex
	bcasts   map[*watchBroadcast]struct{}
	watchers map[*watcher]*watchBroadcast
	updatec  chan *watchBroadcast
	donec    chan struct{}
}

// coalesce — BUG (pre-PR6916): reads wb.nextrev without acquiring wb.mu.
func (wbs *watchBroadcasts) coalesce(wb *watchBroadcast) {
	wbs.mu.Lock()
	defer wbs.mu.Unlock()
	if len(wb.receivers) >= 5 {
		return
	}
	if len(wb.receivers) == 0 {
		return
	}
	_ = wb.nextrev // BUG: racy read while concurrent writer holds wb.mu
}
