package clientv3

import (
	"sync"
	"sync/atomic"
	"testing"

	resolverpkg "go.etcd.io/etcd/client/v3/internal/resolver"
)

// BUG: SetEndpoints writes c.cfg.Endpoints under c.mu.Lock; unsynchronized
// readers of c.cfg.Endpoints (dial path) race.
func TestRace_etcd_13203_cfg_endpoints(t *testing.T) {
	c := &Client{mu: new(sync.RWMutex), resolver: resolverpkg.New("http://127.0.0.1:0")}
	c.cfg.Endpoints = []string{"http://127.0.0.1:0"}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			c.SetEndpoints("http://127.0.0.1:1", "http://127.0.0.1:2")
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = c.cfg.Endpoints
		}
	}()
	wg.Wait()
}
