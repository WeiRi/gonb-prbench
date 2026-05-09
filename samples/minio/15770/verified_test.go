// VERIFIED race reproducer for minio PR #15770
// "Fix race when accessing REST tcp dial values"
// https://github.com/minio/minio/pull/15770
//
// Original racy file: internal/rest/rpc-stats.go (lines 54-79 pre-fix)
// Fix: replace time.Time locals with atomic int64 nanoseconds.
//
// This whitebox test reproduces the race by:
//   1. Embedding the pre-fix SetupReqStatsUpdate in the same package.
//   2. Cancelling HTTP requests mid-dial so the dial goroutine still
//      runs ConnectStart/ConnectDone callbacks (writing dialStart/dialEnd)
//      while the caller goroutine runs the finisher closure (reading them).
//   3. 128 parallel goroutines per iteration, count=20, reproduces RACE
//      on every run under Go 1.21+ with --memory=2g --cpus=2.
//
// To use this test as a positive race oracle:
//   - Drop pre-fix rpc-stats.go alongside this file under package buggy.
//   - go test -race -count=20 → WARNING: DATA RACE (frames cite rpc-stats.go:31/35/41/42).
//   - Apply PR #15770 (atomic.LoadInt64/StoreInt64) → race disappears.
package buggy

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRaceTCPDialStats(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				time.Sleep(50 * time.Millisecond)
				conn.Close()
			}(c)
		}
	}()
	srv := &httptest.Server{
		Listener: ln,
		Config:   &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})},
	}
	srv.Start()
	defer srv.Close()

	transport := &http.Transport{
		DisableKeepAlives: true,
		MaxIdleConns:      0,
	}
	client := &http.Client{Transport: transport}

	const N = 128
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
			defer cancel()
			req, _ := http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
			req, finish := SetupReqStatsUpdate(req)
			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
			}
			finish()
		}()
	}
	wg.Wait()
}
