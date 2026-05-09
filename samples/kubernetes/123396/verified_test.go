package mux

import (
	"fmt"
	"sync"
	"testing"
)

func TestRace_123396_PathRecorder(t *testing.T) {
	m := NewPathRecorderMux()
	const N = 200
	var wg sync.WaitGroup
	for n := 0; n < N; n++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				m.Register(fmt.Sprintf("/p%d", n*100+i))
			}
		}(n)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_ = m.ListedPaths()
			}
		}()
	}
	wg.Wait()
}
