package grid

import (
	"sync"
	"testing"
)

// TestRace_BufferLifecycle reproduces the use-after-free race on byte buffers.
// The bug: PutByteBuffer returns buffer to the pool before cancelFn stops
// goroutines using the buffer. After PutByteBuffer, another goroutine
// can get the same buffer from the pool and write to it, while the
// original goroutine still reads from it.
//
// 70 goroutines get buffers, pass them to other goroutines that read from
// them, return buffers to pool, while new goroutines get and write to the
// same buffers -- triggering the race detector on concurrent read/write.
func TestRace_BufferLifecycle(t *testing.T) {
	var wg sync.WaitGroup
	nWorkers := 70
	nIters := 300

	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				// Get a buffer from the pool
				buf := GetByteBuffer()
				buf = append(buf, byte(id), byte(j), byte(id^j))

				// Pass the buffer to a "handler" goroutine that reads it
				done := make(chan struct{})
				go func(data []byte) {
					// Simulate handler processing: read from the buffer
					_ = data[0] + data[1] + data[2]
					close(done)
				}(buf)

				// Bug: return buffer to pool BEFORE handler is done
				PutByteBuffer(buf)

				// Another goroutine might get the same buffer from pool
				buf2 := GetByteBuffer()
				buf2 = append(buf2, byte(255-id), byte(255-j), byte(255))

				// Write to potentially-the-same-buffer while handler reads
				buf2[0] = byte(id + j)

				PutByteBuffer(buf2)

				// Wait for the handler to finish
				<-done
			}
		}(i)
	}

	wg.Wait()
}
