package transport

import (
	"sync"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc/metadata"
)

// BUG: writeStatus reads s.trailer.Copy() WITHOUT s.hdrMu lock; concurrent
// SetTrailer writes s.trailer under hdrMu. Race on trailer map.
func TestRace_grpc_go_8519_trailer(t *testing.T) {
	s := &ServerStream{Stream: &Stream{trailer: metadata.MD{}}}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = s.trailer.Copy()
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = s.SetTrailer(metadata.MD{"k": []string{"v"}})
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
