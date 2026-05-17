package plugins

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: GetAll's line `if pl, ok := storage.plugins[name]; ok` reads
// storage.plugins map WITHOUT storage.Lock. Concurrent loadWithRetry writes
// the map under Lock. Race on map.
func TestRace_moby_29893_storage_plugins(t *testing.T) {
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			storage.Lock()
			storage.plugins["k"] = &Plugin{}
			storage.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 500000 && atomic.LoadInt32(&done) == 0; j++ {
			_, _ = storage.plugins["k"]
		}
	}()
	wg.Wait()
}
