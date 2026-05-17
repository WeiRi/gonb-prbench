// In-place race test for prometheus-885: package=retrieval, uses upstream Target.
// Bug: target.go — scrape() reads target.httpClient, honorLabels, metricRelabelConfigs
// without RLock, while Update() writes these fields under Lock.
// PR fix: add RLock in scrape() around field reads.
package retrieval

import (
	"sync"
	"testing"
	"net/url"

	"github.com/prometheus/client_golang/model"

	"github.com/prometheus/prometheus/config"
)

func TestRace_885_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	target := &Target{
		url:                  &url.URL{Host: "localhost:9090"},
		httpClient:           nil,
		honorLabels:          false,
		metricRelabelConfigs: nil,
	}

	cfg := &config.ScrapeConfig{
		HonorLabels:           true,
		MetricRelabelConfigs:  nil,
		ScrapeInterval:        1,
		ScrapeTimeout:         1,
	}

	var wg sync.WaitGroup

	// Writer goroutines: call Update() which writes fields under Lock (target.go:195-222)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				target.Update(cfg, model.LabelSet{}, model.LabelSet{})
			}
		}()
	}

	// Reader goroutines: read fields directly WITHOUT lock (simulating scrape() bug, target.go:345/371/390)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = target.httpClient
				_ = target.honorLabels
				_ = len(target.metricRelabelConfigs)
			}
		}()
	}
	wg.Wait()
}
