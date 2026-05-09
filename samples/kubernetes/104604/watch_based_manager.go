// Pre-fix watch_based_manager.go from PR #104604.
// BUG: stopIfIdle does not check store.hasSynced() before closing the reflector,
// so a concurrent Get() that performs PollImmediate(... item.hasSynced) reads
// item.store fields while stopIfIdle may close i.stopCh and unset i.store.
// The resulting data race is on i.store / i.stopped / i.stopCh fields.
package manager

import (
	"sync"
	"sync/atomic"
	"time"
)

type fakeStore struct {
	synced atomic.Bool
}

func (s *fakeStore) hasSynced() bool { return s.synced.Load() }

type objectCacheItem struct {
	lock           sync.Mutex
	stopped        bool
	store          *fakeStore // racy read in pre-fix Get path (PollImmediate)
	stopCh         chan struct{}
	lastAccessTime time.Time
}

func newObjectCacheItem() *objectCacheItem {
	return &objectCacheItem{
		store:  &fakeStore{},
		stopCh: make(chan struct{}),
	}
}

// stopThreadUnsafe is called with i.lock held.
func (i *objectCacheItem) stopThreadUnsafe() bool {
	if i.stopped {
		return false
	}
	i.stopped = true
	// PRE-FIX: replaces store and closes channel without checking hasSynced.
	close(i.stopCh)
	i.store = nil // <-- BUG: post-stop store=nil race with Get's PollImmediate
	return true
}

// stopIfIdle (PRE-FIX): NO hasSynced guard. watch_based_manager.go:97
func (i *objectCacheItem) stopIfIdle(now time.Time, maxIdleTime time.Duration) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	// PRE-FIX path: missing && i.store.hasSynced()
	if !i.stopped && now.After(i.lastAccessTime.Add(maxIdleTime)) {
		return i.stopThreadUnsafe()
	}
	return false
}

// Get (PRE-FIX): calls store.hasSynced inside PollImmediate without holding lock.
// watch_based_manager.go:294
func (i *objectCacheItem) Get() bool {
	// pre-fix order: PollImmediate uses item.hasSynced WITHOUT lock,
	// then setLastAccessTime is set AFTER (racy with stopIfIdle reading lastAccessTime).
	// Simulate the unprotected store access.
	s := i.store          // line 294: racy read of i.store field
	if s == nil {
		return false
	}
	_ = s.hasSynced()
	i.setLastAccessTime(time.Now())
	return true
}

func (i *objectCacheItem) setLastAccessTime(t time.Time) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.lastAccessTime = t
}
