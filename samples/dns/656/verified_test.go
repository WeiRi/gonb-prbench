package dns

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BUG: Conn has rtt + t fields shared across concurrent exchanges → race.
func TestRace_dns_656_conn_t(t *testing.T) {
	co := &Conn{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && atomic.LoadInt32(&done) == 0; i++ {
			co.t = time.Now()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = co.t
			_ = co.rtt
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
