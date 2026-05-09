// Race-trigger test for kubernetes-95664; see README.md for usage.

package record

import (
	"sync"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestRaceShutdownVsEventf(t *testing.T) {
	for iter := 0; iter < 50; iter++ {
		caster := NewBroadcasterForTests(0)
		fakeClock := clock.NewFakeClock(time.Now())
		recorder := recorderWithFakeClock(v1.EventSource{Component: "raceTest"}, caster, fakeClock)

		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				recorder.Eventf(&v1.ObjectReference{}, v1.EventTypeNormal, "Started", "blah")
			}()
		}

		caster.Shutdown()
		wg.Wait()
	}
}
