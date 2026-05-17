package cache

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestRace_137804_PodGroupStateMaps(t *testing.T) {
	pgs := newPodGroupState()
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{UID: types.UID("p1"), Name: "p1"},
	}
	var wg sync.WaitGroup
	const N = 200
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			pgs.addPod(pod)
			pgs.deletePod(pod.UID)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = pgs.empty()
		}
	}()
	wg.Wait()
}
