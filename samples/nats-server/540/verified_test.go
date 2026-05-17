package server

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BUG: HandleRoutez releases r.mu before reading r.nc, while another goroutine
// may mutate r.nc under r.mu. Race on r.nc field.
// PR #540 moves r.mu.Unlock() AFTER the r.nc switch case.
func TestRace_540_routez_nc(t *testing.T) {
	s := &Server{
		routes:       map[uint64]*client{},
		httpReqStats: map[string]uint64{},
	}
	c := &client{
		cid:   1,
		start: time.Now(),
		route: &route{},
		subs:  map[string]*subscription{},
	}
	s.routes[1] = c

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		req := httptest.NewRequest(http.MethodGet, "/routez", nil)
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			rec := httptest.NewRecorder()
			s.HandleRoutez(rec, req)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			c.mu.Lock()
			c.nc = nil
			c.mu.Unlock()
		}
	}()
	wg.Wait()
}
