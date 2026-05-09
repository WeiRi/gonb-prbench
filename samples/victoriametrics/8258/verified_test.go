package notifier

import (
	"sync"
	"testing"
)

func TestRace_victoriametrics_8258(t *testing.T) {
	const N = 64
	const ITERS = 1000
	shared := make([]Alert, N)
	for i := range shared {
		shared[i] = Alert{Name: "alertX", Value: float64(i)}
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < ITERS; i++ {
			_ = send(shared)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < ITERS; i++ {
			_ = send(shared)
		}
	}()
	wg.Wait()
}
