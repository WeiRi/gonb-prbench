package logger

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: RingLogger.Close() drains buffer and calls r.l.Log for each msg, then calls
// r.l.Close, but doesn't wait for run() goroutine to exit. run() may still be inside
// r.l.Log when Close's drain loop also calls r.l.Log → race on the wrapped Logger.
// PR #42320 adds wg.WaitGroup so Close waits for run() before drain.
type rgFakeLog42320 struct {
	counter int64
}

func (f *rgFakeLog42320) Log(m *Message) error  { f.counter++; return nil }
func (f *rgFakeLog42320) Close() error          { return nil }
func (f *rgFakeLog42320) Name() string          { return "fake-42320" }

func TestRace_42320_run_vs_close_drain(t *testing.T) {
	fl := &rgFakeLog42320{}
	rl := newRingLogger(fl, Info{ContainerID: "x"}, 1024*1024)

	var done int32
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		msg := &Message{Line: []byte("hello")}
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = rl.Log(msg)
		}
	}()

	// give producer a head start
	for i := 0; i < 100; i++ {
		_ = rl.Log(&Message{Line: []byte("seed")})
	}

	_ = rl.Close()
	atomic.StoreInt32(&done, 1)
	wg.Wait()
	_ = fl.counter
}
