package metrics

import (
	"sync"
	"testing"

	"github.com/blang/semver/v4"
)

func TestRace_132063(t *testing.T) {
	// Set up label allow list for the metric so that the Do body actually
	// writes to LabelValueAllowLists (triggering the race).
	SetLabelAllowListFromCLI(map[string]string{
		"test_namespace_test_subsystem_race_counter,a": "a_val",
		"test_namespace_test_subsystem_race_counter,b": "b_val",
	})

	// Create a single shared CounterVec. The race is on concurrent
	// read of LabelValueAllowLists (if != nil at counter.go:215)
	// vs write inside sync.Once (v.LabelValueAllowLists = allowList at counter.go:221).
	opts := &CounterOpts{
		Namespace:   "test_namespace",
		Subsystem:   "test_subsystem",
		Name:        "race_counter",
		Help:        "test counter for race detection",
		ConstLabels: map[string]string{},
	}
	cv := NewCounterVec(opts, []string{"a", "b"})
	cv.Create(&semver.Version{Major: 1, Minor: 30})

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: v.LabelValueAllowLists read races with write inside Do.
				// After first successful Do, only reads happen (no race),
				// but initial concurrent calls trigger the race detector.
				cv.WithLabelValues("a_val", "b_val")
			}
		}()
	}

	wg.Wait()
}
