package v2store

import "sync/atomic"

const (
	GetSuccess    = "GetSuccess"
	SetSuccess    = "SetSuccess"
	DeleteSuccess = "DeleteSuccess"
)

type Stats struct {
	GetSuccess    uint64
	SetSuccess    uint64
	DeleteSuccess uint64
}

func newStats() *Stats { return &Stats{} }

func (s *Stats) Inc(field string) {
	switch field {
	case GetSuccess:
		atomic.AddUint64(&s.GetSuccess, 1)
	case SetSuccess:
		atomic.AddUint64(&s.SetSuccess, 1)
	case DeleteSuccess:
		atomic.AddUint64(&s.DeleteSuccess, 1)
	}
}

// FIX (PR11982): use atomic.LoadUint64 for reads to match Inc's atomic writes.
func (s *Stats) clone() Stats {
	return Stats{
		GetSuccess:    atomic.LoadUint64(&s.GetSuccess),
		SetSuccess:    atomic.LoadUint64(&s.SetSuccess),
		DeleteSuccess: atomic.LoadUint64(&s.DeleteSuccess),
	}
}
