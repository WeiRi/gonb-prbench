package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_nats_5146_consumer_isAssigned(t *testing.T) {
	acc := &Account{Name: "test"}
	o := &consumer{acc: acc, stream: "s1", name: "c1"}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			o.mu.Lock()
			o.acc = acc
			o.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			o.mu.RLock()
			_ = o.acc
			o.mu.RUnlock()
		}
	}()
	wg.Wait()
}
