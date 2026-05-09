package reconciler

import (
	"sync"
	"testing"
)

// TestRace_serving_3357 reproduces serving PR #3357: reconciler.NewBase leaks a
// broadcaster goroutine with no StopChannel; concurrent test teardown
// vs. the goroutine reading/writing shared state -> data race.
// ALL race frames hit reconciler.go (production code).
func TestRace_serving_3357(t *testing.T) {
	for iter := 0; iter < 30; iter++ {
		b := NewBase()

		var wg sync.WaitGroup
		// Writer goroutines — call Increment on production code
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					b.Increment()
				}
			}()
		}
		// Reader goroutines — call GetCount on production code
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					_ = b.GetCount()
				}
			}()
		}
		wg.Wait()
		b.Stop()
	}
}
