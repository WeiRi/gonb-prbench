package events

import (
	"sync"
	"testing"
	"time"
)

func TestRace_114236_EventBroadcasterShared(t *testing.T) {
	bc := &eventBroadcasterImpl{eventCache: map[eventKey]*Event{}}
	ev := &Event{Name: "ev1"}

	const N = 100
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			bc.recordToSink(ev)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			bc.recordToSink(ev)
		}
	}()
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
}
