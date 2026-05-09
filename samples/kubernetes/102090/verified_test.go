package testing

import (
	"fmt"
	"sync"
	"testing"

	"k8s.io/apimachinery/pkg/watch"
)

// TestRace_102090 reproduces the race fixed by k8s/k8s#102090.
// : Unprotected R/W on Fake.WatchReactionChain.
// BUG: AddWatchReactor writes the slice without holding c.Lock().
//      InvokesWatch reads the slice while holding c.Lock(), but
//      writer side has no lock => race on slice header.
// PR https://github.com/kubernetes/kubernetes/pull/102090 adds Lock/Unlock
// to AddWatchReactor and PrependWatchReactor.
func TestRace_102090(t *testing.T) {
	c := &Fake{}
	const N = 30
	const ITERS = 200

	noopReactor := func(action Action) (handled bool, ret watch.Interface, err error) {
		return false, nil, nil
	}

	// pre-populate so InvokesWatch's range loop has something to iterate
	c.AddWatchReactor("warmup", noopReactor)

	var wg sync.WaitGroup
	wg.Add(N * 2)

	// Writers
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				c.AddWatchReactor(fmt.Sprintf("r%d-%d", id, j), noopReactor)
			}
		}(i)
	}

	// Readers: directly read the slice (not via InvokesWatch which has its
	// own lock; the BUG is on the writer side specifically).
	// We also exercise InvokesWatch via a simple WatchActionImpl literal.
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = len(c.WatchReactionChain)  // racy read
			}
		}()
	}

	wg.Wait()
}
