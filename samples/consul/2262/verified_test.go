package agent

import (
	"io"
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: GatedWriter.Write holds RLock while appending to w.buf.
// Concurrent goroutines under RLock → concurrent slice append → race on buf header.
// FIX (PR #2262): upgrade to Lock for the append path.
func TestRace_2262_gatedwriter_write(t *testing.T) {
	gw := &GatedWriter{Writer: io.Discard}

	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			payload := []byte("payload")
			for j := 0; j < 2000 && atomic.LoadInt32(&done) == 0; j++ {
				_, _ = gw.Write(payload)
			}
		}()
	}
	wg.Wait()
	atomic.StoreInt32(&done, 1)
}
