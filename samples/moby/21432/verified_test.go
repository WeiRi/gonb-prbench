package store

import (
	"fmt"
	"sync"
	"testing"

	volumedrivers "github.com/docker/docker/volume/drivers"
	volumetestutils "github.com/docker/docker/volume/testutils"
)

func TestRace_21432(t *testing.T) {
	// Bug 1 (PR 21432): List() calls setNamed() which populates
	// the cache. If Remove() deletes a volume between List's getNamed
	// and setNamed, we get a stale cache entry.
	// Bug 2 (PR 21624): create() writes to s.labels[name] without
	// globalLock. getVolume() reads s.names[name] without globalLock.
	// Concurrent Create/Get with different volume names race on the
	// shared Go map internals.

	vs, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	// Register a fake driver
	fakeDriver := volumetestutils.NewFakeDriver("fake")
	volumedrivers.Register(fakeDriver, "fake")

	var wg sync.WaitGroup
	nGoroutines := 50
	nIters := 200

	// Goroutines that create volumes with different names
	for i := 0; i < nGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				name := fmt.Sprintf("vol-%d-%d", id, j)
				vs.Create(name, "fake", nil, map[string]string{"key": "value"})
			}
		}(i)
	}

	// Goroutines that get/remove volumes (reads maps via getVolume/List)
	for i := 0; i < nGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				name := fmt.Sprintf("vol-%d-%d", id, j)
				vs.Get(name)
			}
		}(i)
	}

	wg.Wait()
}
