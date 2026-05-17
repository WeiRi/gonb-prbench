package modes

import (
	"encoding"
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	"github.com/fxamacker/cbor/v2"
)

type SafeCyclicTypeA struct {
	As []SafeCyclicTypeA
}

// Race: BUG checkers.getCheckerInternal does `c = checker{...}` and closes
// over the named return `c` in func literals. Concurrent getChecker calls
// on the same cyclic type race on shared `c`.
// FIX uses `placeholder := checker{...}` so each caller has own instance.
func TestRace_kubernetes_129170_lazy_checker_init(t *testing.T) {
	cache := checkers{
		cborInterface: reflect.TypeOf((*cbor.Marshaler)(nil)).Elem(),
		nonCBORInterfaces: []reflect.Type{
			reflect.TypeOf((*json.Marshaler)(nil)).Elem(),
			reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem(),
		},
	}

	rt := reflect.TypeOf(SafeCyclicTypeA{})
	begin := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-begin
			cache.getChecker(rt)
		}()
	}
	close(begin)
	wg.Wait()
}
