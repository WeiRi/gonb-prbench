// Race-trigger test for syncthing-10430; see README.md for usage.

package fs

import (
	"sync"
	"testing"
)

func TestCaseCacheRace(t *testing.T) {
	cache := newCaseCache()
	fs1 := &defaultRealCaser{cache: cache}
	fs2 := &defaultRealCaser{cache: cache}
	for iter := 0; iter < 200; iter++ {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ { fs1.getExpireAdd("test") }
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ { fs2.getExpireAdd("test") }
		}()
		wg.Wait()
	}
}
