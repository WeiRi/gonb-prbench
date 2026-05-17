package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_nats_6356_cc_meta(t *testing.T) {
	cc := &jetStreamCluster{}
	js := &jetStream{cluster: cc}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			js.mu.Lock()
			cc.meta = nil
			js.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = cc.meta
		}
	}()
	wg.Wait()
}
