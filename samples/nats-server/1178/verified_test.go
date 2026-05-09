// VERIFIED race reproducer for nats-server PR #1178
// "Fixed data race" (resolves issue #1176)
// https://github.com/nats-io/nats-server/pull/1178
//
// Original racy file: server/events.go (debugSubscribers, line 1260 pre-fix)
// Fix: replace `nsubs` bare read with `atomic.LoadInt32(&nsubs)` in deferred dispatch.
//
// Race recipe (whitebox):
//   - Sublist match callbacks atomically increment nsubs from goroutines.
//   - The deferred internal-account-msg dispatcher reads nsubs without atomic.
//   - Pre-fix: bare int32 read races with atomic.AddInt32.
//
// 200 invocations x 4 add-goroutines, count=10, golang:1.21 -race => race in 1st iter.
package buggy

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRaceDebugSubscribers(t *testing.T) {
	s := &Server{results: map[string]int32{}}
	const N = 200
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.debugSubscribers(fmt.Sprintf("r%d", i))
		}(i)
	}
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
}
