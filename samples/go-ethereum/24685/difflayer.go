// PR #24685 - core/state/snapshot/difflayer.go - data race on diffLayer.parent.
// Pre-fix: Parent() returns dl.parent without holding dl.lock; rebloom/flatten
// goroutines reassign parent under dl.lock.Lock(). PR fix adds RLock around
// the read in Parent().
// Production-code path: core/state/snapshot/difflayer.go (pre-fix line ~260).
package snapshot

import "sync"

type snapshot interface{}

type diffLayer struct {
	lock   sync.RWMutex
	parent snapshot
}

func newDiffLayer(p snapshot) *diffLayer { return &diffLayer{parent: p} }

// Parent — pre-fix: reads dl.parent WITHOUT lock.
// Upstream: core/state/snapshot/difflayer.go (pre-fix line ~260).
func (dl *diffLayer) Parent() snapshot {
	return dl.parent
}

// flatten reassigns dl.parent (mimicking the rebloom/flatten path that holds
// dl.lock.Lock()).
func (dl *diffLayer) flatten(np snapshot) {
	dl.lock.Lock()
	defer dl.lock.Unlock()
	dl.parent = np
}
