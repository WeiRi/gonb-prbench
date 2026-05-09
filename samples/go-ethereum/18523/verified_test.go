package pss

import (
	"sync"
	"testing"
)

// TestRace_18523_topicHandlerCaps: Register writes topicHandlerCaps map
// without lock while HandlePssMsg reads same map — concurrent map
// read+write fires under -race.
func TestRace_18523_topicHandlerCaps(t *testing.T) {
	p := NewPss()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			p.Register(Topic(i%16), i%2 == 0)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = p.HandlePssMsg(Topic(i % 16))
		}
	}()
	wg.Wait()
}
