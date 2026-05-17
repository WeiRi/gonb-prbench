package cache

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestRace_8214_InPlace reproduces the data race in lruCache.Stats()
// (pkg/cache/lruCache.go) where Stats() returns c.stats directly as a
// struct copy without locking, racing with SetWithExpiration which does
// atomic.AddUint64 on c.stats.Writes outside the lock.
func TestRace_8214_InPlace(t *testing.T) {
	c := &lruCache{}

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				atomic.AddUint64(&c.stats.Writes, 1)
			}
		}()
	}

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				s := c.Stats()
				_ = s.Writes
			}
		}()
	}

	close(ready)
	wg.Wait()
}
