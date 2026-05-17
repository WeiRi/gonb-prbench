package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_nats_6983_cfg_compare(t *testing.T) {
	js := &jetStream{}
	cfg := StreamConfig{Name: "s1"}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			js.mu.Lock()
			cfg.Name = "s2"
			js.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = cfg.Name
		}
	}()
	wg.Wait()
}
