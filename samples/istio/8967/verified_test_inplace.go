package fs

import (
	"sync"
	"testing"
)

// TestRace_8967_InPlace reproduces the data race in fsSource.Stop()
// (galley/pkg/fs/fssource.go) where Stop() sets fileResorceKeys, shas,
// and donec to nil without synchronization, racing with other goroutines
// reading those same fields.
//
// Bug: Stop() does:
//   s.fileResorceKeys = nil
//   s.shas = nil
//   s.donec = nil
// without holding any lock, while other goroutines read these fields.
//
// Fix: remove the nil assignments (closing donec is sufficient).
func TestRace_8967_InPlace(t *testing.T) {
	s := &fsSource{
		fileResorceKeys: make(map[string][]*fileResourceKey),
		shas:            make(map[string][20]byte),
		donec:           make(chan struct{}),
	}

	numReaders := 100
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Writers: call Stop() which sets fields to nil without lock
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// Create a fresh fsSource each time since Stop() destroys it
				src := &fsSource{
					fileResorceKeys: make(map[string][]*fileResourceKey),
					shas:            make(map[string][20]byte),
					donec:           make(chan struct{}),
				}
				func() {
					defer func() { recover() }()
					src.Stop()
				}()
				// Also trigger race on the shared object
				func() {
					defer func() { recover() }()
					s.Stop()
				}()
			}
		}()
	}

	// Readers: read the shared fields while Stop() writes nil
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				_ = s.fileResorceKeys
				_ = s.shas
				_ = s.donec
			}
		}()
	}

	close(ready)
	wg.Wait()
}
