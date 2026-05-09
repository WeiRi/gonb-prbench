// consul-499: Fix potential race condition on shutdown (pool.reap/server.handleConsulConn)
// Original race: ConnPool.shutdown bool read without lock in reap() (pool.go:361)
// Fix: replaced racy bool check with channel select (<-p.shutdownCh).
// Original diff files: consul/pool.go, consul/rpc.go
// Original frame hits: consul/pool.go:361, consul/rpc.go:152

package consul499

import (
	"sync"
	"testing"
)

// ConnPool replicates the racy ConnPool struct from consul/pool.go
type ConnPool struct {
	mu         sync.Mutex
	pool       map[string]*Conn
	shutdown   bool
	shutdownCh chan struct{}
}

type Conn struct {
	refCount int32
}

func NewConnPool() *ConnPool {
	p := &ConnPool{
		pool:       make(map[string]*Conn),
		shutdownCh: make(chan struct{}),
	}
	return p
}

// Shutdown writes shutdown=true with lock held
func (p *ConnPool) Shutdown() {
	p.mu.Lock()
	if p.shutdown {
		p.mu.Unlock()
		return
	}
	p.shutdown = true
	p.mu.Unlock()
	close(p.shutdownCh)
}

// reap reads p.shutdown WITHOUT the lock (RACY - original line pool.go:361)
func (p *ConnPool) reap() {
	for !p.shutdown {
		select {
		case <-p.shutdownCh:
			return
		default:
		}
		p.mu.Lock()
		_ = len(p.pool)
		p.mu.Unlock()
	}
}

// TestRace reproduces the race on ConnPool.shutdown bool.
// 60 goroutines call Shutdown() concurrently while reap() reads shutdown.
func TestRace(t *testing.T) {
	iterations := 300

	for i := 0; i < iterations; i++ {
		p := NewConnPool()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.reap()
		}()

		for g := 0; g < 60; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.Shutdown()
			}()
		}

		wg.Wait()
	}
}
