// Whitebox PoC for kubernetes-118189: data race on TopologyCache fields
// (hintsPopulatedByService, endpointsByService) between AddHints (write)
// and HasPopulatedHints (read). The fix adds proper locking.
// Production code in topologycache.go.
package topologycache

import (
	"sync"
	"testing"

	discovery "k8s.io/api/discovery/v1"
	"k8s.io/klog/v2/ktesting"
)

func TestRace_118189_TopologyCache(t *testing.T) {
	const numGoroutines = 50
	const iterations = 200

	cache := NewTopologyCache()
	logger, _ := ktesting.NewTestContext(t)

	var wg sync.WaitGroup
	for g := 0; g < numGoroutines; g++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				si := &SliceInfo{
					ServiceKey:  "ns/svc",
					AddressType: discovery.AddressTypeIPv4,
				}
				cache.AddHints(logger, si)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				cache.HasPopulatedHints("ns/svc")
			}
		}()
	}
	wg.Wait()
}
