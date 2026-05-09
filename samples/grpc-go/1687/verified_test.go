// Race-trigger test for grpc-go-1687; see README.md for usage.

package transport

import (
	"sync"
	"testing"
)

func TestRace_PR1687_WriteStatusVsClose(t *testing.T) {
	const iters = 200
	for i := 0; i < iters; i++ {
		ht := newSHT()

		go ht.runStream()

		var wg sync.WaitGroup

		// Goroutine A: WriteStatus path — eventually does close(ht.writes).
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			_ = ht.WriteStatus("ok")
		}()

		// Goroutine B: Write path — calls do() -> ht.writes <- fn.
		// Races against A's close(ht.writes).
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			_ = ht.Write([]byte("payload"))
		}()

		wg.Wait()
		_ = ht.Close()
	}
}
