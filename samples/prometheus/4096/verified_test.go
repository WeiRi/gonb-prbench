package scrape

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

// nopAppendable and nopAppender are helpers for TestScrapePoolRaces.
type nopAppendable struct{}

func (a nopAppendable) Appender() (storage.Appender, error) { return nopAppender{}, nil }

type nopAppender struct{}

func (a nopAppender) Add(labels.Labels, int64, float64) (uint64, error)   { return 0, nil }
func (a nopAppender) AddFast(labels.Labels, uint64, int64, float64) error { return nil }
func (a nopAppender) Commit() error                                       { return nil }
func (a nopAppender) Rollback() error                                     { return nil }

// TestRace_PR4096_ScrapePoolReloadConfig triggers the data race where
// scrapePool.reload() sets sp.config = cfg (line 214) while concurrent
// goroutines in scrapeLoop read sp.config.SampleLimit (line 404) without
// synchronization.
//
// The bug: reload() replaces sp.config and sp.client under mtx.Lock(),
// but scrape loops access sp.config.SampleLimit outside that lock.
// When reload() is called with a new config, the old config (including
// SampleLimit) may be garbage-collected or modified while a scrape loop
// is still reading it.
//
// Fix in PR 4096: reload() correctly holds the mutex while setting config,
// but the underlying issue was that scrape loops read config fields
// without synchronizing with reload. The PR adds proper synchronization
// so the scrape loop captures config values before releasing the lock.
func TestRace_PR4096_ScrapePoolReloadConfig(t *testing.T) {
	interval, _ := model.ParseDuration("500ms")
	timeout, _ := model.ParseDuration("1s")
	newConfig := func() *config.ScrapeConfig {
		return &config.ScrapeConfig{ScrapeInterval: interval, ScrapeTimeout: timeout}
	}
	sp := newScrapePool(newConfig(), &nopAppendable{}, nil)
	tgts := []*targetgroup.Group{
		{
			Targets: []model.LabelSet{
				{model.AddressLabel: "127.0.0.1:9090"},
				{model.AddressLabel: "127.0.0.2:9090"},
				{model.AddressLabel: "127.0.0.3:9090"},
				{model.AddressLabel: "127.0.0.4:9090"},
				{model.AddressLabel: "127.0.0.5:9090"},
				{model.AddressLabel: "127.0.0.6:9090"},
				{model.AddressLabel: "127.0.0.7:9090"},
				{model.AddressLabel: "127.0.0.8:9090"},
			},
		},
	}

	active, dropped := sp.Sync(tgts)
	expectedActive, expectedDropped := len(tgts[0].Targets), 0
	if len(active) != expectedActive {
		t.Fatalf("Invalid number of active targets: expected %v, got %v", expectedActive, len(active))
	}
	if len(dropped) != expectedDropped {
		t.Fatalf("Invalid number of dropped targets: expected %v, got %v", expectedDropped, len(dropped))
	}

	var wg sync.WaitGroup
	// Concurrently reload config while scrape loops may be reading it.
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sp.reload(newConfig())
				time.Sleep(time.Duration(10) * time.Microsecond)
			}
		}()
	}
	// Also race by reading config from multiple goroutines
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = sp.appender() // reads sp.config.SampleLimit
				time.Sleep(time.Duration(10) * time.Microsecond)
			}
		}()
	}
	wg.Wait()
	sp.stop()
}
