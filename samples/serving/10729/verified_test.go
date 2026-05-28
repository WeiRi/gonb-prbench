// Regression test for knative/serving#10729

package handler

import (
	"sync"
	"testing"
	"time"

	network "knative.dev/networking/pkg"
)

func TestConcurrencyReporterRace_10729(t *testing.T) {
	cr, _, cancel := newTestReporter(t)
	defer cancel()
	base := time.Now()

	const W = 10
	const R = 3
	stop := time.After(3 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < W; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j := 0
			for {
				select {
				case <-stop:
					return
				default:
					cr.handleEvent(network.ReqEvent{Time: base.Add(time.Duration(j) * time.Microsecond), Type: network.ReqIn, Key: rev1})
					cr.handleEvent(network.ReqEvent{Time: base.Add(time.Duration(j) * time.Microsecond), Type: network.ReqOut, Key: rev1})
					j++
				}
			}
		}()
	}

	for i := 0; i < R; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j := 0
			for {
				select {
				case <-stop:
					return
				default:
					_ = cr.report(base.Add(time.Duration(j) * time.Millisecond))
					j++
				}
			}
		}()
	}

	wg.Wait()
}
