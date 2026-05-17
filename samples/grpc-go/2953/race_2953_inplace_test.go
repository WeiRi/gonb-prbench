package transport

import (
	"sync"
	"testing"
)

func TestRace_2953_InPlace(t *testing.T) {
	f := &trInFlow{}
	const N = 50
	const ITERS = 100
	for trial := 0; trial < ITERS; trial++ {
		var wg sync.WaitGroup
		wg.Add(N * 2)
		for i := 0; i < N; i++ {
			go func() {
				defer wg.Done()
				f.onData(100)
			}()
			go func() {
				defer wg.Done()
				f.newLimit(1000)
			}()
		}
		wg.Wait()
	}
}
