package grpc

import (
	"sync"
	"testing"
)

func TestRace_1688(t *testing.T) {
	const N = 50
	const ITERS = 200

	ccb := &ccBalancerWrapper{
		subConns: make(map[*acBalancerWrapper]struct{}),
	}

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ccb.AddSC(&acBalancerWrapper{})
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = ccb.IterateSC()
			}
		}()
	}
	wg.Wait()
}
