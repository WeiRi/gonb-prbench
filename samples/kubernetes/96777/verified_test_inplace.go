// Regression test for kubernetes-96777
// PR: https://github.com/kubernetes/kubernetes/pull/96777
// Race: stateData.podVolumesByNode map accessed without proper locking.
//        cs.Lock() locks the CycleState, but Filter can be called on
//        different CycleState instances sharing the same stateData
//        (which happens when the scheduler clones CycleState for parallel Filter).
//        The fix adds a sync.Mutex to stateData to protect the map.
// Uses real upstream types: VolumeBinding, stateData, CycleState, NodeInfo.
// Uses real upstream method: VolumeBinding.Filter.
package volumebinding

import (
	"context"
	"fmt"
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/controller/volume/scheduling"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// mockBinder implements scheduling.SchedulerVolumeBinder for testing
type mockBinder struct{}

func (m *mockBinder) GetPodVolumes(pod *v1.Pod) ([]*v1.PersistentVolumeClaim, []*v1.PersistentVolumeClaim, []*v1.PersistentVolumeClaim, error) {
	return nil, nil, nil, nil
}

func (m *mockBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*v1.PersistentVolumeClaim, node *v1.Node) (*scheduling.PodVolumes, scheduling.ConflictReasons, error) {
	return &scheduling.PodVolumes{}, nil, nil
}

func (m *mockBinder) AssumePodVolumes(assumedPod *v1.Pod, nodeName string, podVolumes *scheduling.PodVolumes) (bool, error) {
	return true, nil
}

func (m *mockBinder) RevertAssumedPodVolumes(podVolumes *scheduling.PodVolumes) {}

func (m *mockBinder) BindPodVolumes(assumedPod *v1.Pod, podVolumes *scheduling.PodVolumes) error {
	return nil
}

func TestRace_96777_InPlace(t *testing.T) {
	pl := &VolumeBinding{
		Binder: &mockBinder{},
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "default"},
	}

	const N = 50
	const ITERS = 100

	for trial := 0; trial < ITERS; trial++ {
		// Create stateData with the shared map (the race target)
		sd := &stateData{
			boundClaims:      nil,
			claimsToBind:     nil,
			podVolumesByNode: make(map[string]*scheduling.PodVolumes),
		}

		// Create two CycleState instances that share the same stateData
		// (simulating CycleState duplication in the scheduler framework)
		cs1 := framework.NewCycleState()
		cs1.Write(stateKey, sd)
		// Manually create a second CycleState with the SAME stateData
		// (Clone() would create a new CycleState but share the same sd pointer)
		cs2 := framework.NewCycleState()
		cs2.Write(stateKey, sd)

		var wg sync.WaitGroup
		wg.Add(N * 2)

		for i := 0; i < N; i++ {
			// Writer: calls Filter which writes to state.podVolumesByNode[node.Name]
			// Uses cs1 as the CycleState
			go func(id int) {
				defer wg.Done()
				node := &v1.Node{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("node-%d", id%10)},
				}
				nodeInfo := framework.NewNodeInfo()
				nodeInfo.SetNode(node)
				// Filter writes: state.podVolumesByNode[node.Name] = podVolumes
				// cs1.Lock() does NOT protect sd because cs2.Lock() locks a different mutex
				pl.Filter(context.Background(), cs1, pod, nodeInfo)
			}(i)

			// Writer: calls Filter concurrently with cs2 (different CycleState mutex)
			go func(id int) {
				defer wg.Done()
				node := &v1.Node{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("node-%d", id%10)},
				}
				nodeInfo := framework.NewNodeInfo()
				nodeInfo.SetNode(node)
				// cs2.Lock() only locks cs2, not cs1 — both write to the same sd map
				pl.Filter(context.Background(), cs2, pod, nodeInfo)
			}(i)
		}
		wg.Wait()
	}
}
