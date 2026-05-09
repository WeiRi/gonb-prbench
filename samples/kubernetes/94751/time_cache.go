package kubelet

import "time"

// timeCache mirrors the racy structure introduced in pkg/kubelet/time_cache.go
// (kubernetes PR #94751). The original used an LRU cache without a mutex,
// causing a data race on concurrent Add/Get. We keep a minimal map-based stub
// here to reproduce the same race in prod code without pulling k8s deps.

type timeCache struct {
	entries map[string]time.Time
}

func newTimeCache() *timeCache {
	return &timeCache{entries: make(map[string]time.Time)}
}

// Add stores a timestamp for uid. Non-synchronized: races with Get.
func (t *timeCache) Add(uid string, ts time.Time) {
	t.entries[uid] = ts // RACE write
}

// Get reads a timestamp for uid. Non-synchronized: races with Add.
func (t *timeCache) Get(uid string) (time.Time, bool) {
	ts, ok := t.entries[uid] // RACE read
	return ts, ok
}
