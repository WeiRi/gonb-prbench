package quota

import (
	"sync"
	"testing"
)

func TestRace_39644(t *testing.T) {
	const N = 50
	const ITERS = 200

	q := NewControl()

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				q.SetQuota("target")
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_, _ = q.GetQuota("target")
			}
		}()
	}
	wg.Wait()
}
