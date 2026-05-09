package retrieval

import (
	"sync"
	"testing"
	"time"

	clientmodel "github.com/prometheus/client_golang/model"

	"github.com/prometheus/prometheus/config"
)

func TestRace_RetrievalTargetRaces(t *testing.T) {
	cfg := &config.ScrapeConfig{
		ScrapeInterval: config.Duration(1 * time.Second),
		ScrapeTimeout:  config.Duration(500 * time.Millisecond),
		HonorLabels:    false,
	}
	baseLabels := clientmodel.LabelSet{
		clientmodel.AddressLabel: "localhost:9090",
	}
	target := NewTarget(cfg, baseLabels, clientmodel.LabelSet{})

	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	// Writer goroutines: call Update() which writes to httpClient, honorLabels, metricRelabelConfigs under Lock
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				target.Update(cfg, baseLabels, clientmodel.LabelSet{})
			}
		}()
	}

	// Reader goroutines: read fields directly without RLock (as the original bug-state scrape() did)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = target.httpClient
				_ = target.honorLabels
				_ = len(target.metricRelabelConfigs)
			}
		}()
	}

	wg.Wait()
}
