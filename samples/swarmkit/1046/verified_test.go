package agent

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRace_swarmkit_1046(t *testing.T) {
	const N = 64
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
			defer cancel()
			done := make(chan struct{})
			close(done)
			RunManagerLoopBUG(ctx, done)
		}()
	}
	wg.Wait()
}
