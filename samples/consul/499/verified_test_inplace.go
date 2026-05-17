// In-place race test for consul-499: pool.go close vs ConnectionCount race
// BUG: pool.ConnectionCount() reads p.consul map without lock
// while pool.Close() writes p.shutdown flag and closes p.consulChan
package consul

import (
	"sync"
	"testing"
)

type connPool struct {
	mu       sync.Mutex
	consul   map[string]int
	shutdown bool
	conns    chan struct{}
}

func newConnPool() *connPool {
	return &connPool{consul: make(map[string]int), conns: make(chan struct{}, 10)}
}

// BUG: Close writes shutdown without lock, closes conns channel
func (p *connPool) Close() {
	p.shutdown = true
	close(p.conns)
}

// BUG: ConnectionCount reads consul map without lock
func (p *connPool) ConnectionCount() int {
	return len(p.consul)
}

func TestRace_499_InPlace(t *testing.T) {
	p := newConnPool()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); p.Close() }()
		go func() { defer wg.Done(); _ = p.ConnectionCount() }()
	}
	wg.Wait()
}
