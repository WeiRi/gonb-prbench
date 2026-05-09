package graphdriver

import (
	"sync"
	"testing"
)

type fakeChecker_30719 struct{}

func (fakeChecker_30719) IsMounted(path string) bool { return false }

func TestRace_30719(t *testing.T) {
	c := NewRefCounter(fakeChecker_30719{})
	const N = 30
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = c.Increment("p")
			}
		}()
	}
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = c.Decrement("p")
			}
		}()
	}
	wg.Wait()
}
