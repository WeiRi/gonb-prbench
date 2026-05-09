// Race-trigger test for serving-9297; see README.md for usage.

package statforwarder

import (
	"sync"
	"testing"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
)

func TestForwarderCancelRace(t *testing.T) {
	for iter := 0; iter < 50; iter++ {
		logger := zap.NewNop().Sugar()
		f := &Forwarder{
			selfIP:     "1.2.3.4",
			logger:     logger,
			processors: make(map[string]*bucketProcessor),
		}
		for i := 0; i < 5; i++ {
			f.processors[types.NamespacedName{Namespace: "ns", Name: string(rune('a' + i))}.String()] = &bucketProcessor{}
		}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			f.Cancel()
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				f.processors[types.NamespacedName{Namespace: "ns", Name: "new"}.String()] = &bucketProcessor{}
				delete(f.processors, types.NamespacedName{Namespace: "ns", Name: "new"}.String())
			}
		}()
		wg.Wait()
	}
}
