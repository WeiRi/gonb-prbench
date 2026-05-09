package pq_1000_poc

import (
	"sync"
	"testing"
)

func TestRace_pq_1000(t *testing.T) {
	c := &conn{}
	const N = 50
	const ITERS = 500
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				c.setBad()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = c.getBad()
			}
		}()
	}
	wg.Wait()
}
