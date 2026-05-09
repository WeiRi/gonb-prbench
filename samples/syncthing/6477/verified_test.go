package syncthing6477repro

import (
	"sync"
	"testing"
)

func TestRaceStoppedField(t *testing.T) {
	const N = 500
	var wg sync.WaitGroup
	s := newService()
	for i := 0; i < N; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); s.Serve() }()
		go func() { defer wg.Done(); s.Stop() }()
	}
	wg.Wait()
}
