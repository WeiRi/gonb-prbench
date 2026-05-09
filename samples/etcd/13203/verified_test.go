package clientv3

import (
	"sync"
	"testing"

	"ase/etcd-13203/resolver"
)

func TestRace_PR13203_Endpoints(t *testing.T) {
	c := &Client{
		mu:       new(sync.RWMutex),
		cfg:      Config{Endpoints: []string{"localhost:2379"}},
		resolver: resolver.New("localhost:2379"),
	}
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// Direct read of cfg.Endpoints without lock (buggy dial() pattern)
				_ = c.cfg.Endpoints
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// SetEndpoints writes cfg.Endpoints under lock (in client.go)
				c.SetEndpoints("localhost:2379")
			}
		}()
	}
	wg.Wait()
}
