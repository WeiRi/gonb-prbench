package stats

import (
	"sync"
	"testing"
	"time"
)

func TestRace_1316(t *testing.T) {
	fs := &FollowerStats{}
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)

	// Both goroutine groups call Succ() — both modify the SAME Latency+Counts fields
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				fs.Succ(time.Millisecond)
			}
		}()
	}
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				fs.Succ(time.Microsecond)
			}
		}()
	}

	wg.Wait()
}
