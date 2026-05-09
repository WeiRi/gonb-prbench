// Race-trigger test for grpc-go-3062; see README.md for usage.

package transport

import (
	"sync"
	"testing"
)

func TestRace_PR3062_WaitOnHeaderRecvCompress(t *testing.T) {
	const N = 200
	for i := 0; i < N; i++ {
		s := &Stream{
			headerChan: make(chan struct{}),
			ctxDone:    make(chan struct{}),
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			close(s.ctxDone)
			_ = s.RecvCompress()
		}()
		go func() {
			defer wg.Done()
			s.operateHeaders("gzip")
		}()
		wg.Wait()
	}
}
