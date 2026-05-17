package server

import (
	"sync"
	"testing"
)

// TestRace_4068_InPlace reproduces the data race in addAllServiceImportSubs
// (server/accounts.go) by calling the REAL upstream function that iterates
// a.imports.services without a.mu.RLock(), racing with concurrent writes to
// the same map.
//
// Bug: addAllServiceImportSubs does:
//   for _, si := range a.imports.services { ... }
// without a.mu.RLock(), while addServiceImport does:
//   a.mu.Lock(); a.imports.services[from] = si; a.mu.Unlock()
//
// Fix: copy services into local slice under a.mu.RLock() then iterate.
func TestRace_4068_InPlace(t *testing.T) {
	acc := NewAccount("test4068")

	// Pre-populate with service import entries
	acc.mu.Lock()
	if acc.imports.services == nil {
		acc.imports.services = make(map[string]*serviceImport)
	}
	for i := 0; i < 10; i++ {
		key := string(rune('A' + i))
		acc.imports.services[key] = &serviceImport{from: key}
	}
	acc.mu.Unlock()

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Readers: call the REAL buggy addAllServiceImportSubs
	// This iterates a.imports.services WITHOUT holding a.mu.RLock()
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				acc.addAllServiceImportSubs()
			}
		}()
	}

	// Writers: modify the map under lock (simulating addServiceImport writes)
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-ready
			key := string(rune('A' + 20 + id))
			si := &serviceImport{from: key}
			for j := 0; j < iterations; j++ {
				acc.mu.Lock()
				acc.imports.services[key] = si
				acc.mu.Unlock()
				acc.mu.Lock()
				delete(acc.imports.services, key)
				acc.mu.Unlock()
			}
		}(i)
	}

	close(ready)
	wg.Wait()
}
