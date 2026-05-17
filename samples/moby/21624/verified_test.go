package store

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	volumedrivers "github.com/docker/docker/volume/drivers"
	vtest "github.com/docker/docker/volume/testutils"
)

// BUG: VolumeStore.create() at L270 writes s.labels[name] = labels WITHOUT s.globalLock,
// concurrently with purge() at L97 that deletes s.labels[name] UNDER s.globalLock.
// FIX (PR #21624): wrap the create() write with s.globalLock.Lock/Unlock.
func TestRace_21624_create_vs_purge(t *testing.T) {
	volumedrivers.Register(vtest.NewFakeDriver("fake-21624"), "fake-21624")
	s, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500 && atomic.LoadInt32(&done) == 0; i++ {
			_, _ = s.Create(fmt.Sprintf("v%d", i%50), "fake-21624", nil, map[string]string{"k": "v"})
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500 && atomic.LoadInt32(&done) == 0; i++ {
			s.purge(fmt.Sprintf("v%d", i%50))
		}
	}()
	wg.Wait()
}
