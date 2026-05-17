package fiber

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: DefaultCtx.Context() dereferences c.fasthttp without nil check.
// Calling Context() after Release (c.fasthttp = nil) panics with
// nil-pointer-dereference.
// Multiple goroutines hitting this path see panic stacks.
func TestRace_gofiber_4062_nil_fasthttp(t *testing.T) {
	c := &DefaultCtx{} // fasthttp is nil

	var done int32
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				// expected: nil pointer deref panic — re-raise so test fails
				panic(r)
			}
		}()
		for j := 0; j < 1 && atomic.LoadInt32(&done) == 0; j++ {
			_ = c.Context()
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
