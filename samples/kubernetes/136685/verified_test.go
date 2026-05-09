package v1

import (
	"sync"
	"testing"
)

func TestRace_136685_SharedVerbsSlice(t *testing.T) {
	sharedVerbs := []string{"get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"}
	numGoroutines := 50
	iterations := 200
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				builder := NewRule(sharedVerbs...)
				builder.Resources("pods")
				builder.Groups("")
				_, _ = builder.Rule()
			}
		}()
	}
	wg.Wait()
}
