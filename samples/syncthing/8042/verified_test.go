// Race-trigger test for syncthing-8042; see README.md for usage.

package nat

import (
	"sync"
	"testing"
	"time"
)

func TestMappingRace(t *testing.T) {
	for iter := 0; iter < 200; iter++ {
		m := &Mapping{extAddresses: make(map[string]Address), expires: time.Now().Add(time.Hour)}
		s := &Service{mappings: []*Mapping{m}}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ { s.process() }
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ { s.updateMapping(m) }
		}()
		wg.Wait()
	}
}
