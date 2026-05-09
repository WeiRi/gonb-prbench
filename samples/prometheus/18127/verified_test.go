package prometheus

import (
	"sync"
	"testing"
)

// TestRace_18127 reproduces the race condition in scrapeFailureLogger.
// Concurrent SetScrapeFailureLogger (write) and getScrapeFailureLogger (read).
func TestRace_18127(t *testing.T) {
	sp := &scrapePool{
		scrapeFailureLogger: noopFailureLogger{},
	}

	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	// Writers: call SetScrapeFailureLogger (write scrapeFailureLogger field)
	for g := 0; g < numGoroutines/2; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				sp.SetScrapeFailureLogger(noopFailureLogger{})
				sp.SetScrapeFailureLogger(nil)
			}
		}()
	}

	// Readers: call getScrapeFailureLogger (read scrapeFailureLogger field)
	for g := 0; g < numGoroutines/2; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = sp.getScrapeFailureLogger()
			}
		}()
	}

	wg.Wait()
}
