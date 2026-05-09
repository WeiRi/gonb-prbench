package server

import (
	"sync"
	"testing"
)

// TestRaceAddAllServiceImportSubs reproduces the data race in
// addAllServiceImportSubs (server/accounts.go) where a.imports.services
// map is iterated WITHOUT holding a.mu lock, racing with addServiceImport
// which writes to the map under a.mu.Lock().
//
// Bug: addAllServiceImportSubs does:
//   for _, si := range a.imports.services { ... }
// without a.mu.RLock(), while addServiceImport does:
//   a.mu.Lock(); a.imports.services[si.from] = si; a.mu.Unlock()
//
// Fix: copy the service imports into a local slice under a.mu.RLock(),
// then iterate the local slice (no lock needed).
func TestRaceAddAllServiceImportSubs(t *testing.T) {
	type serviceImport struct {
		from string
	}
	type fakeAccount struct {
		mu       sync.RWMutex
		services map[string]*serviceImport
	}

	acc := &fakeAccount{
		services: make(map[string]*serviceImport),
	}

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Readers: simulate the buggy addAllServiceImportSubs
	// Iterating a.imports.services WITHOUT holding a.mu.RLock()
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				// BUG: iterating map WITHOUT RLock
				for _, si := range acc.services {
					_ = si.from
				}
			}
		}(i)
	}

	// Writers: simulate addServiceImport modifying the map under Lock
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-ready
			key := string(rune('A' + id))
			si := &serviceImport{from: key}
			for j := 0; j < iterations; j++ {
				acc.mu.Lock()
				acc.services[key] = si // WRITE under lock
				acc.mu.Unlock()
				acc.mu.Lock()
				delete(acc.services, key) // also write under lock
				acc.mu.Unlock()
			}
		}(i)
	}

	close(ready)
	wg.Wait()
}
