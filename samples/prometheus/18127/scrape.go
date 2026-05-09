// Minimal extraction from prometheus scrape.go reproducing:
// race between SetScrapeFailureLogger (write) and getScrapeFailureLogger (read)
// on scrapePool.scrapeFailureLogger field.
package prometheus

import "io"

// FailureLogger is the interface for scrape failure logging.
// In the full prometheus this is slog.Handler + io.Closer.
type FailureLogger interface {
	io.Closer
}

// noopFailureLogger implements FailureLogger as a no-op.
type noopFailureLogger struct{}

func (noopFailureLogger) Close() error { return nil }

// scrapePool manages scrapes for sets of targets.
// BUGGY version: scrapeFailureLogger field accessed without mutex protection.
type scrapePool struct {
	scrapeFailureLogger FailureLogger
}

// SetScrapeFailureLogger sets the failure logger for the scrape pool.
// BUGGY: no mutex protection.
func (sp *scrapePool) SetScrapeFailureLogger(l FailureLogger) {
	sp.scrapeFailureLogger = l
}

// getScrapeFailureLogger returns the current failure logger.
// BUGGY: no mutex protection.
func (sp *scrapePool) getScrapeFailureLogger() FailureLogger {
	return sp.scrapeFailureLogger
}
