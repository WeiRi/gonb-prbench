package store

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/docker/docker/volume"
)

// BUG: VolumeStore.globalLock is sync.Mutex. list()/GetWithRef/FilterByDriver
// read s.labels[name] / s.options[name] WITHOUT lock; setLabels writes under
// Lock. Race on map.
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
	// writer (under Lock)
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			s.globalLock.Lock()
			s.labels["k"] = map[string]string{"a": "1"}
			s.globalLock.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	// BUG reader (no lock)
	go func() {
		defer wg.Done()
		for j := 0; j < 500000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = s.labels["k"]
		}
	}()
	wg.Wait()
}
