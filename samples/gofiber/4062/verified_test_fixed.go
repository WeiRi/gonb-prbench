package fiber

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: Context() handles nil c.fasthttp by returning context.Background().
// No panic.
func TestRace_gofiber_4062_nil_fasthttp(t *testing.T) {
	c := &DefaultCtx{}

	var done int32
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100 && atomic.LoadInt32(&done) == 0; j++ {
			_ = c.Context()
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
