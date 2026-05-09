package vm10860repro

import (
	"sync"
	"testing"
)

func TestPanicWaitGroupMisuse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered (panic = ORDER oracle hit): %v", r)
		}
	}()
	bus := newBackendURLs()
	const N = 200
	bs := make([]*backendURL, N)
	for i := 0; i < N; i++ {
		bs[i] = bus.add()
	}
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(bu *backendURL) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				bu.broken.Store(false)
				bu.setBroken()
			}
		}(bs[i])
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		bus.stopHealthChecks()
	}()
	wg.Wait()
}
