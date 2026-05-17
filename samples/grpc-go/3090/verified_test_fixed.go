package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX (PR #3090): resolverMu sync.Mutex protects resolver/done/curState.
func TestRace_3090_resolver_field(t *testing.T) {
	w := &ccResolverWrapper{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			w.resolverMu.Lock()
			w.resolver = nil
			w.resolverMu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			w.resolverMu.Lock()
			_ = w.resolver
			w.resolverMu.Unlock()
		}
	}()
	wg.Wait()
}
