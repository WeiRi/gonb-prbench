package whisperv6

import (
	"sync"
	"testing"
)

// TestRace_16146_bloomFilter_concurrent: drives concurrent bloomMatch (read)
// vs setBloomFilter (write) on a shared *Peer to provoke the data race that
// upstream PR #16146 fixes by adding peer.bloomMu sync.Mutex.
func TestRace_16146_bloomFilter_concurrent(t *testing.T) {
	peer := &Peer{}
	env := &Envelope{bloom: make([]byte, 64)}

	var wg sync.WaitGroup
	const N = 200

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			b := make([]byte, 64)
			for j := range b {
				b[j] = byte(i)
			}
			peer.setBloomFilter(b)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = peer.bloomMatch(env)
		}
	}()
	wg.Wait()
}
