package grpc

import (
	"sync"
	"testing"
)

func TestRace_1666(t *testing.T) {
	const N = 50
	const ITERS = 200

	ac := &addrConn{}

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = ac.ReadAcbw()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ac.WriteAcbwUnlocked(nil)
			}
		}()
	}
	wg.Wait()
}
