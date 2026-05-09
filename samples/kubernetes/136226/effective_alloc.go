// Production stub for kubernetes-136226.
// Pre-PR: getEffectiveAllocatedResources started with `allocatedResources :=
// allocatedPod.Spec.Resources` (alias, not copy). Subsequent assignment to
// allocatedResources.Requests then mutates pod.Spec.Resources.Requests via
// the alias — concurrent calls race on that field.
package kubelet

import (
	v1 "k8s.io/api/core/v1"
)

// getEffectiveAllocatedResources without DeepCopy — RACE.
func getEffectiveAllocatedResources(pod *v1.Pod) v1.ResourceRequirements {
	if pod.Spec.Resources == nil {
		return v1.ResourceRequirements{}
	}
	allocatedResources := *pod.Spec.Resources // shallow copy: maps are aliased
	// Now mutate Requests — this writes through the aliased map.
	if pod.Spec.Resources.Requests != nil {
		allocatedResources.Requests = pod.Spec.Resources.Requests // alias
		pod.Spec.Resources.Requests = allocatedResources.Requests // RACE: write
	}
	return allocatedResources
}
