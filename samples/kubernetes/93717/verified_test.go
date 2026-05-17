package kubelet

import (
	"sync"
	"sync/atomic"
	"testing"

	cadvisorapi "github.com/google/cadvisor/info/v1"
)

// BUG: Kubelet.machineInfo is read by GetCachedMachineInfo and written
// directly (no lock). Concurrent write/read race.
func TestRace_kubernetes_93717_machineinfo(t *testing.T) {
	kl := &Kubelet{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			kl.machineInfo = &cadvisorapi.MachineInfo{NumCores: j}
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
