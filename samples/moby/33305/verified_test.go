package loggerutils

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG LogPath() reads w.f without lock; Write triggers
// checkCapacityAndRotate which reassigns w.f when capacity exceeded.
// FIX adds w.mu.Lock to LogPath.
func TestRace_moby_33305_rotatefile_logpath(t *testing.T) {
	tmp := t.TempDir() + "/test.log"
	// tiny capacity → every Write triggers rotation → w.f reassigned
	w, err := NewRotateFileWriter(tmp, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	var done atomic.Bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		data := []byte("hello-world-hello-world")
		for i := 0; i < 200 && !done.Load(); i++ {
			_, _ = w.Write(data)
		}
		done.Store(true)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 200000 && !done.Load(); i++ {
			_ = w.LogPath()
		}
	}()
	wg.Wait()
}
