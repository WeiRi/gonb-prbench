// Regression test for kubernetes#136226
// Bug: getEffectiveAllocatedResources started with `allocatedResources := allocatedPod.Spec.Resources`
// (alias, not a copy). Subsequent line `allocatedResources.Requests = resourcehelper.PodRequests(...)`
// then mutates pod.Spec.Resources.Requests via the alias as a side-effect. Two concurrent calls
// race on that field.
// Fix: DeepCopy first → mutations target an independent copy.
package kubelet

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestGetEffectiveAllocatedResources_NoSharedAlias_136226(t *testing.T) {
	pod := &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "c1",
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("100m"),
							v1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Limits: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("200m"),
							v1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},
			Resources: &v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("100m"),
					v1.ResourceMemory: resource.MustParse("128Mi"),
				},
				Limits: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("200m"),
					v1.ResourceMemory: resource.MustParse("256Mi"),
				},
			},
		},
	}

	const N = 200
	var wg sync.WaitGroup
	wg.Add(2)

	// Two concurrent goroutines both call getEffectiveAllocatedResources.
	// In BUG state, both share `pod.Spec.Resources` via alias and the
	// `allocatedResources.Requests = PodRequests(...)` line writes through it
	// → race on pod.Spec.Resources.Requests.
	// In FIX state DeepCopy isolates each call's writes.
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				r := getEffectiveAllocatedResources(pod)
				_ = r.Requests
				_ = r.Limits
			}
		}()
	}
	wg.Wait()
}

