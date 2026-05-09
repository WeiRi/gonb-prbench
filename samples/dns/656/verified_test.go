package dns

import (
	"testing"
)

func TestRaceConnSharedMutableState(t *testing.T) {
	// Handler that echoes back the query with correct ID
	handler := func(w ResponseWriter, req *Msg) {
		m := new(Msg)
		m.SetReply(req)
		w.WriteMsg(m)
	}

	HandleFunc("race656.", handler)
	defer HandleRemove("race656.")

	s, addrstr, err := RunLocalUDPServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer s.Shutdown()

	// Create a single shared connection
	co, err := Dial("udp", addrstr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer co.Close()

	m := new(Msg)
	m.SetQuestion("race656.", TypeA)

	numGoroutines := 60
	iterations := 100
	done := make(chan struct{})
	ready := make(chan struct{})

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			<-ready
			for j := 0; j < iterations; j++ {
				// WriteMsg writes co.t, ReadMsg reads co.t and writes co.rtt
				// Concurrent accesses race on shared Conn fields
				if err := co.WriteMsg(m); err != nil {
					continue
				}
				_, err := co.ReadMsg()
				if err != nil {
					continue
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
