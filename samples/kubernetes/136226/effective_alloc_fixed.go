// FIX version for kubernetes-136226 (PR #136226):
// DeepCopy pod.Spec.Resources before mutation so the alias doesn't leak.
package kubelet

import (
	v1 "k8s.io/api/core/v1"
)

func getEffectiveAllocatedResources(pod *v1.Pod) v1.ResourceRequirements {
	if pod.Spec.Resources == nil {
		return v1.ResourceRequirements{}
	}
	allocatedResources := pod.Spec.Resources.DeepCopy()
	return *allocatedResources
}
