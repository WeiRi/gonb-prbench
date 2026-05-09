package dial

import (
	"sync"
	"testing"
)

// TestRace_23434_conn_flags: register-side read of c.flags races with
// concurrent SetFlag write — fires under -race.
func TestRace_23434_conn_flags(t *testing.T) {
	c := &conn{}
	d := newScheduler()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.SetFlag(connFlag(i))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			d.Register(i, c)
		}
	}()
	wg.Wait()
}
