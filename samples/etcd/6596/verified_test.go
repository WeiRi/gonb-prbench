package lease

import (
	"os"
	"sync"
	"testing"
)

func TestRace_PR6596_LeaseTTL(t *testing.T) {
	dir, be := NewTestBackend(t)
	defer os.RemoveAll(dir)
	defer be.Close()

	le := newLessor(be, minLeaseTTL)
	le.Promote(0)

	l, err := le.Grant(1, 100)
	if err != nil {
		t.Fatal(err)
	}

	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// Renew reads l.TTL without lock in buggy code (at return l.TTL, nil)
				le.Renew(l.ID)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// Direct write to l.TTL
				l.TTL = int64(j % 100)
			}
		}()
	}
	wg.Wait()
}
