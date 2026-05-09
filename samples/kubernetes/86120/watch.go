package watch

import "sync"

// FakeWatcher is a stripped reproduction of staging/src/k8s.io/apimachinery/pkg/watch/watch.go.
// BUG (pre-PR #86120): Stopped is read directly while Stop()/Reset() write it under a mutex.
type FakeWatcher struct {
	Stopped bool         // racy field
	mu      sync.Mutex   // guards writes only — reads are unprotected (BUG)
}

func NewFake() *FakeWatcher {
	return &FakeWatcher{}
}

func (f *FakeWatcher) Stop() {
	f.mu.Lock()
	f.Stopped = true
	f.mu.Unlock()
}

func (f *FakeWatcher) Reset() {
	f.mu.Lock()
	f.Stopped = false
	f.mu.Unlock()
}
