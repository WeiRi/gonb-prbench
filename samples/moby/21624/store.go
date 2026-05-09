// Production stub modeling moby PR #21624 — volume/store/store.go data race
// on s.names map (and globalLock not held during create).
package main

import "sync"

type VolumeStore struct {
	globalLock sync.Mutex
	names      map[string]string
	labels     map[string]map[string]string
}

func NewVolumeStore() *VolumeStore {
	return &VolumeStore{
		names:  make(map[string]string),
		labels: make(map[string]map[string]string),
	}
}

// create writes into s.labels and s.names without holding globalLock — racy.
func (s *VolumeStore) create(name string, lbl map[string]string) { // RACE write site
	s.labels[name] = lbl
	s.names[name] = name
}

// getVolume reads s.names without holding globalLock — racy.
func (s *VolumeStore) getVolume(name string) string { // RACE read site
	return s.names[name]
}
