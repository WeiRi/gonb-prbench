package handlers

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRace_132049_TimeoutWasCreated(t *testing.T) {
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Microsecond)
				_ = patchResource(ctx)
				cancel()
			}
		}()
	}
	wg.Wait()
}
