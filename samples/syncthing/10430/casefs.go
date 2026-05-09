// Production stub for syncthing lib/fs/casefs.go (PR #10430).
// Pre-PR: getExpireAdd locks per-filesystem mutex but operates on shared caseCache.
package fs

import (
	"sync"
	"time"
)

type caseCache struct {
	entries map[string]time.Time
}

func newCaseCache() *caseCache {
	return &caseCache{entries: make(map[string]time.Time)}
}

type defaultRealCaser struct {
	mut   sync.Mutex // per-filesystem mutex (BUG: shared cache should have its own)
	cache *caseCache
}

// getExpireAdd locks the per-fs mutex but writes a SHARED caseCache.
// Two defaultRealCaser instances sharing the same cache hit the cache's
// map concurrently with no real protection.
func (d *defaultRealCaser) getExpireAdd(name string) time.Time {
	d.mut.Lock()
	defer d.mut.Unlock()
	if t, ok := d.cache.entries[name]; ok {
		return t
	}
	t := time.Now().Add(time.Minute)
	d.cache.entries[name] = t // RACE: concurrent map write across instances
	return t
}
