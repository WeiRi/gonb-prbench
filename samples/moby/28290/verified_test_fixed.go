package store

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/docker/docker/volume"
)

// FIX: VolumeStore.globalLock is sync.RWMutex. Readers take RLock; writers
// take Lock. No race.
func TestRace_moby_28290_labels_map(t *testing.T) {
	s := &VolumeStore{
		names:   map[string]volume.Volume{},
		refs:    map[string][]string{},
		labels:  map[string]map[string]string{},
		options: map[string]map[string]string{},
	}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			s.globalLock.Lock()
			s.labels["k"] = map[string]string{"a": "1"}
			s.globalLock.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 500000 && atomic.LoadInt32(&done) == 0; j++ {
			s.globalLock.RLock()
			_ = s.labels["k"]
			s.globalLock.RUnlock()
		}
	}()
	wg.Wait()
}
