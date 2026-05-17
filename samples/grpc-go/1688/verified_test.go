package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: ccBalancerWrapper has no `mu` field; concurrent writes to
// ccb.subConns map (in NewSubConn / RemoveSubConn) race.
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
				ccb.subConns[k] = struct{}{}
				delete(ccb.subConns, k)
			}
			atomic.StoreInt32(&done, 1)
		}(i)
	}
	wg.Wait()
}
