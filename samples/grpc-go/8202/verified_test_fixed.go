package delegatingresolver

import (
	"sync"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc/resolver"
)

type nopRes struct{}

func (nopRes) ResolveNow(resolver.ResolveNowOptions) {}
func (nopRes) Close()                                {}

// FIX: r.childMu serializes access to targetResolver / proxyResolver.
// ResolveNow + Close both take childMu.
func TestRace_grpc_go_8202_resolver_close_resolve(t *testing.T) {
	r := &delegatingResolver{
		targetResolver: nopRes{},
		proxyResolver:  nopRes{},
	}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			r.ResolveNow(resolver.ResolveNowOptions{})
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000 && atomic.LoadInt32(&done) == 0; j++ {
			r.childMu.Lock()
			r.targetResolver = nopRes{}
			r.proxyResolver = nopRes{}
			r.childMu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
