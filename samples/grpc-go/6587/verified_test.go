// Race-trigger test for grpc-go-6587; see README.md for usage.

package leastrequest

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_PR6587_LeastRequestPicker(t *testing.T) {
	const NSC = 4
	counters := make([]int32, NSC)
	scs := make([]scWithRPCCount, NSC)
	for i := 0; i < NSC; i++ {
		scs[i] = NewSCWithCount(i, &counters[i])
	}
	p := NewPicker(2, scs)

	var wg sync.WaitGroup
	const G = 8
	const N = 5000

	for g := 0; g < G; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				res, _ := p.Pick()
				res.Done(DoneInfo{})
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < N*G; i++ {
			res, _ := p.Pick()
			_ = res
		}
	}()

	wg.Wait()
	for i, c := range counters {
		_ = atomic.LoadInt32(&c)
		_ = i
	}
}
