package apiserver

import "sync"

type warningRecorder struct {
	mu       sync.Mutex
	warnings []string
}

func (r *warningRecorder) AddWarning(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.warnings = append(r.warnings, s)
}

func (r *warningRecorder) Snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.warnings...)
}

type Chain struct {
	rec *warningRecorder
}

func (c *Chain) ServeRequest() {
	done := make(chan struct{})
	go func() {
		for i := 0; i < 200; i++ {
			c.rec.AddWarning("w")
		}
		close(done)
	}()
	for i := 0; i < 200; i++ {
		_ = c.rec.Snapshot()
	}
	<-done
}

func DefaultBuildHandlerChain() *Chain {
	return &Chain{rec: &warningRecorder{}}
}
