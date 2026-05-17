// Race test for prometheus-11296 — bstreamReader reads last byte while writer modifies it
package chunkenc

import (
	"sync"
	"testing"
)

func TestRace_11296_BStreamReaderLastByte(t *testing.T) {
	for round := 0; round < 100; round++ {
		shared := make([]byte, 16)
		for i := range shared {
			shared[i] = 0xff
		}
		// Pre-create the bstreamReader BEFORE concurrent goroutines
		r := newBReader(shared)
		var wg sync.WaitGroup
		wg.Add(2)
		// Writer: modifies last byte only (per PR's documented assumption)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				shared[len(shared)-1] = byte(i)
			}
		}()
		// Reader: drain bits to force loadNextBuffer eventually reading last byte
		go func() {
			defer wg.Done()
			for j := 0; j < 8*16; j++ {
				_, _ = r.readBit()
			}
		}()
		wg.Wait()
	}
}
