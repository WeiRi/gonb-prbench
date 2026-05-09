package lease

import (
	"sync"
	"testing"
)

func Test156275Race(t *testing.T) {
	m := &Manager{}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				m.TestingSetDisable(idx%4 == 0)
			} else {
				_ = m.WatchHandler()
			}
		}(i)
	}
	wg.Wait()
}
