// Whitebox PoC for kubernetes-104604: data race on objectCacheItem.store between
// Get (unprotected hasSynced read) and stopIfIdle (writes store=nil).
// Production code in watch_based_manager.go.
package manager

import (
	"sync"
	"testing"
	"time"
)

func TestRace_104604(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	for n := 0; n < N; n++ {
		item := newObjectCacheItem()
		item.lastAccessTime = time.Now().Add(-1 * time.Hour)
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_ = item.Get()
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				item.stopIfIdle(time.Now(), 100*time.Millisecond)
			}
		}()
	}
	wg.Wait()
}
