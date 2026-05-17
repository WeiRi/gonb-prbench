// FIXED: WaitGroup ensures goroutine completes BEFORE Reset(). cancel() moved
// to fire FIRST so goroutine sees ctx.Done() and exits before Reset.
package main

import (
	"context"
	"sync"
	"time"
)

type ResponseBuffer struct {
	data []byte
	done bool
}

func (b *ResponseBuffer) Reset() {
	b.data = nil
}

func (b *ResponseBuffer) Write(p []byte) {
	b.data = append(b.data, p...)
}

func handleMessageTunnel(buf *ResponseBuffer) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer func() {
		cancel()    // FIX: cancel FIRST so goroutine exits
		wg.Wait()   // FIX: wait for goroutine to finish
		buf.Reset() // safe now — goroutine is done
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ctx.Err() == nil {
			select {
			case <-ctx.Done():
				return
			default:
				buf.Write([]byte("msg"))
				time.Sleep(50 * time.Microsecond)
			}
		}
	}()
	time.Sleep(200 * time.Microsecond)
}
