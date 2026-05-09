package freeport

import (
	"sync"
	"testing"
)

// TestFreeportResetRace reproduces the race condition in freeport where
// reset() modifies shared state (freePorts, pendingPorts, total) while
// the checkFreedPorts() background goroutine is still running and accessing
// those same variables.
// PR 7567: add stopCh + WaitGroup for safe goroutine lifecycle in freeport.
func TestFreeportResetRace(t *testing.T) {
	numGoroutines := 60
	iterations := 200

	for i := 0; i < iterations; i++ {
		// Initialize the freeport state (starts checkFreedPorts goroutine)
		initialize()

		var wg sync.WaitGroup

		// Concurrently take/return ports while reset happens
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Try to Take a port - this interacts with freePorts list
				// which checkFreedPorts is also accessing
				ports, err := Take(1)
				if err == nil && len(ports) > 0 {
					Return(ports)
				}
			}()
		}

		// Reset the state while goroutines are running
		// This modifies freePorts, pendingPorts, total without
		// properly stopping the checkFreedPorts goroutine
		reset()

		wg.Wait()
	}
}
