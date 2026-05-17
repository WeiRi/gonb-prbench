package integration

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/coreos/etcd/clientv3"
)

// FIX (PR #6197): pmu sync.Mutex serializes proxies map accesses.
func TestRace_6197_proxies(t *testing.T) {
	c1 := &clientv3.Client{}
	c2 := &clientv3.Client{}
	pmu.Lock()
	proxies = map[*clientv3.Client]grpcAPI{}
	pmu.Unlock()

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			pmu.Lock()
			proxies[c1] = grpcAPI{}
			pmu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			pmu.Lock()
			_ = proxies[c2]
			pmu.Unlock()
		}
	}()
	wg.Wait()
}
