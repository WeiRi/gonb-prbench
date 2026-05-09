package common

import (
	"context"
	"testing"
)

func TestRaceVirtualizationCacheConcurrent(t *testing.T) {
	iterations := 400
	numGoroutines := 80

	done := make(chan struct{})

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			ctx := context.Background()
			for j := 0; j < iterations; j++ {
				// Each call reads/writes virtualizationCache concurrently
				VirtualizationWithContext(ctx)
			}
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
