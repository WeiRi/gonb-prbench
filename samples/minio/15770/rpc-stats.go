// Pre-PR rpc-stats.go for minio-15770 race reproducer.
// BUG: dialStart/dialEnd are time.Time values written by ConnectStart/
// ConnectDone callbacks in a separate goroutine, while finisher reads them.
// Concurrent dial-cancel + finisher invocation produces a data race.
// Fix in PR #15770: use atomic.LoadInt64 / StoreInt64 of nanoseconds.
package buggy

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

// SetupReqStatsUpdate wraps the request with a httptrace.ClientTrace that
// records dialStart/dialEnd. The finisher returned reads them. In the buggy
// pre-fix version, those time.Time values are written by the dial goroutine
// without synchronization, racing with the finisher.
func SetupReqStatsUpdate(req *http.Request) (*http.Request, func()) {
	var dialStart, dialEnd time.Time
	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			dialStart = time.Now() // RACE: write
		},
		ConnectDone: func(network, addr string, err error) {
			dialEnd = time.Now() // RACE: write
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	finisher := func() {
		_ = dialStart // RACE: read
		_ = dialEnd   // RACE: read
	}
	return req, finisher
}
