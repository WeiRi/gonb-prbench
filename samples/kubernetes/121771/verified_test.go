// Race-trigger test for kubernetes-121771; see README.md for usage.

package runtime

import (
	"bytes"
	"sync"
	"testing"
)

type fakeObj struct {
	gvk GroupVersionKind
}

func (o *fakeObj) GetGroupVersionKind() GroupVersionKind  { return o.gvk }
func (o *fakeObj) SetGroupVersionKind(g GroupVersionKind) { o.gvk = g }

// TestWithVersionEncoderRace: same obj already has target GVK; multiple
// goroutines encode it in parallel. BUG state still touches GVK -> race;
// FIX state short-circuits.
func TestWithVersionEncoderRace(t *testing.T) {
	target := GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	obj := &fakeObj{gvk: target}
	enc := WithVersionEncoder{GroupVersion: target, Encoder: fakeEncoder{}}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var buf bytes.Buffer
			for j := 0; j < 200; j++ {
				_ = enc.Encode(obj, &buf)
			}
		}()
	}
	wg.Wait()
}
