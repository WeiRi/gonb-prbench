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

// BUG: r.ResolveNow reads r.targetResolver / r.proxyResolver while r.Close
// concurrently writes them to nil — no lock.
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
			r.targetResolver = nopRes{}
			r.proxyResolver = nopRes{}
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
