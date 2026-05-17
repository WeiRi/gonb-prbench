package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"google.golang.org/grpc/status"
)

type flusherRecorder struct {
	*httptest.ResponseRecorder
}

func (f *flusherRecorder) Flush() {}
func (f *flusherRecorder) CloseNotify() <-chan bool {
	return make(chan bool, 1)
}

// BUG: WriteStatus calls close(ht.writes) but DOES NOT close ht.closedCh first.
// A subsequent ht.do() (called from Write/WriteHeader) does select on
// `ht.writes <- fn` (BUG: no closedCh, so picks send) → panic: send on closed channel.
// PR #1687 calls ht.Close() (closes closedCh) BEFORE close(writes) so do() sees
// <-closedCh and returns ErrConnClosing.
func TestRace_1687_use_after_writestatus(t *testing.T) {
	for iter := 0; iter < 50; iter++ {
		rw := &flusherRecorder{httptest.NewRecorder()}
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		ht := &serverHandlerTransport{
			rw:       rw,
			req:      req,
			writes:   make(chan func(), 1),
			closedCh: make(chan struct{}),
		}
		// drain writes (must be running so do() can send)
		drainDone := make(chan struct{})
		go func() {
			defer close(drainDone)
			for fn := range ht.writes {
				fn()
			}
		}()

		// Sequence: WriteStatus first, then concurrent do() calls
		st := status.New(0, "")
		s := &Stream{ctx: context.Background()}
		_ = ht.WriteStatus(s, st)
		<-drainDone

		// Now ht.writes is closed. Subsequent do() should panic in BUG.
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ht.do(func() {})
		}()
		wg.Wait()
	}
}
