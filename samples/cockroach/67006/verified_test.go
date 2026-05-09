package kvserver

import (
	"sync"
	"testing"
)

func Test67006Race(t *testing.T) {
	rs := &replicaScanner{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		rs.Start(&Stopper{id: 1})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = rs.Monitor()
	}()
	wg.Wait()
}
