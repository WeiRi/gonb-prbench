// +build !providerless

// Regression test for kubernetes#94537
// Bug: TimedCache.getInternal: concurrent Get on same uncached key calls
// Getter multiple times (semantic race; not memory race).
//
// PANIC oracle on getter-call-count. To make race_report.txt parseable, we
// install a custom getter (calls go through prod code path:
// cache.Get → cache.getInternal → getter), and the getter itself dumps all
// goroutine stacks at the moment racy call count > 1 — at that instant the
// dumping goroutine's stack contains azure_cache.go frames.

package cache

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheNoConcurrentGet_94537(t *testing.T) {
	for iter := 0; iter < 10; iter++ {
		var calls int64
		var dumped int32
		val := &fakeDataObj{}
		getter := func(key string) (interface{}, error) {
			c := atomic.AddInt64(&calls, 1)
			if c > 1 && atomic.CompareAndSwapInt32(&dumped, 0, 1) {
				buf := make([]byte, 1<<20)
				n := runtime.Stack(buf, true)
				fmt.Fprintln(os.Stderr, "goroutine dump (PANIC oracle 94537):")
				fmt.Fprintln(os.Stderr, string(buf[:n]))
			}
			return val, nil
		}

		cache, err := NewTimedcache(fakeCacheTTL, getter)
		if err != nil {
			t.Fatalf("NewTimedcache: %v", err)
		}
		_ = time.Now()
		key := "race-key-94537"

		const G = 30
		var wg sync.WaitGroup
		wg.Add(G)
		start := make(chan struct{})
		for i := 0; i < G; i++ {
			go func() {
				defer wg.Done()
				<-start
				_, _ = cache.Get(key, CacheReadTypeDefault)
			}()
		}
		close(start)
		wg.Wait()

		if got := atomic.LoadInt64(&calls); got > 1 {
			t.Fatalf("PANIC oracle fired: iter=%d getter called %d times (expected 1)", iter, got)
		}
	}
}
