// In-place race test for prometheus-4066: package=scrape, uses upstream Target.
// Bug: target.go — SetDiscoveredLabels writes target.discoveredLabels,
// DiscoveredLabels reads it, both without lock. Race on labels slice.
// PR fix: add proper locking around discoveredLabels access.
package scrape

import (
	"sync"
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
)

func TestRace_4066_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	target := &Target{
		discoveredLabels: labels.Labels{{Name: "label1", Value: "value1"}},
	}

	var wg sync.WaitGroup

	// Writer goroutines: call SetDiscoveredLabels (write discoveredLabels, scrape/target.go:119-120)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				target.SetDiscoveredLabels(labels.Labels{{Name: "label", Value: "value"}})
			}
		}()
	}

	// Reader goroutines: call DiscoveredLabels (read discoveredLabels, scrape/target.go:112-114)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = target.DiscoveredLabels()
			}
		}()
	}
	wg.Wait()
}
