// In-place race test for nomad-14120: package=nomad, uses upstream Server.
// Bug: serf.go:76,81 — serfEventHandler writes s.peers/s.localPeers under peerLock,
// but other code reads these maps without holding the lock.
// PR fix: add proper locking around map accesses.
package nomad

import (
	"sync"
	"testing"

	"github.com/hashicorp/raft"
)

func TestRace_14120_InPlace(t *testing.T) {
	const N = 60
	const ITERS = 500

	s := &Server{
		peers:      make(map[string][]*serverParts),
		localPeers: make(map[raft.ServerAddress]*serverParts),
	}

	var wg sync.WaitGroup

	// Writer goroutine: simulates serfEventHandler/reader code that writes maps without proper lock
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < ITERS; i++ {
			s.peerLock.Lock()
			s.peers["region-1"] = append(s.peers["region-1"], &serverParts{})
			s.localPeers["addr-1"] = &serverParts{}
			s.peerLock.Unlock()
		}
	}()

	// Reader goroutines: read maps WITHOUT holding lock (BUG behavior, serf.go reads bare)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = len(s.peers)      // RACE READ on bare map
				_ = len(s.localPeers) // RACE READ on bare map
			}
		}()
	}
	wg.Wait()
}
