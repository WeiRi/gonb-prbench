package auth

import (
	"sync"
	"testing"
)

func TestRace_PR7664_Revision(t *testing.T) {
	as := &authStore{}
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				as.revision = uint64(j) // direct write to the unprotected field
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = as.Revision() // read via unexported method (returns as.revision directly in buggy code)
			}
		}()
	}
	wg.Wait()
}
