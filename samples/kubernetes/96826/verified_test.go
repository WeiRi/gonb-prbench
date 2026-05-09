package volume

import (
	"sync"
	"testing"
)

func TestRace_96826(t *testing.T) {
	const N = 50
	const ITERS = 200

	pm := &VolumePluginMgr{plugins: map[string]VolumePlugin{}}

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				pm.AddPlugin("test", nil)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_, _ = pm.GetPlugin("test")
			}
		}()
	}
	wg.Wait()
}
