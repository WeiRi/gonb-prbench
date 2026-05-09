package registrytest

import (
	"context"
	"sync"
	"testing"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
)

// TestRace_84487 reproduces the data race in NodeRegistry.WatchNodes.
// BUG: WatchNodes reads r.Err without holding r.Lock(),
// racing with other goroutines that write r.Err with the lock held.
func TestRace_84487(t *testing.T) {
	numGoroutines := 50
	iterations := 200

	for g := 0; g < numGoroutines; g++ {
		var wg sync.WaitGroup
		// Note: NodeRegistry embeds sync.Mutex
		r := &NodeRegistry{}

		// Reader: WatchNodes reads r.Err without lock (BUG)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_, _ = r.WatchNodes(context.Background(), &metainternalversion.ListOptions{})
			}
		}()

		// Writer: directly writes r.Err (race with WatchNodes read)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// Simulate error being set by another goroutine
				r.Err = nil
			}
		}()

		wg.Wait()
	}
}
