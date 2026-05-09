package topologycache

import (
	"fmt"
	"sync"
	"testing"
)

func TestRace_118189_TopologyCache(t *testing.T) {
	const N = 100
	var wg sync.WaitGroup
	tc := NewTopologyCache()
	for n := 0; n < N; n++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				tc.AddHints(fmt.Sprintf("svc%d", n%5), "ipv4")
			}
		}(n)
		go func(n int) {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				_ = tc.HasPopulatedHints(fmt.Sprintf("svc%d", n%5))
			}
		}(n)
	}
	wg.Wait()
}
