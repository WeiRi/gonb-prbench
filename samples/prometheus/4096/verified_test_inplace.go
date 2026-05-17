// In-place race test for prometheus-4096: package=scrape, uses upstream scrapePool.
// Bug: scrape.go — reload() sets sp.config, while appender() reads sp.config.SampleLimit,
// both without lock. Race on config pointer field.
// PR fix: ensure config reads happen under lock.
package scrape

import (
	"sync"
	"testing"

	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/pkg/labels"
)

// mockAppendable provides a minimal Appendable/Appender for testing
type mockAppendable struct{}

func (m mockAppendable) Appender() (storage.Appender, error) { return mockAppenderImpl{}, nil }

type mockAppenderImpl struct{}

func (m mockAppenderImpl) Add(l labels.Labels, t int64, v float64) (uint64, error) { return 0, nil }
func (m mockAppenderImpl) AddFast(l labels.Labels, ref uint64, t int64, v float64) error { return nil }
func (m mockAppenderImpl) Commit() error { return nil }
func (m mockAppenderImpl) Rollback() error { return nil }

func TestRace_4096_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	sp := &scrapePool{
		appendable: mockAppendable{},
		config: &config.ScrapeConfig{
			SampleLimit:    100,
			ScrapeInterval: 1,
			ScrapeTimeout:  1,
		},
	}

	var wg sync.WaitGroup

	// Writer goroutines: call reload() (writes sp.config, scrape/scrape.go:214)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				sp.reload(&config.ScrapeConfig{
					SampleLimit:    uint(j),
					ScrapeInterval: 1,
					ScrapeTimeout:  1,
				})
			}
		}()
	}

	// Reader goroutines: call appender() which reads sp.config.SampleLimit without lock (scrape/scrape.go:404)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = sp.appender()
			}
		}()
	}
	wg.Wait()
}
