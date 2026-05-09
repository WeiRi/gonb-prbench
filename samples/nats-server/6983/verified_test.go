package server

import (
	"sync"
	"testing"
)

// TestRaceCopyStreamMetadata reproduces the data race in
// addStreamWithAssignment (server/stream.go) where copyStreamMetadata
// writes to a shared cfg struct without holding js.mu lock, racing
// with reflect.DeepEqual which reads the same struct.
//
// Bug: addStreamWithAssignment does:
//   copyStreamMetadata(cfg, &ocfg)         // WRITE to ocfg - no lock
//   reflect.DeepEqual(cfg, &ocfg)          // READ from ocfg - no lock
// While another goroutine could be simultaneously reading/writing ocfg.
//
// Fix: wrap both copyStreamMetadata and reflect.DeepEqual under js.mu.Lock().
func TestRaceCopyStreamMetadata(t *testing.T) {
	// Simulate shared metadata struct that's copied to and compared
	type metadata struct {
		fields [16]int
	}
	shared := &metadata{}

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Readers: simulate reflect.DeepEqual reading the shared struct
	// without holding the lock
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// BUG: reading shared fields without lock
				// Simulates reflect.DeepEqual(cfg, &ocfg)
				for k := 0; k < len(shared.fields); k++ {
					_ = shared.fields[k]
				}
			}
		}()
	}

	// Writers: simulate copyStreamMetadata writing to the shared struct
	// without holding the lock
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// BUG: writing to shared fields without lock
				// Simulates copyStreamMetadata(cfg, &ocfg)
				for k := 0; k < len(shared.fields); k++ {
					shared.fields[k] = j + k
				}
			}
		}()
	}

	close(ready)
	wg.Wait()
}
