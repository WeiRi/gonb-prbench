package admission

import (
	"sync"
	"testing"
)

func Test131292Race(t *testing.T) {
	q := NewWorkQueue()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				q.Push(idx)
			} else {
				_ = q.Admit()
				_ = q.Granted()
			}
		}(i)
	}
	wg.Wait()
}
