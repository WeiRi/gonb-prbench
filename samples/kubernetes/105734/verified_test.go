package httplog

import (
	"sync"
	"testing"
)

func TestRace_105734_HttplogConcurrent(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	for n := 0; n < N; n++ {
		rl := newRespLogger()
		wg.Add(3)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				rl.AddKeyValue("k", i)
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				rl.Addf("info")
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_ = rl.Log()
			}
		}()
	}
	wg.Wait()
}
