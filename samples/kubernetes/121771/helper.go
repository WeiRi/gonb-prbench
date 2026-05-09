package runtime

import "io"

// Stripped reproduction of staging/src/k8s.io/apimachinery/pkg/runtime/helper.go pre-PR #121771.
// BUG: WithVersionEncoder.Encode unconditionally Set/Encode/Set even when gvk == oldGVK.
// Concurrent encoders of the same object then race on obj.gvk.

type GroupVersionKind struct {
	Group, Version, Kind string
}

type Object interface {
	GetGroupVersionKind() GroupVersionKind
	SetGroupVersionKind(GroupVersionKind)
}

type Encoder interface {
	Encode(obj Object, w io.Writer) error
}

type fakeEncoder struct{}

func (fakeEncoder) Encode(obj Object, w io.Writer) error {
	_ = obj.GetGroupVersionKind()        // reads obj.gvk
	return nil
}

type WithVersionEncoder struct {
	GroupVersion GroupVersionKind
	Encoder      Encoder
}

// BUG body: always writes obj.gvk, even when oldGVK == GroupVersion.
func (e WithVersionEncoder) Encode(obj Object, stream io.Writer) error {
	gvk := e.GroupVersion
	oldGVK := obj.GetGroupVersionKind()  // line 35 (read)
	obj.SetGroupVersionKind(gvk)         // line 36 (write — RACE)
	err := e.Encoder.Encode(obj, stream) // line 37
	obj.SetGroupVersionKind(oldGVK)      // line 38 (write — RACE)
	return err
}
