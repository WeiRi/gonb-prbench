package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: ccResolverWrapper.resolver field accessed without lock from multiple goroutines.
// PR #3090 adds resolverMu sync.Mutex to protect resolver/done/curState.
func TestRace_3090_resolver_field(t *testing.T) {
	w := &ccResolverWrapper{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			w.resolver = nil
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = w.resolver
		}
	}()
	wg.Wait()
}
