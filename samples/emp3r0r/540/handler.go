// Stub reproduction of emp3r0r-540 PR:
// handleMessageTunnel spawns a goroutine that writes to a shared buffer.
// BUG: handler returns without waiting; goroutine writes after caller already
// freed/reused the buffer → data race on buffer state.
package main

import (
	"context"
	"time"
)

type ResponseBuffer struct {
	data []byte
	done bool
}

func (b *ResponseBuffer) Reset() {
	b.data = nil // race write
}

func (b *ResponseBuffer) Write(p []byte) {
	b.data = append(b.data, p...) // race write
}

// BUG: spawns goroutine but doesn't wait for it. cancel() runs LAST in defer
// (after Reset/free), so goroutine may still run and write after Reset.
func handleMessageTunnel(buf *ResponseBuffer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		buf.Reset() // caller "frees" buffer
		cancel()    // BUG: cancel runs AFTER Reset → goroutine still active
	}()

	go func() {
		defer cancel()
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
	// "handler" returns after a brief moment, before goroutine done
	time.Sleep(200 * time.Microsecond)
}
