package events

import (
	"sync"
	"testing"
)

func TestRace_114236_EventBroadcasterSharedCache(t *testing.T) {
	const N = 200
	for n := 0; n < N; n++ {
		bc := NewBroadcaster()
		// Pre-populate cache.
		bc.eventCache["k"] = &Event{Reason: "init"}

		var wg sync.WaitGroup
		wg.Add(2)
		// goroutine A: records (mutates returned event)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				ev := bc.recordToSink("k", &Event{Reason: "x"})
				attemptRecording(ev)
			}
		}()
		// goroutine B: records concurrently — same eventKey -> same cached object
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				ev := bc.recordToSink("k", &Event{Reason: "y"})
				attemptRecording(ev)
			}
		}()
		wg.Wait()
	}
}
