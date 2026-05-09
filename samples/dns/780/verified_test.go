package dns

import (
	"sync"
	"testing"
)

// TestRace_dns780_started_vs_shutdown: ReadTCP reads srv.started without
// srv.lock while ShutdownContext writes started under srv.lock — fires
// under -race detector.
func TestRace_dns780_started_vs_shutdown(t *testing.T) {
	srv := NewServer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			srv.ReadTCP()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			srv.ShutdownContext()
			srv.lock.Lock()
			srv.started = true
			srv.lock.Unlock()
		}
	}()
	wg.Wait()
}
