package cache

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestRace_8144_InPlace reproduces the data race in ttlCache.Stats()
// (pkg/cache/ttlCache.go) where Stats() returns c.stats directly as a
// struct copy without using atomic loads, racing with atomic writes
// in functions like IncrementHits, IncrementWrites, etc.
//
// Bug: Stats() returns c.stats (struct copy without atomic loads)
// while other functions do atomic.AddUint64(&c.stats.Hits, 1) etc.
//
// Fix: use atomic.LoadUint64 for each field in the returned Stats struct.
func TestRace_8144_InPlace(t *testing.T) {
	c := &ttlCache{}

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Writers: modify stats via atomic operations
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				atomic.AddUint64(&c.stats.Hits, 1)
				atomic.AddUint64(&c.stats.Writes, 1)
				atomic.AddUint64(&c.stats.Evictions, 1)
			}
		}()
	}

	// Readers: call the BUGGY Stats() which does plain struct copy
	// without atomic loads, racing with the atomic writes above
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				s := c.Stats()
				_ = s.Evictions + s.Hits + s.Writes
			}
		}()
	}

	close(ready)
	wg.Wait()
}
