package scrape

import (
	"net/url"
	"sync"
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
)

func TestRace_ScrapePoolDiscoveredLabelsRaces(t *testing.T) {
	target := NewTarget(
		labels.FromStrings("__name__", "test_metric", "instance", "localhost:9090"),
		labels.FromStrings("__discovered_label__", "value"),
		url.Values{},
	)

	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	// Writer goroutines: set discovered labels concurrently
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				target.SetDiscoveredLabels(
					labels.FromStrings(
						"__discovered_label__", "value",
						"__goroutine_id__", string(rune(id)),
					),
				)
			}
		}(i)
	}

	// Reader goroutines: read discovered labels concurrently
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = target.DiscoveredLabels()
			}
		}()
	}

	wg.Wait()
}
