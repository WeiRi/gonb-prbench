// consul-1214: lock.go: fix another race condition
// Original race: killChild() reads c.child WITHOUT holding c.childLock,
// while startChild() writes c.child under the lock. The mutex was not
// held in killChild when reading c.child.
// Fix: moved c.childLock.Lock() to the caller BEFORE killChild() call.
// Original diff file: command/lock.go
// Original frame hits: command/lock.go:321 (child read without lock)

package consul1214

import (
	"sync"
	"testing"
)

// TestRace reproduces the race between killChild (racy read of c.child)
// and startChild (write to c.child). 60 goroutines concurrently execute.
func TestRace(t *testing.T) {
	iterations := 300

	for i := 0; i < iterations; i++ {
		lc := &LockCommand{}

		var wg sync.WaitGroup

		// Start a child (writes c.child under lock)
		wg.Add(1)
		go func() {
			defer wg.Done()
			lc.startChild()
		}()

		// Multiple goroutines calling killChild (racy: reads c.child without lock)
		for g := 0; g < 60; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lc.killChild()
			}()
		}

		wg.Wait()
	}
}
