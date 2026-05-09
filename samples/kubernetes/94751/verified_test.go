package kubelet

import (
	"sync"
	"testing"
	"time"
)

func TestRace_94751(t *testing.T) {
	const N = 50
	const ITERS = 200

	tc := newTimeCache()
	uid := "test-uid"
	tc.Add(uid, time.Now())

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				tc.Get(uid)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				tc.Add(uid, time.Now())
			}
		}()
	}
	wg.Wait()
}
