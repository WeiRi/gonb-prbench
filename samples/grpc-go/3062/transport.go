package transport

type Stream struct {
	headerChan   chan struct{}
	ctxDone      chan struct{}
	recvCompress string
}

// RecvCompress — BUG (pre-PR3062): in ctx.Done branch returns without waiting
// on headerChan, so concurrent writer in operateHeaders is still mid-write.
func (s *Stream) RecvCompress() string {
	select {
	case <-s.headerChan:
	case <-s.ctxDone:
		return s.recvCompress // line 32 BUG: read while writer still running
	}
	return s.recvCompress // line 37
}

// operateHeaders — writes recvCompress.
func (s *Stream) operateHeaders(c string) {
	s.recvCompress = c
	close(s.headerChan)
}
