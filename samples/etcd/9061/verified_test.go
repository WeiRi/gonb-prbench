package leasing

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRace_PR9061_WaitSession(t *testing.T) {
	lkv := &leasingKV{
		ctx:      context.Background(),
		sessionc: make(chan struct{}),
		leases:   leaseCache{mu: sync.RWMutex{}},
	}
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
				lkv.waitSession(ctx) // calls waitSession which reads lkv.sessionc without lock
				cancel()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				lkv.sessionc = make(chan struct{}) // writes to sessionc field
			}
		}()
	}
	wg.Wait()
}
