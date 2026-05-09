package transport

type HF struct {
	Name  string
	Value string
}

type Stream struct {
	header       map[string][]string
	headerChan   chan struct{}
	recvCompress string
}

func operateHeaders(s *Stream, hfs []HF) {
	for _, hf := range hfs {
		if hf.Name == "grpc-encoding" {
			s.recvCompress = hf.Value // line 25 BUG mid-parse write
		}
		s.header[hf.Name] = []string{hf.Value}
	}
	close(s.headerChan)
}

func readerObservesAfterChan(s *Stream) string {
	<-s.headerChan
	return s.recvCompress // line 46
}

func raceWriter(s *Stream) {
	s.recvCompress = "racing-writer"
}
