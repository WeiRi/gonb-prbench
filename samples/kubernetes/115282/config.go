// Production stub for k8s staging/src/k8s.io/apiserver/pkg/server/config.go (PR #115282).
// Models DefaultBuildHandlerChain returning a handler chain with a shared
// warningRecorder accessed concurrently by the request goroutine and the
// timeout goroutine. BUG: no lock around warnings slice.
package apiserver

type warningRecorder struct {
	warnings []string
}

func (r *warningRecorder) AddWarning(s string) {
	r.warnings = append(r.warnings, s)
}

func (r *warningRecorder) Snapshot() []string {
	return append([]string(nil), r.warnings...)
}

// Chain is the value returned by DefaultBuildHandlerChain.
type Chain struct {
	rec *warningRecorder
}

// ServeRequest emulates a request running through the chain: the request
// goroutine adds warnings, the timeout goroutine snapshots warnings. BUG:
// both unsynchronized -> race on the warnings slice.
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

// DefaultBuildHandlerChain wires WithWarningRecorder into the chain in the
// position that exposes the racy slice to both timeout and request paths.
func DefaultBuildHandlerChain() *Chain {
	return &Chain{rec: &warningRecorder{}}
}
