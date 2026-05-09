package grpc

import (
	"fmt"
	"sync"
	"testing"
)

// TestRace_1948_SliceBackingArray reproduces the data race from PR #1948.
//
// BUG: append(cc.dopts.callOptions, opts...) shares the backing array
// when callOptions has extra capacity. Concurrent Invoke() calls race.
func TestRace_1948_SliceBackingArray(t *testing.T) {
	const N = 500
	const G = 10

	var wg sync.WaitGroup
	for g := 0; g < G; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Create callOptions with extra capacity to trigger the race.
			// When a slice has cap > len, append writes to the same backing array.
			callOpts := make([]CallOption, 0, 10)
			cc := &ClientConn{dopts: dopts{callOptions: callOpts}}

			var inner sync.WaitGroup
			inner.Add(3)

			for j := 0; j < 3; j++ {
				go func(id int) {
					defer inner.Done()
					for k := 0; k < N; k++ {
						cc.Invoke(CallOption{Key: fmt.Sprintf("opt-%d-%d", id, k), Val: k})
					}
				}(j)
			}

			inner.Wait()
		}()
	}
	wg.Wait()
}
