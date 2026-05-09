package metrics

import (
	"context"
	"sync"
	"testing"

	"github.com/blang/semver/v4"
)

func TestRace_133307(t *testing.T) {
	opts := &CounterOpts{
		Namespace:   "test_namespace",
		Subsystem:   "test_subsystem",
		Name:        "race_counter",
		Help:        "test counter for race detection",
		ConstLabels: map[string]string{},
	}
	c := NewCounter(opts)
	c.Create(&semver.Version{Major: 1, Minor: 30})

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Race: WithContext writes c.ctx, while Add/Inc read c.ctx
				// through withExemplar -> exemplarCounterMetric.ctx (embedded *Counter)
				c.WithContext(context.Background())
				c.Add(1)
			}
		}()
	}

	wg.Wait()
}
