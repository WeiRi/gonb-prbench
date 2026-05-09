package testrace

import (
	"testing"
)

// Minimal reproduction of istio/istio#8967: Stop() sets shared fields to nil
// without synchronization while other goroutines read them.
// BUG: Stop() sets fileResourceKeys, shas, donec to nil without lock.
// FIX: remove nil assignments (closing donec is sufficient).

type fsSource struct {
	fileResourceKeys []string
	shas             map[string]string
	donec            chan struct{}
}

func (s *fsSource) readKeys() []string {
	return s.fileResourceKeys // READ concurrent with Stop() WRITE
}

func (s *fsSource) readShas() map[string]string {
	return s.shas // READ concurrent with Stop() WRITE
}

func (s *fsSource) StopBuggy() {
	s.fileResourceKeys = nil // WRITE
	s.shas = nil             // WRITE
	close(s.donec)
	s.donec = nil // WRITE
}

func TestRace(t *testing.T) {
	done := make(chan struct{}, 200)

	for iter := 0; iter < 200; iter++ {
		s := &fsSource{
			fileResourceKeys: []string{"key1", "key2", "key3"},
			shas:             map[string]string{"k1": "v1", "k2": "v2"},
			donec:            make(chan struct{}),
		}

		// Reader: reads the SAME object
		go func(src *fsSource) {
			for j := 0; j < 100; j++ {
				_ = src.readKeys()
				_ = src.readShas()
				_ = src.donec
			}
			done <- struct{}{}
		}(s)

		// Writer: calls Stop on the SAME object
		go func(src *fsSource) {
			src.StopBuggy()
			done <- struct{}{}
		}(s)
	}

	for i := 0; i < 200*2; i++ {
		<-done
	}
}
