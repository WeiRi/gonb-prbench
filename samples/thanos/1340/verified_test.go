// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package pool

import (
	"sync"
	"testing"
)

// TestRace_1340 triggers the data race between atomic writes to BytesPool.usedTotal
// in Get() and the plain (non-atomic) read of BytesPool.usedTotal in Put().
//
// BUG: Put() calls atomic.AddUint64(&p.usedTotal, ^uint64(p.usedTotal-1)).
// The expression ^uint64(p.usedTotal-1) reads p.usedTotal as a PLAIN (non-atomic) read.
// This races with the atomic writes (atomic.AddUint64) performed by concurrent Get() calls.
//
// FIX: Replace atomic operations with a sync.Mutex protecting all reads and writes to usedTotal.
func TestRace_1340(t *testing.T) {
	p, err := NewBytesPool(10, 1000, 2, 0) // no maxTotal limit
	if err != nil {
		t.Fatal(err)
	}

	numGoroutines := 50
	iterations := 200

	var wg sync.WaitGroup

	// Goroutines doing Get+Put cycles to trigger concurrent access to usedTotal.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				b, err := p.Get(64)
				if err != nil {
					continue
				}
				p.Put(b)
			}
		}()
	}

	wg.Wait()
}
