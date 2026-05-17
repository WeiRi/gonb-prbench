package pss

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: Pss.topicHandlerCaps map accessed without lock from Register / handlePssMsg /
// processSym. Concurrent calls race on the map. PR #18523 adds topicHandlerCapsMu
// and getTopicHandlerCaps/setTopicHandlerCaps helpers.
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
			p.topicHandlerCaps[t1] = &handlerCaps{}
			p.topicHandlerCaps[t2] = &handlerCaps{}
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = p.topicHandlerCaps[t1]
			_ = p.topicHandlerCaps[t2]
		}
	}()
	wg.Wait()
}
