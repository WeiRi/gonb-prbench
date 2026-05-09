package modes

import (
	"reflect"
	"sync"
	"testing"
)

type raceTriggerType struct{}
func (raceTriggerType) MarshalText() ([]byte, error) { return nil, nil }

// Forces the named-return-value race in getCheckerInternal by clearing the
// cache between iterations so each round has a fresh publisher path.
func TestRace_PR129170_LazyCheckerInit(t *testing.T) {
	const numGoroutines = 50
	const iterations = 50

	rt := reflect.TypeOf(raceTriggerType{})

	for it := 0; it < iterations; it++ {
		ResetCache(rt)
		var wg sync.WaitGroup
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = marshalerCache.getChecker(rt)
			}()
		}
		wg.Wait()
	}
}
