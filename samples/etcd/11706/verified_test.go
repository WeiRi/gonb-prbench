package clientv3

import (
	"context"
	"sync"
	"testing"
)

func TestRace_11706(t *testing.T) {
	md := PairsLocal("key1", "val1")

	const ITERS = 500

	var wg sync.WaitGroup
	wg.Add(8)
	for i := 0; i < 4; i++ {
		go func() {
			for j := 0; j < ITERS; j++ {
				md.Set("key1", "val1")
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < ITERS; j++ {
				ctx := NewOutgoingContextLocal(context.Background(), md)
				_ = WithRequireLeader(ctx)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
