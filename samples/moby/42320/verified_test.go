// Race-trigger test for moby-42320; see README.md for usage.

package pocmoby42320

import (
	"sync"
	"sync/atomic"
)

// BUG version of moby daemon/logger/ring.go (pre-PR-42320):
// Close() returns before run() goroutine exits.  Caller may then close the
// underlying logger's channel; run() still pushes to it -> "send on closed
// channel" panic.  PR adds sync.WaitGroup to make Close() wait for run().

type Message struct{ Line []byte }

// underlying is a stand-in for fluentd-style downstream logger that
// will panic if Log() is called after the user closes its input channel.
type underlying struct {
	ch     chan Message
	closed int32
}

func newUnderlying() *underlying { return &underlying{ch: make(chan Message, 64)} }

// Log mimics fluentd-go-logger: sends to internal channel; PANICs if channel closed.
func (u *underlying) Log(m Message) error {
	u.ch <- m // <- panic: send on closed channel
	return nil
}

// CloseChannel closes the input chan (called by orchestrator after RingLogger.Close()
// has returned, on the assumption that no more Log() calls will come).
func (u *underlying) CloseChannel() {
	if atomic.CompareAndSwapInt32(&u.closed, 0, 1) {
		close(u.ch)
	}
}

func (u *underlying) drain() {
	for range u.ch {
	}
}

type RingLogger struct {
	buffer    chan Message
	l         *underlying
	closeFlag int32

	runPanic     *int32
	runStackOnce *sync.Once
	runStackPtr  *string
}

func newRingLogger(driver *underlying, runPanic *int32, once *sync.Once, stackPtr *string) *RingLogger {
	l := &RingLogger{
		buffer:       make(chan Message, 16),
		l:            driver,
		runPanic:     runPanic,
		runStackOnce: once,
		runStackPtr:  stackPtr,
	}
	go l.run()
	return l
}

func (r *RingLogger) Log(m Message) error {
	if atomic.LoadInt32(&r.closeFlag) == 1 {
		return nil
	}
	defer func() { recover() }()
	r.buffer <- m
	return nil
}

func (r *RingLogger) closed() bool { return atomic.LoadInt32(&r.closeFlag) == 1 }

func (r *RingLogger) Close() error {
	atomic.StoreInt32(&r.closeFlag, 1)
	defer func() { recover() }()
	close(r.buffer)
	// BUG: does NOT wait for run() to finish. Returns immediately.
	return nil
}

// run drains buffer to the underlying driver.
func (r *RingLogger) run() {
	defer func() {
		if rec := recover(); rec != nil {
			atomic.StoreInt32(r.runPanic, 1)
			r.runStackOnce.Do(func() {
				buf := make([]byte, 4096)
				n := runtimeStack(buf)
				*r.runStackPtr = string(buf[:n])
			})
		}
	}()
	for m := range r.buffer {
		_ = r.l.Log(m) // may panic if caller closes underlying after Close()
	}
}


import "runtime"

func runtimeStack(buf []byte) int { return runtime.Stack(buf, false) }


import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Reproduces moby PR 42320 lifecycle-ordering panic:
// RingLogger.Close() returns while run() goroutine is still pumping into the
// underlying driver. Caller, observing Close() returned, closes the
// underlying driver's channel - run() then panics on send-on-closed-channel.
//
// Oracle: PANIC (caught inside run() goroutine). Frame target: ring.go.
func TestRace_moby42320(t *testing.T) {
	const iters = 200
	var runPanic int32
	var stack string
	var once sync.Once

	for it := 0; it < iters; it++ {
		u := newUnderlying()
		go u.drain()
		r := newRingLogger(u, &runPanic, &once, &stack)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				_ = r.Log(Message{Line: []byte("x")})
			}
		}()
		go func() {
			_ = r.Close()
			u.CloseChannel()
		}()

		wg.Wait()
		// Brief settle for run() to react.
		time.Sleep(2 * time.Millisecond)
		if atomic.LoadInt32(&runPanic) == 1 {
			t.Logf("iter %d: run() panic stack:\n%s", it, stack)
			t.Fatal("RingLogger.run() panic via underlying.Log -> send on closed channel reproduced")
		}
	}
}

