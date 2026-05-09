package server

import (
	"sync"
	"testing"
)

// TestRaceHandleRoutezNc reproduces the data race in HandleRoutez
// (server/monitor.go) where r.nc is read at line ~325 AFTER r.mu.Unlock()
// at line ~323, racing with closeConnection() which sets r.nc = nil under
// r.mu.Lock() at client.go line ~1315.
//
// Bug pattern:
//   HandleRoutez:     r.mu.Lock(); ... r.mu.Unlock(); _ = r.nc.(type)
//   closeConnection:  c.mu.Lock(); c.nc = nil; c.mu.Unlock()
//
// The r.nc read is outside the lock, so it races with closeConnection's
// write c.nc = nil.
func TestRaceHandleRoutezNc(t *testing.T) {
	type raceClient struct {
		mu sync.Mutex
		nc int // simulate net.Conn pointer (0=nil, 1=valid)
	}

	numReaders := 100
	numWriters := 50
	iterations := 1000
	var wg sync.WaitGroup
	ready := make(chan struct{})

	// Shared race client -- like a route entry in s.routes
	rc := &raceClient{nc: 1}

	// Readers: simulate HandleRoutez pattern
	//   r.mu.Lock(); read some fields; r.mu.Unlock();
	// BUG: then read r.nc WITHOUT re-acquiring lock
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				rc.mu.Lock()
				// Read some fields under lock (line ~303-321)
				_ = rc.nc
				rc.mu.Unlock()
				// BUG: read rc.nc AGAIN after unlock (line ~325)
				// The original code does:
				//   if ip, ok := r.nc.(*net.TCPConn); ok { ... }
				// This reads r.nc WITHOUT holding r.mu
				_ = rc.nc
			}
		}()
	}

	// Writers: simulate closeConnection setting nc = nil under lock
	//   c.mu.Lock(); c.nc = nil; c.mu.Unlock()  (client.go:1315)
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				rc.mu.Lock()
				rc.nc = 0 // simulate c.nc = nil
				rc.mu.Unlock()
				rc.mu.Lock()
				rc.nc = 1 // restore
				rc.mu.Unlock()
			}
		}()
	}

	close(ready)
	wg.Wait()
}
