package disk

import (
	"sync"
	"time"
)

type Stats struct {
	BytesRead    int64
	BytesWritten int64
	Updated      time.Time
}

type Monitor struct {
	tracer *monitorTracer
}

type monitorTracer struct {
	mu    sync.Mutex
	stats Stats
}

func NewMonitor() *Monitor {
	m := &Monitor{tracer: &monitorTracer{}}
	go m.tracer.tracker()
	return m
}

func (m *Monitor) IncrementalStats() Stats {
	m.tracer.mu.Lock()
	defer m.tracer.mu.Unlock()
	return m.tracer.stats
}

func (t *monitorTracer) tracker() {
	for i := 0; i < 1000; i++ {
		t.mu.Lock()
		t.stats.BytesRead = int64(i)
		t.stats.BytesWritten = int64(i * 2)
		t.stats.Updated = time.Now()
		t.mu.Unlock()
	}
}
