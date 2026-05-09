// Race-trigger test for kubernetes-119729; see README.md for usage.

package scheduler

import (
	"sync"
	"testing"
)

func TestScheduleOneBindingFailureRace(t *testing.T) {
	s := &scheduler{}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pi := &PodInfo{Pod: &Pod{UID: "u", Name: "n"}}
			s.scheduleOne(pi)
		}()
	}
	wg.Wait()
	s.wg.Wait()
}
