package kubelet

import (
	"sync"
	"testing"

	cadvisorapi "github.com/google/cadvisor/info/v1"
)

func TestRace_93717(t *testing.T) {
	kl := &Kubelet{}

	var wg sync.WaitGroup
	n := 50

	// Readers: call GetCachedMachineInfo which reads kl.machineInfo without lock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				kl.GetCachedMachineInfo()
			}
		}()
	}

	// Writers: directly set kl.machineInfo without lock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				kl.machineInfo = &cadvisorapi.MachineInfo{
					NumCores:       j,
					MemoryCapacity: uint64(j * 1000),
				}
			}
		}()
	}

	wg.Wait()
}
