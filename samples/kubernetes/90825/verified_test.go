// VERIFIED whitebox PoC for kubernetes-90825 (rebuilt, ORDER oracle).
package pock90825

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestRace_kubernetes90825 exercises the "broadcast not under f.lock" bug.
// Multiple Pops park on cond.Wait; Close() then runs without holding f.lock,
// so a waiter that has just returned from Wait may re-enter Wait before
// Close's Broadcast hits — and no further broadcast ever comes.
func TestRace_kubernetes90825(t *testing.T) {
	const iters = 30
	const jobs = 200

	for it := 0; it < iters; it++ {
		f := NewFIFO()
		var done int64
		var wg sync.WaitGroup
		wg.Add(jobs)
		for i := 0; i < jobs; i++ {
			go func() {
				defer wg.Done()
				_, _ = f.Pop(func(string) error { return nil })
				atomic.AddInt64(&done, 1)
			}()
		}

		// Tickle: bursty add/close to widen race window
		go func() {
			for k := 0; k < 8; k++ {
				f.Add("item")
				runtime.Gosched()
			}
			f.Close()
		}()

		doneCh := make(chan struct{})
		go func() { wg.Wait(); close(doneCh) }()
		select {
		case <-doneCh:
		case <-time.After(2 * time.Second):
			buf := make([]byte, 1<<16)
			n := runtime.Stack(buf, true)
			t.Logf("iter %d goroutine dump:\n%s", it, string(buf[:n]))
			t.Fatalf("iter %d: ORDER oracle: %d/%d Pops still blocked %v after Close", it, jobs-int(atomic.LoadInt64(&done)), jobs, 2*time.Second)
		}
	}
}
