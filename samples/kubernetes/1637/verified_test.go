// Race-trigger test for kubernetes-1637; see README.md for usage.

package kubelet

import (
	"sync"
	"testing"
)

func TestRace_PR1637_RunOnceLoopVarCapture(t *testing.T) {
	const N = 256
	pods := make([]Pod, 0, N)
	for i := 0; i < N; i++ {
		pods = append(pods, Pod{
			Name: "pod-x",
			UID:  uint64(i),
			A:    i, B: i + 1, C: i + 2, D: i + 3,
		})
	}

	var wg sync.WaitGroup
	const iters = 16
	wg.Add(iters)
	for i := 0; i < iters; i++ {
		go func() {
			defer wg.Done()
			kl := &Kubelet{}
			_, _ = kl.runOnce(pods)
		}()
	}
	wg.Wait()
}
