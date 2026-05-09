// Production stub for kubernetes-93717.
// Pre-PR: Kubelet.machineInfo is read by GetCachedMachineInfo and written by
// direct assignment without any synchronization, racing across goroutines.
package kubelet

import (
	cadvisorapi "github.com/google/cadvisor/info/v1"
)

type Kubelet struct {
	// machineInfo is unguarded in pre-fix version (RACE).
	machineInfo *cadvisorapi.MachineInfo
}

// GetCachedMachineInfo reads kl.machineInfo without any lock (RACE).
func (kl *Kubelet) GetCachedMachineInfo() (*cadvisorapi.MachineInfo, error) {
	return kl.machineInfo, nil
}
