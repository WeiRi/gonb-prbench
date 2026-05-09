package filters

import (
	"net/http"
)

// Stripped reproduction of staging/src/k8s.io/apiserver/pkg/server/filters/timeout.go pre-PR #107452.
// BUG: Header() returns the underlying writer's header map, so handler goroutine's Header().Set
// races with timeout goroutine's WriteHeader (which reads the same header map).

type baseTimeoutWriter struct {
	w http.ResponseWriter
}

// Header — BUG: returns the underlying writer's header map, no clone.
func (tw *baseTimeoutWriter) Header() http.Header {
	return tw.w.Header()
}

// WriteHeaderTimeout — fired by the timeout goroutine.
func (tw *baseTimeoutWriter) WriteHeaderTimeout(code int) {
	// internal http path: ranges over Header() (a map read) and writes to underlying writer.
	for k, v := range tw.w.Header() {
		_ = k
		_ = v
	}
	tw.w.WriteHeader(code)
}
