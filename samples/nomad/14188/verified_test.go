// Race test for nomad-14188 — loop variable pointer race in EventBroker.Publish
// fix changes aclCh from chan *structs.Event to chan structs.Event
// BUG: Publish loops `for _, event := range ...` and sends &event (pointer to reused loop var)
// concurrent receivers race on event fields
package stream

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/nomad/nomad/structs"
)

func TestRace_14188_PublishLoopVar(t *testing.T) {
	broker, err := NewEventBroker(context.Background(), nil, EventBrokerCfg{})
	if err != nil {
		t.Fatalf("NewEventBroker: %v", err)
	}

	// Drain aclCh to prevent buffer fill blocking
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			case e := <-broker.aclCh:
				_ = e
			}
		}
	}()

	// Build event sets that always trigger the ACL branch
	makeEvents := func() *structs.Events {
		return &structs.Events{
			Events: []structs.Event{
				{Topic: structs.TopicACLToken, Type: "x"},
				{Topic: structs.TopicACLPolicy, Type: "y"},
				{Topic: structs.TopicACLToken, Type: "z"},
			},
		}
	}

	var wg sync.WaitGroup
	const N = 30
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				broker.Publish(makeEvents())
			}
		}()
	}
	wg.Wait()
}
