// Race test for kubernetes-117249
// Targets ONLY the race fix.diff addresses: TopologyCache.getAllocations
// reads t.cpuRatiosByZone without lock; SetNodes writes it under lock.
package topologycache

import (
	"fmt"
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestRace117249_GetAllocations(t *testing.T) {
	cache := NewTopologyCache()

	makeNodes := func(seed int) []*v1.Node {
		return []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{v1.LabelTopologyZone: fmt.Sprintf("zone-%d", seed)}},
				Status: v1.NodeStatus{
					Allocatable: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1000m")},
					Conditions:  []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionTrue}},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{v1.LabelTopologyZone: fmt.Sprintf("zone-%d-b", seed)}},
				Status: v1.NodeStatus{
					Allocatable: v1.ResourceList{v1.ResourceCPU: resource.MustParse("2000m")},
					Conditions:  []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionTrue}},
				},
			},
		}
	}

	makeSliceInfo := func(seed int) *SliceInfo {
		return &SliceInfo{
			ServiceKey:  fmt.Sprintf("ns/svc-%d", seed),
			AddressType: discovery.AddressTypeIPv4,
			ToCreate: []*discovery.EndpointSlice{{
				Endpoints: []discovery.Endpoint{{
					Addresses:  []string{"10.1.2.3"},
					Zone:       pointer.String(fmt.Sprintf("zone-%d", seed)),
					Conditions: discovery.EndpointConditions{Ready: pointer.Bool(true)},
				}, {
					Addresses:  []string{"10.1.2.4"},
					Zone:       pointer.String(fmt.Sprintf("zone-%d-b", seed)),
					Conditions: discovery.EndpointConditions{Ready: pointer.Bool(true)},
				}},
			}},
		}
	}

	var wg sync.WaitGroup
	const iterations = 100

	// Writers: SetNodes writes t.cpuRatiosByZone under lock
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.SetNodes(makeNodes(id))
		}(i)
	}

	// Readers: AddHints calls getAllocations which (in BUG) reads cpuRatiosByZone without lock.
	// Each goroutine uses its OWN sliceInfo, so RemoveHintsFromSlices doesn't race across goroutines.
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.AddHints(makeSliceInfo(id))
		}(i)
	}

	wg.Wait()
}
