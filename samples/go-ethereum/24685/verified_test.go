package snapshot

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG diffLayer.Parent() reads dl.parent WITHOUT lock; writer paths
// (snapshot.go:484/505 in diff.parent = flattened/base) hold tree lock but
// not dl.lock. To mirror the production race, writer takes dl.lock and
// reader doesn't (BUG). FIX makes Parent take RLock.
func TestRace_go_ethereum_24685_difflayer_parent(t *testing.T) {
	dl := &diffLayer{}
	other := &diffLayer{}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			dl.lock.Lock()
			dl.parent = other
			dl.lock.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 500000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = dl.Parent()
		}
	}()
	wg.Wait()
}
