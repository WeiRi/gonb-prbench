// Production stub for kubernetes-137804.
// Pre-PR: addPod writes to podGroupStateData fields without pgs.lock; empty
// reads them concurrently without lock.
package cache

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

type podGroupStateData struct {
	pods   map[string]struct{}
	count  int
}

type podGroupState struct {
	lock sync.Mutex
	data *podGroupStateData
}

func newPodGroupState() *podGroupState {
	return &podGroupState{
		data: &podGroupStateData{pods: map[string]struct{}{}},
	}
}

// addPod intentionally does NOT take pgs.lock — RACE with empty().
func (p *podGroupState) addPod(pod *v1.Pod) {
	p.data.pods[pod.Name] = struct{}{}
	p.data.count++
}

// empty intentionally does NOT take pgs.lock — RACE with addPod().
func (p *podGroupState) empty() bool {
	return p.data.count == 0
}
