package tracing

import (
	"sync"
	"testing"
)

func Test65339Race(t *testing.T) {
	s := &crdbSpan{operation: "init"}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				s.SetOperationName("op")
			} else {
				_ = s.GetRecording()
			}
		}(i)
	}
	wg.Wait()
}
