package grpcproxy

import (
	"sync"
	"testing"
)

func TestRace_6906(t *testing.T) {
	wbs := &watchBroadcasts{
		bcasts:   make(map[*watchBroadcast]struct{}),
		watchers: make(map[*watcher]*watchBroadcast),
		updatec:  make(chan *watchBroadcast, 1),
		donec:    make(chan struct{}),
	}
	close(wbs.donec)

	const N = 8
	const ITERS = 500

	var wg sync.WaitGroup
	wg.Add(N * 2)
	// Writers
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				wb := &watchBroadcast{}
				wbs.mu.Lock()
				wbs.bcasts[wb] = struct{}{}
				wbs.mu.Unlock()
				wbs.mu.Lock()
				delete(wbs.bcasts, wb)
				wbs.mu.Unlock()
			}
		}()
	}
	// Readers (bounded ITERS so wg.Wait completes naturally)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS*4; j++ {
				_ = wbs.empty()
			}
		}()
	}
	wg.Wait()
}
