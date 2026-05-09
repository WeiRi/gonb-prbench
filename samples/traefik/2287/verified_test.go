package healthcheck

import (
	"context"
	"sync"
	"testing"
)

// TestRace_TR2287_SetBackendsRace triggers the data race on the Backends map
// in HealthCheck. SetBackendsConfiguration writes hc.Backends and then
// iterates over it without any mutex protection. Concurrent calls to
// SetBackendsConfiguration can race on the map assignment and iteration.
//
// Bug: HealthCheck has no mutex. SetBackendsConfiguration assigns
// hc.Backends (line 78) and iterates over it (line 85) without
// synchronization.
//
// Fix: add mutex to protect Backends map in SetBackendsConfiguration.
func TestRace_TR2287_SetBackendsRace(t *testing.T) {
	hc := newHealthCheck()
	ctx := context.Background()

	var wg sync.WaitGroup
	const numGoroutines = 100

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				backends := make(map[string]*BackendHealthCheck)
				backends["test"] = NewBackendHealthCheck(Options{})
				// SetBackendsConfiguration writes hc.Backends (line 78)
				// and iterates (line 85) without mutex.
				hc.SetBackendsConfiguration(ctx, backends)
			}
		}(g)
	}

	wg.Wait()
}
