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

// Inc — writes via atomic.AddUint64.
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

// clone — BUG (pre-PR11982): plain reads of fields race with atomic.AddUint64.
func (s *Stats) clone() Stats {
	return Stats{
		GetSuccess:    s.GetSuccess,    // line 90 BUG
		SetSuccess:    s.SetSuccess,
		DeleteSuccess: s.DeleteSuccess, // line 115 region
	}
}
