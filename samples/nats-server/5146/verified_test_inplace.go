package server

import (
	"sync"
	"testing"
)

// TestRace_5146_InPlace reproduces the data race in consumer.stopWithFlags
// (server/consumer.go) where o.acc.js is read without holding acc.mu,
// racing with account cleanup that sets acc.js = nil under acc.mu.Lock().
//
// Bug: stopWithFlags at the buggy commit reads o.acc.js without acc.mu:
//   if o.acc == nil || o.acc.js == nil || !o.acc.js.consumerAssigned(...) {
// while account deletion does acc.mu.Lock(); acc.js = nil; acc.mu.Unlock()
//
// Fix: capture o.acc under o.mu.RLock(), then check acc.js under acc.mu.RLock().
func TestRace_5146_InPlace(t *testing.T) {
	acc := NewAccount("test5146")
	acc.js = &jsAccount{}

	// Create a minimal consumer needed to reach the buggy code path
	o := &consumer{
		acc:    acc,
		closed: false,
		name:   "testconsumer",
		stream: "teststream",
	}
	o.leader.Store(true)

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Readers: call the real stopWithFlags which reads o.acc.js without acc.mu
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// Create a fresh consumer each iteration since stopWithFlags sets closed=true
				c := &consumer{
					acc:    acc,
					closed: false,
					name:   "testconsumer",
					stream: "teststream",
				}
				c.leader.Store(true)
				func() {
					defer func() { recover() }()
					c.stopWithFlags(true, false, false, false)
				}()
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
				acc.js = nil
				acc.mu.Unlock()
				acc.mu.Lock()
				acc.js = &jsAccount{}
				acc.mu.Unlock()
			}
		}()
	}

	close(ready)
	wg.Wait()
}
