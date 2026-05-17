package config

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
)

// Race: BUG inner closure of CacheBuster's compiledSource captures outer
// `match` variable: `match = targetRe.MatchString(ss)`. Multiple goroutines
// calling the returned matcher race on the shared `match` variable.
// FIX uses `match := targetRe.MatchString(ss)` (local).
func TestRace_hugo_14446_cachebuster_closure_match(t *testing.T) {
	cb := &CacheBuster{
		Source: "(css|styles|scss|sass)",
		Target: "$1",
	}
	if err := cb.CompileConfig(loggers.NewDefault()); err != nil {
		t.Fatal(err)
	}
	// Call compiledSource("css") which returns the racy inner matcher.
	matcher := cb.compiledSource("css")
	if matcher == nil {
		t.Fatal("matcher is nil")
	}

	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sample := "foo"
			if idx%2 == 0 {
				sample = "bar"
			}
			for j := 0; j < 10000 && atomic.LoadInt32(&done) == 0; j++ {
				_ = matcher(sample)
			}
			atomic.StoreInt32(&done, 1)
		}(i)
	}
	wg.Wait()
}
