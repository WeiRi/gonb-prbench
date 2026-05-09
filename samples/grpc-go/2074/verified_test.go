// Race-trigger test for grpc-go-2074; see README.md for usage.

package transport

import (
	"sync"
	"testing"
)

func TestRace_PR2074_WriteStatusVsWriteHeader(t *testing.T) {
	const N = 50
	for i := 0; i < N; i++ {
		s := NewStream()
		s.header["pre"] = []string{"1"}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = s.WriteHeader(MD{"k1": {"v1"}, "k2": {"v2"}})
		}()
		go func() {
			defer wg.Done()
			_ = s.WriteStatus()
		}()
		wg.Wait()
	}
}
