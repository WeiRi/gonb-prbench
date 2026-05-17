package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_nats_6620_sys_access(t *testing.T) {
	acc1 := &Account{Name: "a1"}
	acc2 := &Account{Name: "a2"}
	s := &Server{sys: &internal{account: acc1}}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			s.mu.Lock()
			if j%2 == 0 {
				s.sys = &internal{account: acc1}
			} else {
				s.sys = &internal{account: acc2}
			}
			s.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			s.mu.Lock()
			s.mu.Unlock()
			if s.sys != nil {
				_ = s.sys.account
			}
		}
	}()
	wg.Wait()
}
