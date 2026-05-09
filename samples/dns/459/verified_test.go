package dns

import (
	"testing"
)

func TestRaceExchangeSharedResponseError(t *testing.T) {
	// Server that returns wrong-ID response (causes ErrId: r != nil, err != nil)
	handler := func(w ResponseWriter, req *Msg) {
		m := new(Msg)
		m.SetReply(req)
		m.Id = req.Id + 1
		w.WriteMsg(m)
	}

	HandleFunc("race459.", handler)
	defer HandleRemove("race459.")

	s, addrstr, err := RunLocalUDPServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer s.Shutdown()

	c := new(Client)
	c.SingleInflight = true
	m := new(Msg)
	m.SetQuestion("race459.", TypeA)

	numGoroutines := 80
	done := make(chan struct{})
	ready := make(chan struct{})

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			<-ready
			for j := 0; j < 10; j++ {
				r, _, _ := c.Exchange(m, addrstr)
				if r != nil {
					r.Compress = !r.Compress
					r.Truncated = !r.Truncated
					_ = r.String()
				}
			}
			done <- struct{}{}
		}(i)
	}

	close(ready)

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
