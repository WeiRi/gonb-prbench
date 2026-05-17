// Post-PR rpc-stats.go for minio-15770 race reproducer (FIXED version).
// Fix: use atomic.LoadInt64 / StoreInt64 instead of time.Time.
package buggy

import (
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

func SetupReqStatsUpdate(req *http.Request) (*http.Request, func()) {
	var dialStart, dialEnd int64
	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			atomic.StoreInt64(&dialStart, time.Now().UnixNano())
		},
		ConnectDone: func(network, addr string, err error) {
			if err == nil {
				atomic.StoreInt64(&dialEnd, time.Now().UnixNano())
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	finisher := func() {
		if ds := atomic.LoadInt64(&dialStart); ds > 0 {
			if de := atomic.LoadInt64(&dialEnd); de == 0 {
				// timeout
			} else if de >= ds {
				_ = de - ds
			}
		}
	}
	return req, finisher
}
