package v2

import (
	"sync"
	"testing"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

// BUG: handler.watch's goroutine has `case <-closeChan: stopWatchChan <- true`.
// After watch's deferred close(stopWatchChan), the goroutine's send panics:
// send on closed channel. PR #503 splits with a separate stopWrapChan.
func TestRace_503_panic_send_on_closed_stopWatchChan(t *testing.T) {
	// Drive h.client.Get to fail quickly via unreachable endpoint.
	cli := etcd.NewClient([]string{"http://127.0.0.1:1"})
	h := &handler{client: cli}

	for i := 0; i < 100; i++ {
		closeChan := make(chan bool, 1)
		go func() {
			time.Sleep(1 * time.Millisecond)
			close(closeChan)
		}()
		_ = h.watch("/lock-test", 0, closeChan)
	}
	// Allow goroutine panics (if any) to surface
	time.Sleep(50 * time.Millisecond)
	_ = sync.Mutex{}
}
