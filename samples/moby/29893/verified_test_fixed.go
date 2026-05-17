package plugins

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: GetAll wraps the read with storage.Lock/Unlock. No race.
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
			storage.Lock()
			_, _ = storage.plugins["k"]
			storage.Unlock()
		}
	}()
	wg.Wait()
}
