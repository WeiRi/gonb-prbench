// In-place race test for nomad-14188: package=stream, uses upstream EventBroker.
// Bug: event_broker.go:25/104 -- Publish sends &event (pointer to range-loop variable)
// on aclCh, so receiver sees data from next iteration. Race between range-loop write
// and channel receiver read on loop variable memory.
// PR fix: change channel from chan *Event to chan Event (value type).
// NOTE: Upstream event_broker.go was modified to reintroduce the buggy pointer channel type.
package stream

import (
	"sync"
	"testing"

	"github.com/hashicorp/nomad/nomad/structs"
)

func TestRace_14188_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	var wg sync.WaitGroup

	broker := &EventBroker{
		aclCh:     make(chan *structs.Event, N*ITERS*3),
		publishCh: make(chan *structs.Events, N*ITERS),
	}

	// Receiver goroutines: read from aclCh (read *Event fields)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				evt := <-broker.aclCh
				_ = evt.Topic // RACE READ on loop-variable memory
				_ = evt.Type
				_ = evt.Key
				_ = evt.Index
			}
		}()
	}

	// Publish events: writes &event pointer, overwriting same memory (event_broker.go:104)
	for iter := 0; iter < N*ITERS; iter++ {
		events := &structs.Events{
			Events: []structs.Event{
				{Topic: "ACLToken", Type: "updated", Key: "token-1", Index: uint64(iter)},
				{Topic: "ACLPolicy", Type: "updated", Key: "policy-1", Index: uint64(iter)},
			},
		}
		broker.Publish(events)
	}

	wg.Wait()
}
