package admission

import (
	"fmt"
	"sync"
	"testing"
)

func TestRace_106045_AuditAnnotations(t *testing.T) {
	const N = 100
	var wg sync.WaitGroup
	for n := 0; n < N; n++ {
		h := NewAuditHandler()
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				h.Admit(fmt.Sprintf("k%d", i), "v")
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				h.Admit(fmt.Sprintf("z%d", i), "w")
			}
		}()
	}
	wg.Wait()
}
