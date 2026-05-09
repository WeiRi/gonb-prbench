package queueset

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	fq "k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing"
	"k8s.io/apiserver/pkg/util/flowcontrol/metrics"
	testclock "k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/testing/clock"
)

func TestRace_100638(t *testing.T) {
	metrics.Register()
	now := time.Now()
	clk, counter := testclock.NewFakeEventClock(now, 0, nil)
	qsf := NewQueueSetFactory(clk, counter)
	qCfg := fq.QueuingConfig{
		Name:             "TestRace100638",
		DesiredNumQueues: 9,
		QueueLengthLimit: 8,
		HandSize:         1,
		RequestWaitLimit: 10 * time.Minute,
	}
	qsc, err := qsf.BeginConstruction(qCfg, newObserverPair(clk))
	if err != nil {
		t.Fatal(err)
	}
	qs := qsc.Complete(fq.DispatchingConfig{ConcurrencyLimit: 4})

	var wg sync.WaitGroup
	n := 50

	// Goroutines that call StartRequest - the inner goroutine reads qs.qCfg.Name without lock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				ctx, cancel := context.WithCancel(context.Background())
				req, _ := qs.StartRequest(ctx, uint64(id), "test", "fs", "d1", "d2", nil)
				if req != nil {
					cancel() // trigger the inner goroutine to read qs.qCfg.Name
					req.Finish(func() {})
				} else {
					cancel()
				}
			}
		}(i)
	}

	// Goroutines that reconfigure the queueSet, writing qs.qCfg
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				newCfg := fq.QueuingConfig{
					Name:             fmt.Sprintf("TestRace100638-%d-%d", id, j),
					DesiredNumQueues: 9,
					QueueLengthLimit: 8,
					HandSize:         1,
					RequestWaitLimit: 10 * time.Minute,
				}
				completer, err := qs.BeginConfigChange(newCfg)
				if err != nil {
					continue
				}
				completer.Complete(fq.DispatchingConfig{ConcurrencyLimit: 4})
			}
		}(i)
	}

	wg.Wait()
}
