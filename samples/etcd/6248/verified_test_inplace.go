// TRUE_INPLACE race test for etcd-6248
// Bug: kvstore.go — getChanges() reads+writes s.changes without mutex (lines 641-642).
// concurrent getChanges() calls race on the shared s.changes slice.
package mvcc

import (
	"sync"
	"testing"

	"github.com/coreos/etcd/mvcc/mvccpb"
)

func TestRace_6248_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	s := &store{
		currentRev: revision{main: 1},
		bytesBuf8:  make([]byte, 8),
		changes:    make([]mvccpb.KeyValue, 0, 4),
	}

	var wg sync.WaitGroup
	wg.Add(N * 2)

	// All goroutines call getChanges() -> reads+replaces s.changes (kvstore.go:641-642)
	for i := 0; i < N*2; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = s.getChanges()
			}
		}()
	}

	wg.Wait()
}
