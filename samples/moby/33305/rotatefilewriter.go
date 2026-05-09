// Production stub for moby daemon/logger/loggerutils/rotatefilewriter.go (PR #33305).
// Pre-PR: LogPath reads f.f.Name() while Write rotation replaces f.f without lock.
package loggerutils

import (
	"os"
	"sync"
)

type RotateFileWriter struct {
	mu       sync.Mutex
	f        *os.File
	capacity int64
	maxFiles int
	written  int64
}

func NewRotateFileWriter(path string, capacity int64, maxFiles int) (*RotateFileWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &RotateFileWriter{f: f, capacity: capacity, maxFiles: maxFiles}, nil
}

func (r *RotateFileWriter) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.written+int64(len(p)) > r.capacity {
		// rotate: close, rename, recreate (writes r.f without protecting LogPath readers)
		old := r.f
		old.Close()
		os.Rename(old.Name(), old.Name()+".1")
		nf, err := os.Create(old.Name())
		if err == nil {
			r.f = nf
			r.written = 0
		}
	}
	n, err := r.f.Write(p)
	r.written += int64(n)
	return n, err
}

// LogPath reads r.f.Name() WITHOUT taking r.mu (pre-PR bug).
func (r *RotateFileWriter) LogPath() string {
	return r.f.Name() // RACE: r.f replaced under lock during rotation
}

func (r *RotateFileWriter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f != nil {
		return r.f.Close()
	}
	return nil
}
