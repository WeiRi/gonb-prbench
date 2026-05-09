package cache

import (
	"sync"
	"testing"
)

func TestRace_3985(t *testing.T) {
	it := &item{}
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 3)
	for i := 0; i < N; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				it.updateData([]byte{byte(idx)}, int64(j))
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = it.readData()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = it.readDate()
			}
		}()
	}
	wg.Wait()
}
