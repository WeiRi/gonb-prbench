// Race-trigger test for etcd-7361; see README.md for usage.

package tcpproxy

import (
	"sync"
	"testing"
	"time"
)

func TestRace_TCPProxy__LoopVarCapture(t *testing.T) {
	const N = 64
	remotes := make([]*remote, 0, N)
	for i := 0; i < N; i++ {
		remotes = append(remotes, &remote{addr: "127.0.0.1:1", id: i, active: false})
	}
	tp := &TCPProxy{remotes: remotes}
	var wg sync.WaitGroup
	const iters = 16
	wg.Add(iters)
	for i := 0; i < iters; i++ {
		go func() {
			defer wg.Done()
			for k := 0; k < 32; k++ {
				tp.runMonitorOnce()
			}
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
}
