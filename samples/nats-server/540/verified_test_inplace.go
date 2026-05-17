package server

import (
	"net/http/httptest"
	"sync"
	"testing"
)

// TestRace_540_InPlace reproduces the data race in HandleRoutez
// (server/monitor.go) where r.nc is read AFTER r.mu.Unlock(),
// racing with closeConnection() which sets c.nc = nil under c.mu.Lock().
//
// Bug: HandleRoutez does:
//   r.mu.Unlock()
//   if ip, ok := r.nc.(*net.TCPConn); ok { ... }
// The r.nc read at line ~325 is outside the lock.
//
// Fix: move r.mu.Unlock() after the r.nc read.
func TestRace_540_InPlace(t *testing.T) {
	s := &Server{
		httpReqStats: make(map[string]uint64),
	}

	r := &client{
		cid:   1,
		route: &route{remoteID: "test", routeType: Explicit},
		subs:  make(map[string]*subscription),
		nc:    nil,
	}
	s.routes = map[uint64]*client{1: r}

	req := httptest.NewRequest("GET", "/routez", nil)
	w := httptest.NewRecorder()

	numReaders := 80
	numWriters := 20
	iterations := 500
	var wg sync.WaitGroup
	ready := make(chan struct{})

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				func() {
					defer func() { recover() }()
					s.HandleRoutez(w, req)
				}()
			}
		}()
	}

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				r.mu.Lock()
				r.nc = nil
				r.mu.Unlock()
			}
		}()
	}

	close(ready)
	wg.Wait()
}
