package container

import (
	"sync"
	"testing"
)

// TestRace_39645_HealthString triggers the race where Health.String()
// reads s.Health.Status directly (health.go:25) without holding s.mu,
// while SetStatus() writes s.Health.Status under s.mu.Lock() (health.go:20).
//
// Fix: changed "return s.Health.Status" to "return status" (local variable
// already read under s.mu.Lock() via Status() call at top of String).
func TestRace_39645_HealthString(t *testing.T) {
	const numGoroutines = 50
	const iterations = 200

	h := &Health{}

	var wg sync.WaitGroup

	// Writer goroutines: call SetStatus() which writes s.Health.Status under lock
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				h.SetStatus("healthy")
			}
		}()
	}

	// Reader goroutines: call String() which reads s.Health.Status WITHOUT lock
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = h.String()
			}
		}()
	}

	wg.Wait()
}
