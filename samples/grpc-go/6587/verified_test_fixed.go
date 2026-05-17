package leastrequest

import (
	"sync"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc/balancer"
)

// FIX: scWithRPCCount.numRPCs is *atomic.Int32. picker.Pick uses Load/Add.
// No race.
func TestRace_grpc_go_6587_picker(t *testing.T) {
	var n1, n2 atomic.Int32
	p := &picker{
		choiceCount: 2,
		subConns:    []scWithRPCCount{{sc: nil, numRPCs: &n1}, {sc: nil, numRPCs: &n2}},
	}
	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
				_, _ = p.Pick(balancer.PickInfo{})
			}
			atomic.StoreInt32(&done, 1)
		}()
	}
	wg.Wait()
}
