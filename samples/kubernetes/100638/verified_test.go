package queueset

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	fq "k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing"
	testclock "k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/testing/clock"
	"k8s.io/apiserver/pkg/util/flowcontrol/metrics"
)

func _newObsPair() metrics.TimedObserverPair {
	return metrics.PriorityLevelConcurrencyObserverPairGenerator.Generate(0, 0, []string{"x"})
}

// BUG: queueSet.StartRequest launches a goroutine that captures qs.qCfg.Name
// from the outer scope. PR #100638 captures configName locally INSIDE the lock
// before launching the goroutine. Race on qs.qCfg.Name field.
func TestRace_100638_qcfg_name_capture(t *testing.T) {
	metrics.Register()
	now := time.Now()
	clk, counter := testclock.NewFakeEventClock(now, 0, nil)
	qsf := NewQueueSetFactory(clk, counter)
	qCfg := fq.QueuingConfig{
		Name:             "race-100638",
		DesiredNumQueues: 4,
		QueueLengthLimit: 8,
		HandSize:         1,
		RequestWaitLimit: 10 * time.Second,
	}
	qsc, err := qsf.BeginConstruction(qCfg, _newObsPair())
	if err != nil {
		t.Fatal(err)
	}
	qs := qsc.Complete(fq.DispatchingConfig{ConcurrencyLimit: 4}).(*queueSet)

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500 && atomic.LoadInt32(&done) == 0; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			_, _ = qs.StartRequest(ctx, uint64(i), "flow", "fs", "d1", "d2", nil)
			cancel()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			qs.lock.Lock()
			qs.qCfg.Name = "mutated"
			qs.lock.Unlock()
		}
	}()
	wg.Wait()
}
