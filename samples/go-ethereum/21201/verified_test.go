package downloader

import (
	"sync"
	"testing"
)

// TestRace_21201_Downloader_mode: synchronise() writes d.mode while
// Progress() reads it concurrently — plain field, no atomic, fires under -race.
func TestRace_21201_Downloader_mode(t *testing.T) {
	d := NewDownloader()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			d.Synchronise(SyncMode(i % 3))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = d.Progress()
		}
	}()
	wg.Wait()
}
