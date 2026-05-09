package chunkenc

import (
	"sync"
	"testing"
)

// TestRace_PR11296_BstreamRace triggers the data race between bstreamReader
// (used by XOR chunk iterators) reading the last byte of the stream and a
// concurrent Appender (xorAppender) writing to it.
//
// The bug: bstreamReader.loadNextBuffer() reads from the stream slice while
// bstream.writeBit/writeByte concurrently modify the last byte of the same
// slice. No synchronization between reader and writer in the buggy version.
//
// Fix in PR 11296: copy the last byte at initialization time to avoid the race.
func TestRace_PR11296_BstreamRace(t *testing.T) {
	const numGoroutines = 100
	const iterations = 500

	for iter := 0; iter < 5; iter++ {
		var wg sync.WaitGroup

		// Create a chunk with initial data.
		chk := NewXORChunk()
		app, _ := chk.Appender()
		for i := 0; i < 20; i++ {
			app.Append(int64(i*1000), float64(i)*0.5)
		}

		// Writer goroutines: keep appending samples.
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					app, _ := chk.Appender()
					app.Append(int64((20+i)*1000), float64(i)*0.5)
				}
			}(g)
		}

		// Reader goroutines: read from the chunk via iterator.
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					it := chk.Iterator(nil)
					for it.Next() {
						_, _ = it.At()
					}
				}
			}(g)
		}

		wg.Wait()
	}
}
