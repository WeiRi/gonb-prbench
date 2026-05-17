package graphdriver

import (
	"sync"
	"sync/atomic"
	"testing"
)

type noopChecker struct{}

func (noopChecker) IsMounted(path string) bool { return false }

func TestRace_moby_30719_refcounter(t *testing.T) {
	c := NewRefCounter(noopChecker{})
	var done atomic.Bool
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5000 && !done.Load(); j++ {
				_ = c.Increment("p1")
				_ = c.Decrement("p1")
			}
			done.Store(true)
		}()
	}
	wg.Wait()
}
