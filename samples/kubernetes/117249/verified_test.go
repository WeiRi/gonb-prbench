// PR: https://github.com/kubernetes/kubernetes/pull/117249
// Fix: Move the lock acquisition before the nil check on cpuRatiosByZone
// in getAllocations to prevent data race with SetNodes.
package topologycache

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestRace117249TopologyCache(t *testing.T) {
	cache := NewTopologyCache()

	nodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{v1.LabelTopologyZone: "zone-a"}},
			Status: v1.NodeStatus{
				Allocatable: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1000m")},
				Conditions:  []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionTrue}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{v1.LabelTopologyZone: "zone-b"}},
			Status: v1.NodeStatus{
				Allocatable: v1.ResourceList{v1.ResourceCPU: resource.MustParse("2000m")},
				Conditions:  []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionTrue}},
			},
		},
	}

	sliceInfo := &SliceInfo{
		ServiceKey:  "ns/svc",
		AddressType: discovery.AddressTypeIPv4,
		ToCreate: []*discovery.EndpointSlice{{
			Endpoints: []discovery.Endpoint{{
				Addresses:  []string{"10.1.2.3"},
				Zone:       pointer.String("zone-a"),
				Conditions: discovery.EndpointConditions{Ready: pointer.Bool(true)},
			}, {
				Addresses:  []string{"10.1.2.4"},
				Zone:       pointer.String("zone-b"),
				Conditions: discovery.EndpointConditions{Ready: pointer.Bool(true)},
			}},
		}},
	}

	var wg sync.WaitGroup
	const iterations = 100
	wg.Add(iterations * 2)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			// SetNodes writes cpuRatiosByZone
			cache.SetNodes(nodes)
		}()
		go func() {
			defer wg.Done()
			// AddHints calls getAllocations which reads cpuRatiosByZone
			cache.AddHints(sliceInfo)
		}()
	}

	wg.Wait()
}
