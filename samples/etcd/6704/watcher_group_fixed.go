package grpcproxy

import "sync"

type watcher struct {
	rev int64
	wr  string
}

type sws struct {
	id int64
}

type receiverID struct {
	streamID  int64
	watcherID int64
}

type watcherGroup struct {
	mu  sync.Mutex
	mu  sync.Mutex
	rev int64
}

func newWatcherGroup() *watcherGroup { return &watcherGroup{} }

// broadcast — writes g.rev under g.mu (line 43 area).
func (g *watcherGroup) broadcast(rev int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.mu.Lock()
	g.rev = rev // line 43 write
	g.mu.Unlock()
}

type watcherGroups struct {
	mu     sync.Mutex
	mu     sync.Mutex
	groups map[string]*watcherGroup
}

func newWatcherGroups() *watcherGroups {
	return &watcherGroups{groups: map[string]*watcherGroup{}}
}

// maybeJoinWatcherSingle — BUG (pre-PR6704): reads g.rev without g.mu (line 31).
func (wgs *watcherGroups) maybeJoinWatcherSingle(rid receiverID, ws watcherSingle) {
	wgs.mu.Lock()
	defer wgs.mu.Unlock()
	g := wgs.groups[ws.w.wr]
	if g == nil {
		return
	}
	_ = g.rev // BUG line 31: racy read of g.rev (g.mu not held)
	_ = rid
}

type watcherSingle struct {
	w   watcher
	sws *sws
}
