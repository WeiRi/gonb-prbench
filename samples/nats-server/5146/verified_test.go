package server

import (
	"sync"
	"testing"
)

// TestRaceConsumerStopWithFlagsAccJs reproduces the data race in
// consumer.stopWithFlags (server/consumer.go) where o.acc.js was read
// without holding acc.mu lock, racing with account cleanup that
// modifies acc.js under acc.mu.Lock().
//
// Bug: stopWithFlags does:
//   if o.acc == nil || o.acc.js == nil || !o.acc.js.consumerAssigned(...) {
// without holding acc.mu, while account deletion does:
//   acc.mu.Lock(); acc.js = nil; acc.mu.Unlock()
//
// Fix: capture o.acc under o.mu.RLock(), then check acc.js under
// acc.mu.RLock() before accessing.
func TestRaceConsumerStopWithFlagsAccJs(t *testing.T) {
	type jsAccount struct {
		value int
	}
	type fakeAccount struct {
		mu sync.RWMutex
		js *jsAccount
	}
	type fakeConsumer struct {
		mu  sync.RWMutex
		acc *fakeAccount
	}

	acc := &fakeAccount{js: &jsAccount{value: 1}}
	o := &fakeConsumer{acc: acc}

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Readers: simulate stopWithFlags reading o.acc.js without acc.mu
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// BUG: reading o.acc and o.acc.js without holding acc.mu
				// Concurrent with writer setting acc.js = nil under Lock
				a := o.acc
				if a != nil {
					jsa := a.js
					if jsa != nil {
						_ = jsa.value
					}
				}
			}
		}()
	}

	// Writers: simulate account cleanup setting acc.js = nil under lock
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				acc.mu.Lock()
				acc.js = nil // write under lock
				acc.mu.Unlock()
				acc.mu.Lock()
				acc.js = &jsAccount{value: 1} // restore
				acc.mu.Unlock()
			}
		}()
	}

	close(ready)
	wg.Wait()
}
