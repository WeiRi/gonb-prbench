package nomad

import (
	"sync"
	"testing"
)

// TestRace_PeersMapAccess reproduces the data race on s.peers and
// s.localPeers maps. The serfEventHandler goroutine (started in NewServer)
// writes to s.peers and s.localPeers under peerLock in nodeJoin().
// Test code reads these maps via len() without holding peerLock.
//
// 50+ goroutines concurrently read the peers/localPeers maps while the
// serf event handler goroutine writes to them, triggering the race detector.
func TestRace_PeersMapAccess(t *testing.T) {
	// Create two servers — this starts serfEventHandler goroutines
	s1, cleanupS1 := TestServer(t, func(c *Config) {
		c.BootstrapExpect = 1
	})
	defer cleanupS1()

	s2, cleanupS2 := TestServer(t, func(c *Config) {
		c.BootstrapExpect = 1
	})
	defer cleanupS2()

	TestJoin(t, s1, s2)

	// serfEventHandler goroutines are now running and may call nodeJoin()
	// which writes to s.peers and s.localPeers under peerLock.

	var wg sync.WaitGroup
	nWorkers := 60
	nIters := 500

	// Read goroutines: read peers/localPeers WITHOUT holding peerLock
	// (simulating the bug in serf_test.go:43-53)
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				// Read s1.peers and s1.localPeers without lock
				_ = len(s1.peers)
				_ = len(s1.localPeers)

				// Read s2.peers and s2.localPeers without lock
				_ = len(s2.peers)
				_ = len(s2.localPeers)

				// Also exercise the non-atomic Config read path
				_ = s1.config.BootstrapExpect
				_ = s2.config.BootstrapExpect
			}
		}()
	}

	wg.Wait()
}
