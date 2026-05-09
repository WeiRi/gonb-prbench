package cache

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRace_137804(t *testing.T) {
	pgs := newPodGroupState()

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: addPod writes to podGroupStateData without pgs.lock,
				// while empty() reads podGroupStateData without pgs.lock.
				// This is a data race on the shared podGroupStateData fields.
				pod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-pod",
					},
				}
				pgs.addPod(pod)
				pgs.empty()
			}
		}()
	}

	wg.Wait()
}
