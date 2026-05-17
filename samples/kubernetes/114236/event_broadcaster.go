package events

import "sync"

// Event mirrors eventsv1.Event for kubernetes-114236 stub.
type Event struct {
	Name  string
	Count int32
}

func (e *Event) DeepCopy() *Event { c := *e; return &c }

type eventKey struct{ Name string }

type eventBroadcasterImpl struct {
	mu         sync.Mutex
	eventCache map[eventKey]*Event
}

// attemptRecording mutates the event (simulating real recordEvent behavior).
func (b *eventBroadcasterImpl) attemptRecording(e *Event) {
	e.Count++ // mutation that races
}

// recordToSink — BUG: returns the same cached pointer that attemptRecording mutates.
func (b *eventBroadcasterImpl) recordToSink(ev *Event) {
	var ev2 *Event
	func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		k := eventKey{Name: ev.Name}
		if cached, ok := b.eventCache[k]; ok {
			ev2 = cached // BUG: return same pointer
			return
		}
		b.eventCache[k] = ev
		ev2 = ev // BUG: return same pointer
	}()
	if ev2 != nil {
		go b.attemptRecording(ev2) // mutates ev2 → races with other goroutines reading it
	}
}
