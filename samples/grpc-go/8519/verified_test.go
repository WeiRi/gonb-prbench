package transport

import (
	"sync"
	"testing"
)

func TestRace_HandlerServer_TrailerCopy_8519(t *testing.T) {
	const N = 50
	const ITERS = 200

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := &Stream{header: MD{}, trailer: MD{"k": {"v"}}}
			sht := &serverHandlerTransport{}

			var inner sync.WaitGroup
			inner.Add(2)
			go func() {
				defer inner.Done()
				for j := 0; j < ITERS; j++ {
					s.SetTrailer(MD{"rk": {"rv"}})
				}
			}()
			go func() {
				defer inner.Done()
				for j := 0; j < ITERS; j++ {
					sht.writeStatus(s)
				}
			}()
			inner.Wait()
		}()
	}
	wg.Wait()
}
