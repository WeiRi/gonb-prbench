package contentdb

import (
	"sync"
	"testing"
)

// TestRaceConcurrentInit triggers data races on the buggy Init() function
// where multiple goroutines concurrently call Init() without mutex protection.
// We pre-initialize once so subsequent Init() calls skip loadDB(),
// avoiding fatal concurrent map writes while still racing on isInitialized + extDB.
func TestRaceConcurrentInit(t *testing.T) {
	// Pre-initialize: this makes isInitialized=true so subsequent Init()
	// calls skip loadDB(), avoiding fatal concurrent map write panics.
	if err := Init(); err != nil {
		t.Skipf("Init failed (expected in minimal workspace): %v", err)
	}

	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Race: concurrent Init() calls from multiple goroutines.
				// Bug: no mutex - isInitialized and extDB accessed without sync.
				// Even though isInitialized is already true, Init() still:
				// 1. Writes extDB = make(map[string]string) - RACE
				// 2. Reads isInitialized - RACE
				// 3. Writes isInitialized = true - RACE
				_ = Init()
			}
		}()
	}
	wg.Wait()
}

// TestRaceInitVsLookup triggers race between Init() writing shared state
// and Lookup()/MustLookup() reading it concurrently.
func TestRaceInitVsLookup(t *testing.T) {
	if err := Init(); err != nil {
		t.Skipf("Init failed (expected in minimal workspace): %v", err)
	}

	var wg sync.WaitGroup
	numWriters := 30
	numReaders := 30
	iterations := 200

	// Writers: call Init() which writes extDB + isInitialized
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = Init()
			}
		}()
	}

	// Readers: call Lookup/MustLookup which read isInitialized + extDB
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, _ = Lookup("txt")
				_, _ = Lookup("html")
				_ = MustLookup("json")
			}
		}()
	}

	wg.Wait()
}
