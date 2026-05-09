package connrotation

import (
	"context"
	"net"
	"sync"
	"testing"
)

type fakeConn struct {
	net.Conn
}

func (f *fakeConn) Close() error {
	return nil
}

func TestRace_88079(t *testing.T) {
	d := NewDialer(func(ctx context.Context, network, address string) (net.Conn, error) {
		return &fakeConn{}, nil
	})

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		// Writer: DialContext adds conn to map, then sets onClose
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				d.Dial("tcp", "127.0.0.1:0")
			}
		}()
		// Reader/Writer: CloseAll reads conns, closes connections,
		// which reads onClose. Races with DialContext writing onClose.
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				d.CloseAll()
			}
		}()
	}

	wg.Wait()
}
