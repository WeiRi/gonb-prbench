// Race-trigger test for grpc-go-4641; see README.md for usage.

package transport

import (
	"sync"
	"testing"
)

func TestRace_PR4641_RecvCompressBeforeChanClose(t *testing.T) {
	const N = 100
	for i := 0; i < N; i++ {
		s := &Stream{header: map[string][]string{}, headerChan: make(chan struct{})}
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			_ = readerObservesAfterChan(s)
		}()
		go func() {
			defer wg.Done()
			operateHeaders(s, []HF{{Name: "grpc-encoding", Value: "gzip"}})
		}()
		go func() {
			defer wg.Done()
			raceWriter(s)
		}()
		wg.Wait()
	}
}
