package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: ccBalancerWrapper has ccb.mu sync.RWMutex. NewSubConn/RemoveSubConn
// take Lock around map mutation. No race.
func TestRace_grpc_go_1688_subconns_map(t *testing.T) {
	ccb := &ccBalancerWrapper{
		subConns: map[*acBalancerWrapper]struct{}{},
	}
	keys := make([]*acBalancerWrapper, 16)
	for i := range keys {
		keys[i] = &acBalancerWrapper{}
	}
	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			k := keys[idx]
			for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
				ccb.mu.Lock()
				ccb.subConns[k] = struct{}{}
				ccb.mu.Unlock()
				ccb.mu.Lock()
				delete(ccb.subConns, k)
				ccb.mu.Unlock()
			}
			atomic.StoreInt32(&done, 1)
		}(i)
	}
	wg.Wait()
}
