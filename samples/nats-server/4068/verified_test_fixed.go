package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_nats_4068_services_map(t *testing.T) {
	a := &Account{Name: "a", imports: importMap{services: map[string]*serviceImport{}}}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			a.mu.Lock()
			a.imports.services["k"] = &serviceImport{}
			a.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			var sis [32]*serviceImport
			serviceImports := sis[:0]
			a.mu.RLock()
			for _, si := range a.imports.services {
				serviceImports = append(serviceImports, si)
			}
			a.mu.RUnlock()
			for _, si := range serviceImports {
				_ = si
			}
		}
	}()
	wg.Wait()
}
