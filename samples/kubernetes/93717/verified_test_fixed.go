package kubelet

import (
	"sync"
	"sync/atomic"
	"testing"

	cadvisorapi "github.com/google/cadvisor/info/v1"
)

// FIX: setCachedMachineInfo serializes writes under machineInfoLock; readers
// use GetCachedMachineInfo which takes RLock. No race.
func TestRace_kubernetes_93717_machineinfo(t *testing.T) {
	kl := &Kubelet{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			kl.setCachedMachineInfo(&cadvisorapi.MachineInfo{NumCores: j})
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000000 && atomic.LoadInt32(&done) == 0; j++ {
			_, _ = kl.GetCachedMachineInfo()
		}
	}()
	wg.Wait()
}
