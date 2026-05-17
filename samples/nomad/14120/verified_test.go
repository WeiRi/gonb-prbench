package nomad

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: Config.Bootstrapped is int32 — accessed via plain field assignment in test/init paths
// concurrently with atomic.LoadInt32 from serf goroutines → race detector fires.
// PR #14120 moves the field to Server.bootstrapped as *atomic.Bool.
func TestRace_14120_bootstrapped(t *testing.T) {
	cfg := &Config{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			cfg.Bootstrapped = int32(i & 1)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = atomic.LoadInt32(&cfg.Bootstrapped)
		}
	}()
	wg.Wait()
}
