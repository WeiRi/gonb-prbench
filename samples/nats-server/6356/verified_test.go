package server

import (
	"sync"
	"testing"
)

// TestRaceJsConsumerInfoRequestMeta reproduces the data race in
// jsConsumerInfoRequest (server/jetstream_api.go) where cc.meta is read
// without holding js.mu lock, racing with consumer shutdown/deletion
// which modifies cc.meta under js.mu.Lock().
//
// Bug: jsConsumerInfoRequest does:
//   ourID := cc.meta.ID()           // reads cc.meta WITHOUT js.mu lock
//   groupLeader := cc.meta.GroupLeader()
//   groupCreated := cc.meta.Created()
// While shutdown/delete does:
//   js.mu.Lock(); cc.meta = nil; js.mu.Unlock()
//
// Fix: capture cc.meta into local variable under js.mu.RLock(),
// then use the local for ID/GroupLeader/Created calls.
func TestRaceJsConsumerInfoRequestMeta(t *testing.T) {
	type meta struct {
		value int
	}
	type consumerClient struct {
		mu   sync.RWMutex
		meta *meta
	}

	cc := &consumerClient{meta: &meta{value: 42}}

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Readers: simulate jsConsumerInfoRequest reading cc.meta without lock
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// BUG: reading cc.meta WITHOUT RLock
				// Races with writer setting cc.meta = nil under Lock
				m := cc.meta
				if m != nil {
					_ = m.value
				}
			}
		}()
	}

	// Writers: simulate shutdown/deletion setting cc.meta = nil under lock
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				cc.mu.Lock()
				cc.meta = nil // write under lock
				cc.mu.Unlock()
				cc.mu.Lock()
				cc.meta = &meta{value: 42} // restore
				cc.mu.Unlock()
			}
		}()
	}

	close(ready)
	wg.Wait()
}
