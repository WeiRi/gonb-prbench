package txpool

import (
	"sync"
	"testing"
)

// TestRace_14940_journal_rotate_vs_reset: RotateJournal reads
// pool.pending/locals without pool.mu while Reset writes under pool.mu —
// fires under -race.
func TestRace_14940_journal_rotate_vs_reset(t *testing.T) {
	pool := NewTxPool()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			pool.Reset(i % 16)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = pool.RotateJournal()
		}
	}()
	wg.Wait()
}
