package transport

import "sync"

type MD map[string][]string

func (m MD) Join(other MD) MD {
	out := MD{}
	for k, v := range m {
		out[k] = append(out[k], v...)
	}
	for k, v := range other {
		out[k] = append(out[k], v...)
	}
	return out
}

type Stream struct {
	hdrMu   sync.Mutex
	header  MD
	trailer MD
}

// SetTrailer — writes s.trailer under hdrMu.
func (s *Stream) SetTrailer(md MD) {
	s.hdrMu.Lock()
	s.trailer = s.trailer.Join(md) // line 177 write under hdrMu
	s.hdrMu.Unlock()
}

type serverHandlerTransport struct {
	rwLock sync.Mutex
}

// writeStatus — BUG (pre-PR8519): reads s.trailer in stats handler callback
// WITHOUT holding s.hdrMu (line 282).
func (sht *serverHandlerTransport) writeStatus(s *Stream) {
	for k, v := range s.trailer { // BUG line 282: unlocked iteration
		_ = k
		_ = v
	}
}
