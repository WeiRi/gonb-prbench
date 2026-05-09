// Production stub for nats-server server/filestore.go (PR #3353).
// Pre-PR: populateGlobalPerSubjectInfo calls mb.readPerSubjectInfo without mb.mu.
package buggy

import "sync"

type msgBlock struct {
	mu       sync.Mutex
	perSubj  map[string]int
	dirty    bool
}

type fileStore struct {
	psim map[string]int
}

func (mb *msgBlock) writePerSubject(k string, v int) {
	mb.mu.Lock()
	if mb.perSubj == nil {
		mb.perSubj = make(map[string]int)
	}
	mb.perSubj[k] = v
	mb.dirty = true
	mb.mu.Unlock()
}

func (mb *msgBlock) readPerSubjectInfo(_ bool) map[string]int {
	return mb.perSubj // RACE: read without holding mb.mu
}

// populateGlobalPerSubjectInfo calls readPerSubjectInfo without mb.mu (pre-PR).
func (fs *fileStore) populateGlobalPerSubjectInfo(mb *msgBlock) {
	if fs.psim == nil {
		fs.psim = make(map[string]int)
	}
	psm := mb.readPerSubjectInfo(false)
	for k, v := range psm { // RACE: iterate map vs writePerSubject's writes
		fs.psim[k] += v
	}
}
