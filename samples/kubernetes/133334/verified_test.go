// Race test for k8s-133334: Sync-vs-Sync race on psuc.newCondition
// BUG: Sync reads psuc.newCondition AFTER unlock; concurrent Sync goroutines
// share the same condition pointer which syncStatus mutates (UpdatePodCondition).
// FIX: Sync DeepCopies newCondition under lock, passes local copy to syncStatus.
package apicalls

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func TestRace_133334_SyncOnly(t *testing.T) {
	// ONE shared psuc — multiple goroutines call its Sync concurrently
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{UID: "test-uid", Name: "test-pod"},
	}
	psuc := NewPodStatusPatchCall(pod, &v1.PodCondition{
		Type:    v1.PodScheduled,
		Status:  v1.ConditionTrue,
		Reason:  "initial",
		Message: "init",
	}, &framework.NominatingInfo{
		NominatedNodeName: "node-a",
		NominatingMode:    framework.ModeOverride,
	})

	var wg sync.WaitGroup
	const N = 50
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				p := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{UID: "u", Name: "p"},
				}
				_, _ = psuc.Sync(p)
			}
		}()
	}
	wg.Wait()
}
