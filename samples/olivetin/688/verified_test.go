package entities

import (
	"fmt"
	"sync"
	"testing"
)

// TestRace_688_EntitiesMap reproduces PR #688:
// GetEntities() returns shared map ref; readers iterate while AddEntity writes → map race.
func TestRace_688_EntitiesMap(t *testing.T) {
	// reset contents to fresh state for test isolation
	rwmutex.Lock()
	contents = &variableBase{Entities: make(entitiesByClass, 0)}
	rwmutex.Unlock()

	const N = 30
	const ITERS = 200
	var wg sync.WaitGroup
	// Writers
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				AddEntity("class-a", fmt.Sprintf("k-%d-%d", i, j), j)
			}
		}(i)
	}
	// Readers (iterate map → race with concurrent writes)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				m := GetEntities()
				for _, inst := range m {
					for k := range inst {
						_ = k
					}
				}
			}
		}()
	}
	wg.Wait()
}
