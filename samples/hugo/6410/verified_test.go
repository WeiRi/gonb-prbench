package hugo6410repro

import (
	"fmt"
	"sync"
	"testing"
)

func TestRaceInitLoggers(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			DistinctErrorLog.Println(fmt.Sprintf("msg-%d", i))
		}(i)
	}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			InitLoggers()
		}()
	}
	wg.Wait()
}
