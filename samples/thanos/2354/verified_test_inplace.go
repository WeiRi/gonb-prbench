// Race-trigger test for thanos-2354; self-contained package main.
// Bug: pkg/block/fetcher.go — multiple concurrent Fetch() calls cause
// concurrent map read/write on MetaFetcher.cached map. loadMeta() reads
// s.cached while Fetch() writes s.cached = cached at the end.
// PR fix: introduce BaseFetcher with singleflight.Group to serialize.
package block

import (
	"sync"
	"testing"
)

// Meta_2354 mirrors upstream metadata.Meta (simplified).
type Meta_2354 struct {
	ULID string
}

// MetaFetcher_2354 mirrors upstream MetaFetcher with BUG: cached map
// read/written concurrently without synchronization.
type MetaFetcher_2354 struct {
	cached map[string]*Meta_2354
}

func (f *MetaFetcher_2354) loadMeta_BUG(id string) *Meta_2354 {
	return f.cached[id] // RACE READ on bare map (pkg/block/fetcher.go loadMeta)
}

func (f *MetaFetcher_2354) Fetch_BUG() {
	// Build new cache
	newCached := make(map[string]*Meta_2354)
	for i := 0; i < 30; i++ {
		id := string(rune('a' + i%26)) + string(rune('A'+i%26))
		newCached[id] = &Meta_2354{ULID: id}
	}
	// RACE WRITE: replace cached map without sync (pkg/block/fetcher.go Fetch end)
	f.cached = newCached
}

func TestRace_2354_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200
	f := &MetaFetcher_2354{
		cached: make(map[string]*Meta_2354),
	}
	f.cached["init"] = &Meta_2354{ULID: "init"}
	var wg sync.WaitGroup

	// Writer goroutines: call Fetch which writes f.cached
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				f.Fetch_BUG()
			}
		}()
	}

	// Reader goroutines: read f.cached via loadMeta
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = f.loadMeta_BUG("init")
				_ = f.loadMeta_BUG("aA")
			}
		}()
	}
	wg.Wait()
}
