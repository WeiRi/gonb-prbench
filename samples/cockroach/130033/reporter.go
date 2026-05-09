package diagnostics

import "time"

// Reproduction of PR cockroachdb/cockroach#130031 (and 130032/130033 backports):
// "server: make telemetry timestamp atomic". BUG state: time.Time field
// is read/written by multiple goroutines without synchronization.

type Reporter struct {
	// LastSuccessfulTelemetryPing is read/written without sync (BUG).
	LastSuccessfulTelemetryPing time.Time
}

// ReportDiagnostics writes the timestamp (BUG: unsynchronized).
func (r *Reporter) ReportDiagnostics() {
	r.LastSuccessfulTelemetryPing = time.Now()
}

// Read reads the timestamp (BUG: unsynchronized).
func (r *Reporter) Read() time.Time {
	return r.LastSuccessfulTelemetryPing
}

