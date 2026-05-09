package grpcproxy

import (
	"context"
	"sync"
	"testing"
)

func TestRace_6916(t *testing.T) {
	wbs := &watchBroadcasts{
		bcasts:   make(map[*watchBroadcast]struct{}),
		watchers: make(map[*watcher]*watchBroadcast),
		updatec:  make(chan *watchBroadcast, 1),
		donec:    make(chan struct{}),
	}
	close(wbs.donec)
	// Build two watchBroadcasts with cancel funcs
	mkWb := func() *watchBroadcast {
		_, cancel := context.WithCancel(context.Background())
		return &watchBroadcast{
			cancel:    cancel,
			donec:     make(chan struct{}),
			receivers: make(map[*watcher]struct{}),
			nextrev:   1,
			responses: 1, // so coalesce loop enters branch
		}
	}
	wb := mkWb()
	wbswb := mkWb()
	// Add a fake receiver to wb so size()<5 and not 0 (otherwise empty() short-circuits)
	w1 := &watcher{}
	wb.receivers[w1] = struct{}{}
	close(wb.donec) // skip wb.stop hang
	close(wbswb.donec)
	wbs.bcasts[wb] = struct{}{}
	wbs.bcasts[wbswb] = struct{}{}

	const ITERS = 1000
	var wg sync.WaitGroup
	wg.Add(2)
	// Writer: under wb.mu.Lock() mutates wb.nextrev (simulating bcast loop line 86)
	go func() {
		defer wg.Done()
		for j := 0; j < ITERS; j++ {
			wb.mu.Lock()
			wb.nextrev = int64(j + 2)
			wb.mu.Unlock()
		}
	}()
	// Reader: coalesce(wb) reads wb.nextrev WITHOUT locking wb.mu (BUG)
	go func() {
		defer wg.Done()
		for j := 0; j < ITERS; j++ {
			// Reset wb.receivers so coalesce doesn't drain it forever
			wb.mu.Lock()
			if _, ok := wb.receivers[w1]; !ok {
				wb.receivers = map[*watcher]struct{}{w1: {}}
				wbs.bcasts[wb] = struct{}{}
			}
			wb.mu.Unlock()
			wbs.coalesce(wb)
		}
	}()
	wg.Wait()
}
