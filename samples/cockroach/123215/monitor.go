// Production stub for cockroach pkg/storage/disk/monitor.go (PR #123215).
// Models Monitor struct with mutable stats fields read concurrently without locks.
package disk

import "time"

type Stats struct {
	BytesRead    int64
	BytesWritten int64
	Updated      time.Time
}

type Monitor struct {
	tracer *monitorTracer
}

type monitorTracer struct {
	stats Stats // racy: written by tracker, read by IncrementalStats without lock
}

func NewMonitor() *Monitor {
	m := &Monitor{tracer: &monitorTracer{}}
	go m.tracer.tracker()
	return m
}

// IncrementalStats reads stats without acquiring the tracer mutex (pre-PR bug).
func (m *Monitor) IncrementalStats() Stats {
	return m.tracer.stats
}

// tracker simulates the background updater that writes stats every tick.
func (t *monitorTracer) tracker() {
	for i := 0; i < 1000; i++ {
		t.stats.BytesRead = int64(i)
		t.stats.BytesWritten = int64(i * 2)
		t.stats.Updated = time.Now()
	}
}
