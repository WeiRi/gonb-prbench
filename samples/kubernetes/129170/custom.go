// Stripped reproduction of staging/src/k8s.io/apimachinery/pkg/runtime/serializer/cbor/internal/modes/custom.go pre-PR #129170.
// BUG: getCheckerInternal stores &c (named return value) into a sync.Map via LoadOrStore;
// later writes to c (reassign), so concurrent readers fetch *&c and race with the writer.
package modes

import (
	"reflect"
	"runtime"
	"sync"
)

type checker struct {
	flag bool
	tag  int64
}

type checkers struct {
	m sync.Map
}

var marshalerCache = &checkers{}

func (cache *checkers) getChecker(rt reflect.Type) checker {
	return cache.getCheckerInternal(rt)
}

// BUG-state body: store &c then reassign — exposes named-return-value race.
func (cache *checkers) getCheckerInternal(rt reflect.Type) (c checker) {
	c = checker{flag: false, tag: 0}                   // initial value
	if actual, loaded := cache.m.LoadOrStore(rt, &c); loaded {
		return *actual.(*checker)                       // RACE READ on previous publisher's &c
	}
	// publisher path: simulate "compute final value" with yields so readers can race
	for i := 0; i < 8; i++ {
		runtime.Gosched()
	}
	c.flag = true                                       // RACE WRITE on &c (already published)
	c.tag++
	return c
}

type nonCBORInterface interface {
	MarshalText() ([]byte, error)
}

func RejectCustomMarshalers(v interface{}) error {
	rt := reflect.TypeOf(v)
	_ = marshalerCache.getChecker(rt)
	return nil
}

// resetCacheForType is used by the test to make every iteration trigger the publisher path.
func ResetCache(rt reflect.Type) {
	marshalerCache.m.Delete(rt)
}
