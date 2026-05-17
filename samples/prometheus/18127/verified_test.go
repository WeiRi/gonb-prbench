// Race test for prometheus-18127 — scrapePool.scrapeFailureLogger raced via two mutexes
// BUG: SetScrapeFailureLogger uses scrapeFailureLoggerMtx; restartLoops reads sp.scrapeFailureLogger
//      while holding only targetMtx → race between two mutexes
// FIX: consolidates write under targetMtx, removes scrapeFailureLoggerMtx
package scrape

import (
	"sync"
	"testing"
)

func TestRace_18127_ScrapeFailureLoggerMtx(t *testing.T) {
	sp := &scrapePool{}
	const N = 200
	var wg sync.WaitGroup
	wg.Add(2)
	// Goroutine A: SetScrapeFailureLogger (BUG: scrapeFailureLoggerMtx)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			sp.SetScrapeFailureLogger(nil)
		}
	}()
	// Goroutine B: simulate restartLoops-style read (BUG: only targetMtx)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			sp.targetMtx.Lock()
			_ = sp.scrapeFailureLogger
			sp.targetMtx.Unlock()
		}
	}()
	wg.Wait()
}
