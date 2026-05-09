package main

import (
	"sync"
	"testing"
)

// TestVolumePluginConcurrentRace runs FindPluginByName in multiple goroutines
// to trigger concurrent writes to probedPlugins map via refreshProbedPlugins.
func TestVolumePluginConcurrentRace(t *testing.T) {
	prober := &fakeProber{
		plugin: &testPlugin{name: "testPlugin"},
	}

	pm := NewVolumePluginMgr(prober, nil)

	// Init with empty plugin list — probedPlugins NOT populated yet.
	_ = pm.InitPlugins([]VolumePlugin{}, prober, nil)

	var wg sync.WaitGroup
	n := 50

	for g := 0; g < n; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				_, _ = pm.FindPluginByName("testPlugin")
			}
		}()
	}

	wg.Wait()
}
