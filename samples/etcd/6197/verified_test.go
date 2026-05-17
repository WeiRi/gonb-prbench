package integration

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/coreos/etcd/clientv3"
)

// BUG: package-level proxies map accessed without lock from toGRPC.
// PR #6197 adds pmu sync.Mutex. Direct map writes race with concurrent reads.
func TestRace_6197_proxies(t *testing.T) {
	c1 := &clientv3.Client{}
	c2 := &clientv3.Client{}
	proxies = map[*clientv3.Client]grpcAPI{}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			proxies[c1] = grpcAPI{}
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = proxies[c2]
		}
	}()
	wg.Wait()
}
