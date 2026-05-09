package cache

import "sync"

// Stripped reproduction of staging/src/k8s.io/client-go/tools/cache/delta_fifo.go pre-PR #124344.
// BUG: Replace() invokes the user-supplied transformer that mutates objects in-place,
// while GetByKey returns the same object pointer to readers without coordination.

type Delta struct {
	Type   string
	Object interface{}
}

type Deltas []Delta

type DeltaFIFOOptions struct {
	KeyFunction func(obj interface{}) (string, error)
	Transformer func(obj interface{}) (interface{}, error)
}

type DeltaFIFO struct {
	mu          sync.Mutex
	keyFunc     func(obj interface{}) (string, error)
	transformer func(obj interface{}) (interface{}, error)
	items       map[string]Deltas
	queue       []string
}

func NewDeltaFIFOWithOptions(opts DeltaFIFOOptions) *DeltaFIFO {
	return &DeltaFIFO{
		keyFunc:     opts.KeyFunction,
		transformer: opts.Transformer,
		items:       map[string]Deltas{},
	}
}

// Replace — BUG: calls transformer(obj) which can mutate obj in-place,
// while another goroutine is reading the same obj via GetByKey.
func (f *DeltaFIFO) Replace(list []interface{}, _ string) error {
	for _, obj := range list {
		// BUG: transformer mutates obj in-place WITHOUT holding f.mu
		// (it's invoked while iterating; lock is held but readers also race).
		o := obj
		if f.transformer != nil {
			o, _ = f.transformer(obj)        // line 41 — racing write to obj fields
		}
		k, _ := f.keyFunc(o)
		f.mu.Lock()
		if _, ok := f.items[k]; !ok {
			f.queue = append(f.queue, k)
		}
		f.items[k] = Deltas{{Type: "Sync", Object: o}}
		f.mu.Unlock()
	}
	return nil
}

// GetByKey — BUG: returns the stored object pointer; readers access
// fields without coordination with Replace's transformer mutation.
func (f *DeltaFIFO) GetByKey(key string) (interface{}, bool, error) {
	f.mu.Lock()
	d, ok := f.items[key]
	f.mu.Unlock()
	if !ok {
		return nil, false, nil
	}
	return d, true, nil
}

// Pop drains one item.
func (f *DeltaFIFO) Pop(process func(obj interface{}, isInInitialList bool) error) (interface{}, error) {
	f.mu.Lock()
	if len(f.queue) == 0 {
		f.mu.Unlock()
		return nil, nil
	}
	k := f.queue[0]
	f.queue = f.queue[1:]
	d := f.items[k]
	delete(f.items, k)
	f.mu.Unlock()
	if err := process(d, false); err != nil {
		return nil, err
	}
	return d, nil
}
