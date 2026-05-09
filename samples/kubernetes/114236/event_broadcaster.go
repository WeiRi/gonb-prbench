// Pre-fix event_broadcaster.go from PR #114236.
// BUG: recordToSink returns the cached *Event directly (not DeepCopy), so the
// client and the broadcaster's cache share the same object — a downstream
// recorder mutates fields concurrently with another goroutine's read.
package events

import "sync"

type Event struct {
	mu     sync.Mutex
	Series *EventSeries
	Reason string
	Count  int
}

type EventSeries struct {
	Count            int
	LastObservedTime int64
}

// DeepCopy is the post-fix safe copy.
func (e *Event) DeepCopy() *Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	c := &Event{Reason: e.Reason, Count: e.Count}
	if e.Series != nil {
		c.Series = &EventSeries{Count: e.Series.Count, LastObservedTime: e.Series.LastObservedTime}
	}
	return c
}

type eventBroadcasterImpl struct {
	mu         sync.Mutex
	eventCache map[string]*Event
}

func NewBroadcaster() *eventBroadcasterImpl {
	return &eventBroadcasterImpl{eventCache: map[string]*Event{}}
}

// recordToSink (PRE-FIX): event_broadcaster.go:187 returns isomorphicEvent
// directly without DeepCopy => caller can mutate the cached object.
func (e *eventBroadcasterImpl) recordToSink(eventKey string, ev *Event) *Event {
	e.mu.Lock()
	if cached, ok := e.eventCache[eventKey]; ok {
		// pre-fix: mutate cached.Series, then return cached object (not deep-copied).
		cached.Series = &EventSeries{Count: cached.Count + 1, LastObservedTime: 1}
		evToRecord := cached
		e.mu.Unlock()
		return evToRecord
	}
	e.eventCache[eventKey] = ev
	e.mu.Unlock()
	return ev
}

// attemptRecording -- caller side; mutates fields without lock.
func attemptRecording(ev *Event) {
	ev.Reason = "recorded"
	ev.Count++
}
