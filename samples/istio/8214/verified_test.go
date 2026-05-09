package testrace

import (
	"sync/atomic"
	"testing"
)

// Minimal reproduction of istio/istio#8214: Stats() reads without lock
// while SetWithExpiration writes fields.
// BUG: Stats() does plain return of struct fields without any lock.
// FIX: use sync.RWMutex, Stats() acquires RLock before reading stats.

type Stats struct {
	Writes   uint64
	Hits     uint64
	Evictions uint64
}

type lruCache struct {
	stats Stats
}

func (c *lruCache) SetWithExpiration() {
	// Simulates write to stats
	atomic.AddUint64(&c.stats.Writes, 1) // WRITE
}

func (c *lruCache) StatsBuggy() Stats {
	return c.stats // BUGGY: reads without RLock
}

func TestRace(t *testing.T) {
	c := &lruCache{}
	done := make(chan struct{}, 200)

	// Concurrent writers
	for i := 0; i < 50; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				c.SetWithExpiration()
			}
			done <- struct{}{}
		}()
	}

	// Concurrent readers (no lock)
	for i := 0; i < 50; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				s := c.StatsBuggy()
				_ = s
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
