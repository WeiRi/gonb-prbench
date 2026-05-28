// Minimal reproduction stub for etcd-19600
// Bug: In cancelWatcher, wa.compacted branch comes BEFORE wa.ch == nil,
// so double-decrement on close/cancel race for compacted watchers.
// Race: watcherGauge is a plain int64; Inc() called outside store mutex
// (in newWatcher after s.mu.Unlock()), while Dec() called inside s.mu.Lock()
// (in cancelWatcher). Concurrent Inc+Dec = data race.
//
// Pre-fix version.

package main

import (
	"sync"
	"sync/atomic"
)

var watcherGauge int64

type watcher struct {
	key       []byte
	compacted bool
	ch        chan struct{}
	victim    bool
}

type watchableStore struct {
	mu       sync.Mutex
	unsynced map[*watcher]struct{}
	synced   map[*watcher]struct{}
}

func newWatchableStore() *watchableStore {
	return &watchableStore{
		unsynced: make(map[*watcher]struct{}),
		synced:   make(map[*watcher]struct{}),
	}
}

func (s *watchableStore) newWatcher(compacted bool) (*watcher, func()) {
	wa := &watcher{
		ch:        make(chan struct{}),
		compacted: compacted,
	}
	// watcherGauge.Inc() called outside store mutex (mirrors real code)
	atomic.AddInt64(&watcherGauge, 1)
	return wa, func() { s.cancelWatcher(wa) }
}

func (s *watchableStore) deleteFrom(m map[*watcher]struct{}, wa *watcher) bool {
	if _, ok := m[wa]; ok {
		delete(m, wa)
		return true
	}
	return false
}

// cancelWatcher - PRE-FIX: wa.compacted BEFORE wa.ch == nil.
// Bug: two calls (cancel+close) both hit wa.compacted, double-decrement.
func (s *watchableStore) cancelWatcher(wa *watcher) {
	for {
		s.mu.Lock()
		if s.deleteFrom(s.unsynced, wa) {
			atomic.AddInt64(&watcherGauge, -1)
			break
		} else if s.deleteFrom(s.synced, wa) {
			atomic.AddInt64(&watcherGauge, -1)
			break
		} else if wa.compacted {
			// PRE-FIX BUG: checked before wa.ch == nil
			atomic.AddInt64(&watcherGauge, -1)
			break
		} else if wa.ch == nil {
			break
		}

		if !wa.victim {
			s.mu.Unlock()
			panic("watcher not in any group")
		}
		s.mu.Unlock()
		break
	}

	wa.ch = nil
	s.mu.Unlock()
}
