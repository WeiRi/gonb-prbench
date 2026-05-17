package pss

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX (PR #18523): topicHandlerCaps map accesses go through
// getTopicHandlerCaps/setTopicHandlerCaps which use topicHandlerCapsMu.
func TestRace_18523_topichandlercaps(t *testing.T) {
	p := &Pss{topicHandlerCaps: make(map[Topic]*handlerCaps)}
	t1 := Topic{1}
	t2 := Topic{2}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			p.setTopicHandlerCaps(t1, &handlerCaps{})
			p.setTopicHandlerCaps(t2, &handlerCaps{})
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			_, _ = p.getTopicHandlerCaps(t1)
			_, _ = p.getTopicHandlerCaps(t2)
		}
	}()
	wg.Wait()
}
