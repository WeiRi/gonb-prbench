package transport

import "sync"

type MD map[string][]string

type Stream struct {
	mu       sync.Mutex
	header   MD
	headerOk bool
	trailer  MD
}

func NewStream() *Stream {
	return &Stream{header: MD{}, trailer: MD{}}
}

// WriteHeader — BUG (pre-PR2074): writes s.header / s.headerOk without sync.
func (s *Stream) WriteHeader(md MD) error {
	for k, v := range md { // line 22 read md
		s.header[k] = v // line 25 write s.header
	}
	s.headerOk = true
	return nil
}

// WriteStatus — BUG: also reads s.header / s.headerOk without sync.
func (s *Stream) WriteStatus() error {
	if !s.headerOk { // line 44
		s.headerOk = true // line 45
	}
	return nil
}
