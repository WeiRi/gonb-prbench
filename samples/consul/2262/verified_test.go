package agent

import (
	"sync"
	"testing"
)

// TestGatedWriterRace reproduces the data race in GatedWriter.Write() where
// w.buf = append(w.buf, p2) is done under RLock() instead of Lock().
// PR 2262: upgrade RWMutex RLock to Lock in GatedWriter.Write() for safe buffer append.
func TestGatedWriterRace(t *testing.T) {
	// Use a discard-like writer; in the unbuffered path (flush=true) it would
	// be used, but since we never flush, all writes go to buf.
	var nop nopWriter
	gw := &GatedWriter{
		Writer: &nop,
	}

	data := []byte("hello")

	var wg sync.WaitGroup
	numGoroutines := 60
	iterations := 300

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				gw.Write(data)
			}
		}()
	}

	wg.Wait()

	// After all writes, buf should have numGoroutines * iterations entries
	expected := numGoroutines * iterations
	if len(gw.buf) != expected {
		t.Logf("buf length: %d, expected: %d (race may cause loss)", len(gw.buf), expected)
	}
}

// nopWriter implements io.Writer (discards writes)
type nopWriter struct{}

func (n *nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
