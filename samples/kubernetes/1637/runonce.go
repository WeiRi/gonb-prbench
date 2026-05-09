package kubelet

import "sync/atomic"

// Pod is a stripped-down stub of v1.Pod for the BUG-state runOnce reproduction.
// Original BUG site: pkg/kubelet/runonce.go (PR #1637).
type Pod struct {
	Name string
	UID  uint64
	A, B, C, D int
}

// RunPodResult mirrors the original RunPodResult struct.
type RunPodResult struct {
	Pod  *Pod
	Info int
	Err  error
}

// Kubelet stub.
type Kubelet struct {
	counter atomic.Int64
}

// runPod reads several fields of pod (matches the original behaviour).
func (kl *Kubelet) runPod(pod Pod) (int, error) {
	kl.counter.Add(1)
	// touch all fields so race-detector instrumentation fires on field reads
	return pod.A + pod.B + pod.C + pod.D + int(pod.UID), nil
}

// runOnce — BUG state of the original: closure captures the loop variable
// `pod` from `for _, pod := range pods` while the range loop overwrites
// it on the next iteration => data race on `pod`.
func (kl *Kubelet) runOnce(pods []Pod) (results []RunPodResult, err error) {
	ch := make(chan RunPodResult, len(pods))
	for _, pod := range pods { // line 41 (range loop overwrites pod)
		go func() {            // line 42 (goroutine captures pod by ref)
			info, e := kl.runPod(pod)             // line 43 (read pod fields)
			ch <- RunPodResult{Pod: &pod, Info: info, Err: e} // line 44 (read pod for &pod)
		}()
	}
	for i := 0; i < len(pods); i++ {
		res := <-ch
		results = append(results, res)
	}
	return results, nil
}
