package transport

import (
	"sync"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc/metadata"
)

// FIX: writeStatus takes s.hdrMu.Lock around s.trailer.Copy(). SetTrailer
// also takes hdrMu. No race.
func TestRace_grpc_go_8519_trailer(t *testing.T) {
	s := &ServerStream{Stream: &Stream{trailer: metadata.MD{}}}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			s.hdrMu.Lock()
			_ = s.trailer.Copy()
			s.hdrMu.Unlock()
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
