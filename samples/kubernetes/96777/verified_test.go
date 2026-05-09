package volumebinding

import (
	"sync"
	"testing"

	"k8s.io/kubernetes/pkg/controller/volume/scheduling"
)

func TestRace_96777(t *testing.T) {
	const N = 50
	const ITERS = 200

	state := &stateData{
		podVolumesByNode: make(map[string]*scheduling.PodVolumes),
	}

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				state.podVolumesByNode["nodeA"] = &scheduling.PodVolumes{}
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				state.podVolumesByNode["nodeB"] = &scheduling.PodVolumes{}
			}
		}()
	}
	wg.Wait()
}
