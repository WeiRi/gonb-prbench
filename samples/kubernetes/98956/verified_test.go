// Regression test for kubernetes#98956
package kubelet

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

func TestKillPodFollwedByIsPodPendingTermination_98956(t *testing.T) {
	pk := &podKillerWithChannel{
		podKillingCh:            make(chan *kubecontainer.PodPair, 256),
		podKillingLock:          &sync.RWMutex{},
		mirrorPodTerminationMap: make(map[string]string),
		podTerminationMap:       make(map[string]string),
		killPod: func(pod *v1.Pod, runningPod *kubecontainer.Pod, status *kubecontainer.PodStatus, gracePeriodOverride *int64) error {
			return nil
		},
	}

	const N = 100
	var raceObserved int32
	var wg sync.WaitGroup
	wg.Add(N)
	start := make(chan struct{})

	// Dump watcher: dumps all goroutines including workers still inside
	// KillPod → kubelet_pods.go frames.
	stop := make(chan struct{})
	var dumped int32
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				if atomic.LoadInt32(&raceObserved) > 0 && atomic.CompareAndSwapInt32(&dumped, 0, 1) {
					buf := make([]byte, 1<<20)
					n := runtime.Stack(buf, true)
					fmt.Fprintln(os.Stderr, "goroutine dump (PANIC oracle 98956):")
					fmt.Fprintln(os.Stderr, string(buf[:n]))
					return
				}
				runtime.Gosched()
			}
		}
	}()

	for i := 0; i < N; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-start
			uid := types.UID([]byte{byte(i + 1), 0, 0, 0})
			pod := &kubecontainer.Pod{ID: uid, Name: "p", Namespace: "default"}
			pp := &kubecontainer.PodPair{APIPod: nil, RunningPod: pod}
			pk.KillPod(pp)
			if !pk.IsPodPendingTerminationByUID(uid) {
				atomic.AddInt32(&raceObserved, 1)
			}
		}()
	}
	close(start)
	wg.Wait()
	close(stop)

	if raceObserved > 0 {
		t.Fatalf("PANIC oracle fired: %d/%d KillPod returns left pod NOT pending termination (race)",
			raceObserved, N)
	}
}
