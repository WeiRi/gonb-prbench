// FIX version: read `initialized` under cacheStore.lock to eliminate race.
package manager

import (
	"sync"
	"testing"

	"k8s.io/client-go/tools/cache"
)

func TestRace_104604(t *testing.T) {
	const numGoroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				store := &cacheStore{Store: cache.NewStore(cache.MetaNamespaceKeyFunc)}

				var innerWg sync.WaitGroup
				innerWg.Add(2)
				go func() {
					defer innerWg.Done()
					for j := 0; j < 50; j++ {
						store.lock.Lock()
						_ = store.initialized
						store.lock.Unlock()
					}
				}()
				go func() {
					defer innerWg.Done()
					for j := 0; j < 50; j++ {
						store.lock.Lock()
						store.initialized = false
						store.lock.Unlock()
					}
				}()
				innerWg.Wait()
			}
		}()
	}
	wg.Wait()
}
