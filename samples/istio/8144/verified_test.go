package testrace

import (
	"sync/atomic"
	"testing"
)

// Minimal reproduction of istio/istio#8144: Stats() returns struct
// containing atomically-modified fields without atomic loads.
// BUG: return c.stats directly (struct copy with non-atomic reads).
// FIX: use atomic.LoadUint64 for each field.

type Stats struct {
	Evictions uint64
	Hits      uint64
	Misses    uint64
	Writes    uint64
	Removals  uint64
}

type ttlCache struct {
	stats Stats
}

func (c *ttlCache) StatsBuggy() Stats {
	return c.stats // BUGGY: struct copy races with atomic writes below
}

func (c *ttlCache) IncrementHits() {
	atomic.AddUint64(&c.stats.Hits, 1) // WRITE via atomic
}

func (c *ttlCache) IncrementWrites() {
	atomic.AddUint64(&c.stats.Writes, 1) // WRITE via atomic
}

func (c *ttlCache) IncrementEvictions() {
	atomic.AddUint64(&c.stats.Evictions, 1) // WRITE via atomic
}

func TestRace(t *testing.T) {
	c := &ttlCache{}
	done := make(chan struct{}, 1000)

	// Writers: modify stats via atomic
	for i := 0; i < 50; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				c.IncrementHits()
				c.IncrementWrites()
				c.IncrementEvictions()
			}
			done <- struct{}{}
		}()
	}

	// Readers: read stats via struct copy (non-atomic)
	for i := 0; i < 50; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				s := c.StatsBuggy() // RACE: reads fields without atomic
				_ = s.Evictions + s.Hits + s.Writes
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
