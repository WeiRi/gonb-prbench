package file

import (
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// TestRace_PR3738_CollectTimestampRace triggers the data race between
// TimestampCollector.Collect() iterating over Discovery.timestamps and
// Discovery.writeTimestamp()/deleteTimestamp() modifying the same map.
//
// The bug: Collect() iterates fileSD.timestamps (Discovery's map) without
// holding d.lock, while writeTimestamp()/deleteTimestamp() modify it under
// d.lock. Different locks mean no synchronization.
//
// Race detected at:
//   READ:  file.go:99  (range over fileSD.timestamps in Collect)
//   WRITE: file.go:275 (d.timestamps[filename] = timestamp)
//   WRITE: file.go:281 (delete(d.timestamps, filename))
func TestRace_PR3738_CollectTimestampRace(t *testing.T) {
	d := &Discovery{
		timestamps: make(map[string]float64),
	}
	// Create collector and add discoverer directly (bypassing NewTimestampCollector
	// to avoid the prometheus.Desc which would call MustNewConstMetric).
	tc := &TimestampCollector{
		discoverers: make(map[*Discovery]struct{}),
	}
	tc.addDiscoverer(d)
	// But Collect() requires a valid Description. Create one.
	tc.Description = prometheus.NewDesc(
		"test_metric",
		"test help",
		nil, nil,
	)

	for i := 0; i < 10; i++ {
		d.writeTimestamp("file", float64(i))
	}

	var wg sync.WaitGroup
	const numGoroutines = 100

	for g := 0; g < numGoroutines/2; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				d.writeTimestamp("file", float64(i))
				d.deleteTimestamp("file")
				d.writeTimestamp("file", float64(i))
			}
		}(g)
	}

	for g := 0; g < numGoroutines/2; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			ch := make(chan prometheus.Metric, 100)
			// Drain channel in background.
			go func() {
				for range ch {
				}
			}()
			for i := 0; i < 1000; i++ {
				tc.Collect(ch)
			}
			close(ch)
		}(g)
	}

	wg.Wait()
}
