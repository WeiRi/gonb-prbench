// race test for CacheBuster concurrent inner-closure invocation
package config

import (
	"io"
	"sync"
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/loggers"
)

func TestCacheBuster_CompileConfig_Race(t *testing.T) {
	cb := &CacheBuster{
		Source: `assets/.*\.(scss|css)$`,
		Target: `(css|styles|scss|sass)`,
	}

	logger := loggers.New(loggers.Options{
		DistinctLevel: logg.LevelWarn,
		StoreErrors:   true,
		StdOut:        io.Discard,
		StdErr:        io.Discard,
	})

	if err := cb.CompileConfig(logger); err != nil {
		t.Fatalf("compile: %v", err)
	}

	// Get inner closure once. compiledSource is private; use the getter via reflection-free path.
	inner := cb.compiledSource("assets/main.scss")
	if inner == nil {
		t.Fatalf("expected inner closure for matching source")
	}

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				inner("path/to/styles.css")
			}
		}()
	}
	wg.Wait()
}
