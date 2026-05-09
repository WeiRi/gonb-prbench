package stream

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/nomad/nomad/structs"
)

// TestRaceACLChannelPointerBug triggers the data race where Publish sends
// a pointer to the range-loop variable (&event) on aclCh, causing the
// receiver to see overwritten data from subsequent loop iterations.
// Fix: change channel from chan *structs.Event to chan structs.Event (value type).
func TestRaceACLChannelPointerBug(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker, err := NewEventBroker(ctx, nil, EventBrokerCfg{
		EventBufferSize: 100,
		Logger:          nil,
	})
	if err != nil {
		t.Fatalf("failed to create event broker: %v", err)
	}

	var wg sync.WaitGroup

	// Receiver goroutine: reads from aclCh (bug: receives *Event pointers
	// that point to the range-loop variable which gets overwritten)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			select {
			case evt := <-broker.aclCh:
				// RACE: reading *evt fields while Publish loop
				// modifies the underlying event variable
				_ = evt.Topic
				_ = evt.Type
				_ = evt.Key
				_ = evt.Index
			case <-ctx.Done():
				return
			}
		}
	}()

	// Publish ACL events in a loop. The range-loop variable &event
	// pointer is reused, causing the data race.
	for iter := 0; iter < 500; iter++ {
		events := &structs.Events{
			Index: uint64(iter),
			Events: []structs.Event{
				{Topic: structs.TopicACLToken, Type: "acl-token-updated", Key: "token-1"},
				{Topic: structs.TopicACLPolicy, Type: "acl-policy-updated", Key: "policy-1"},
				{Topic: structs.TopicACLToken, Type: "acl-token-deleted", Key: "token-2"},
			},
		}
		broker.Publish(events)
	}

	wg.Wait()
}
