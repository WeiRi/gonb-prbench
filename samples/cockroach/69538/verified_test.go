package metric

import (
	"sync"
	"testing"
)

// TestRace_69538 reproduces cockroach pkg/util/metric/registry.go race:
// AddMetric writes to r.metrics map, while WriteMetricsMetadata / MarshalJSON
// (modeled here as Each / Get) iterate / read the same map without holding the
// lock. Bug: read paths don't take r.Lock().
func TestRace_69538(t *testing.T) {
	const (
		numGoroutines = 60
		numIterations = 100
	)

	for iter := 0; iter < numIterations; iter++ {
		r := NewRegistry()

		var wg sync.WaitGroup

		// Reader goroutines: iterate / lookup map (BUG: no lock)
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.Each(func(name string, val int64) {})
			}()
		}

		// Writer goroutines: add to map (BUG: no lock)
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				r.AddMetric("metric", int64(gid))
			}(g)
		}

		wg.Wait()
	}
}
