// Pre-fix fixture.go from PR #103487 (client-go testing fixture).
// BUG: tracker.add / Delete pass obj directly into watch callbacks that may
// modify the same Object concurrently with subsequent reads in tests.
package fixture

import "sync"

type Object struct {
	mu   sync.Mutex
	Data map[string]string
}

// DeepCopyObject returns a deep copy.
func (o *Object) DeepCopyObject() *Object {
	o.mu.Lock()
	defer o.mu.Unlock()
	d := make(map[string]string, len(o.Data))
	for k, v := range o.Data {
		d[k] = v
	}
	return &Object{Data: d}
}

type Watcher struct{}

// Modify -- watcher writes to obj.Data without lock (simulates downstream
// listener mutating the received object). This is the racy site in the test.
// fixture.go:295/403 in pre-fix tracker.add path.
func (w *Watcher) Modify(obj *Object) {
	obj.Data["modified"] = "by-watcher" // racy WRITE on shared object map
}

// Add -- fixture.go:419
func (w *Watcher) Add(obj *Object) {
	obj.Data["added"] = "by-watcher" // racy WRITE
}

// Delete -- fixture.go:406 in pre-fix
func (w *Watcher) Delete(obj *Object) {
	obj.Data["deleted"] = "by-watcher" // racy WRITE
}

type tracker struct {
	objects map[string]*Object
	watches []*Watcher
}

func NewTracker() *tracker {
	return &tracker{objects: make(map[string]*Object), watches: []*Watcher{{}, {}}}
}

// add -- pre-fix path: w.Modify(obj) WITHOUT DeepCopy. fixture.go:403
func (t *tracker) add(name string, obj *Object) {
	if _, ok := t.objects[name]; ok {
		for _, w := range t.watches {
			w.Modify(obj) // PRE-FIX: passes obj directly (no DeepCopy) -> race
		}
		t.objects[name] = obj
		return
	}
	t.objects[name] = obj
	for _, w := range t.watches {
		w.Add(obj) // PRE-FIX: fixture.go:419
	}
}

// Delete -- fixture.go:406 pre-fix
func (t *tracker) Delete(name string) {
	if obj, ok := t.objects[name]; ok {
		delete(t.objects, name)
		for _, w := range t.watches {
			w.Delete(obj) // PRE-FIX: passes obj directly
		}
	}
}
