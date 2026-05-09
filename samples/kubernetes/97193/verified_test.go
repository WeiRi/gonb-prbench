package reconciler

import (
	"sync"
	"testing"
)

type fakeHandler97193 struct{}

func TestRace_97193(t *testing.T) {
	rc := &reconciler{handlers: make(map[string]interface{})}
	const N = 30
	const ITERS = 100

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			h := &fakeHandler97193{}
			for j := 0; j < ITERS; j++ {
				rc.AddHandler(string(rune('a'+id))+string(rune('0'+j%10)), h)
			}
		}(i)
	}

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				m := rc.getHandlers()
				for range m {
				}
			}
		}()
	}

	wg.Wait()
}
