package license

import (
	"sync"
	"testing"
)

// Test130092Race reproduces the  (Missing HB) bug fixed by
// cockroachdb/cockroach#130092 / 130093 / 130094 / 130095.
func Test130092Race(t *testing.T) {
	instance = nil
	once = sync.Once{}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			e := GetEnforcerInstance()
			if idx%2 == 0 {
				e.Start("db_obj")
			} else {
				_ = e.DB()
			}
		}(i)
	}
	wg.Wait()
}
