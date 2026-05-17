package main

import (
	"sync"
	"testing"
)

func TestRace_540_HandlerLifecycle(t *testing.T) {
	const N = 30
	const ITERS = 20
	var wg sync.WaitGroup
	for trial := 0; trial < ITERS; trial++ {
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				buf := &ResponseBuffer{}
				handleMessageTunnel(buf)
			}()
		}
		wg.Wait()
	}
}
