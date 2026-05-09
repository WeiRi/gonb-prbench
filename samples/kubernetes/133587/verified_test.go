package experimental

import (
	"context"
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	resourceapi "k8s.io/api/resource/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/dynamic-resource-allocation/structured/internal"
)

type mockDeviceClassLister struct{}

func (m *mockDeviceClassLister) List() ([]*resourceapi.DeviceClass, error) {
	return nil, nil
}

func (m *mockDeviceClassLister) Get(className string) (*resourceapi.DeviceClass, error) {
	return nil, nil
}

func TestRace_133587(t *testing.T) {
	alloc, err := NewAllocator(
		context.Background(),
		SupportedFeatures,
		AllocatedState{
			AllocatedDevices:         sets.New[internal.DeviceID](),
			AllocatedSharedDeviceIDs: sets.New[internal.SharedDeviceID](),
			AggregatedCapacity:       internal.NewConsumedCapacityCollection(),
		},
		&mockDeviceClassLister{},
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: Allocate writes alloc.claimsToAllocate (which is
				// a.claimsToAllocate through the embedded *Allocator) at line 149
				// without synchronization. Multiple concurrent calls race.
				node := &v1.Node{}
				node.Name = "test-node"
				alloc.Allocate(context.Background(), node, nil)
			}
		}()
	}

	wg.Wait()
}
