package slinstance

import (
	"sync"
	"testing"
)

func Test69290Race(t *testing.T) {
	s := &session{id: 1, exp: 1}
	l := &Instance{}
	l.mu.s = s
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				l.ExtendSession(s, int64(idx))
			} else {
				_ = s.Expiration()
			}
		}(i)
	}
	wg.Wait()
}
