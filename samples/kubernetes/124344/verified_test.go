package cache

import (
	"sync"
	"testing"
)

// simpleObj is a minimal type that DeltaFIFO can store and that the
// transformer can mutate in-place to trigger the data race.
type simpleObj struct {
	ID    string
	Value int
}

// TestRace_PR124344_TransformerResync triggers the data race between:
// - Replace()/Resync() calling the TransformFunc that mutates objects in-place, and
// - GetByKey() returning the same object pointer to callers outside the DeltaFIFO lock.
//
// The race: DeltaFIFO.Replace() calls f.transformer(obj) which can mutate 'obj'
// in place (line 468 in delta_fifo.go). GetByKey() returns *Deltas containing
// the same object pointer. When two goroutines concurrently run
// Replace (mutating via transformer) and GetByKey (reading the object), the
// Go race detector catches the unsynchronized access.
//
// Fix: The PR adds transformer idempotency guarantees and coordination with
// readers, ensuring objects already in cache are not re-transformed during
// Replace/Resync for default informer usage.
func TestRace_PR124344_TransformerResync(t *testing.T) {
	const numObjects = 100
	const iterations = 500

	// Create objects that will be mutated by the transformer.
	objs := make([]*simpleObj, numObjects)
	for i := 0; i < numObjects; i++ {
		objs[i] = &simpleObj{ID: string(rune('A' + i%26)), Value: i}
	}

	// Key function: use the object itself as key.
	keyFunc := func(obj interface{}) (string, error) {
		return obj.(*simpleObj).ID, nil
	}

	// Transformer: mutates object in-place (the root cause of the bug).
	// It increments Value, simulating a cleanup/annotation operation.
	transformer := func(obj interface{}) (interface{}, error) {
		o := obj.(*simpleObj)
		o.Value++ // mutate in-place
		return o, nil
	}

	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: keyFunc,
		Transformer: transformer,
	})

	// Initial population.
	ifaceList := make([]interface{}, numObjects)
	for i, o := range objs {
		ifaceList[i] = o
	}
	if err := f.Replace(ifaceList, ""); err != nil {
		t.Fatalf("Initial Replace failed: %v", err)
	}

	var wg sync.WaitGroup

	// Goroutine group 1: Replace() — triggers transformer that mutates objects in-place.
	for g := 0; g < 20; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				list := make([]interface{}, numObjects)
				for j, o := range objs {
					list[j] = o
				}
				f.Replace(list, "")
			}
		}()
	}

	// Goroutine group 2: GetByKey() — reads objects returned from the queue.
	// The returned Deltas contain object pointers that the transformer
	// concurrently mutates.
	for g := 0; g < 20; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				for _, o := range objs {
					item, exists, _ := f.GetByKey(o.ID)
					if exists {
						// Access the object data — races with transformer mutation.
						if deltas, ok := item.(Deltas); ok && len(deltas) > 0 {
							_ = deltas[0].Object.(*simpleObj).Value
						}
					}
				}
			}
		}()
	}

	// Goroutine group 3: Pop() — empties the queue so Replace() can
	// re-enqueue and re-run the transformer.
	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations/10; i++ {
				f.Pop(func(obj interface{}, isInInitialList bool) error {
					return nil
				})
			}
		}()
	}

	wg.Wait()
}
