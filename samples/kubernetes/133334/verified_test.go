package apicalls

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func TestRace_133334(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: Sync reads psuc.newCondition (line 134) after releasing
				// the lock (line 130), while Merge writes psuc.newCondition
				// (line 152) without any lock. This is a data race.
				pod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID:  "test-uid",
						Name: "test-pod",
					},
					Status: v1.PodStatus{},
				}
				// Create psuc with nil condition so Merge will write to it
				psuc := NewPodStatusPatchCall(pod, nil, &framework.NominatingInfo{
					NominatedNodeName: "node-a",
					NominatingMode:    framework.ModeOverride,
				})
				// Spawn goroutine that calls Merge to write newCondition
				go func() {
					psuc2 := NewPodStatusPatchCall(pod, &v1.PodCondition{
						Type:   v1.PodScheduled,
						Status: v1.ConditionTrue,
					}, &framework.NominatingInfo{
						NominatingMode: framework.ModeOverride,
					})
					psuc.Merge(psuc2)
				}()
				// Sync reads psuc.newCondition outside the lock - RACE
				psuc.Sync(pod)
			}
		}()
	}

	wg.Wait()
}
