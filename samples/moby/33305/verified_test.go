package loggerutils

import (
	"os"
	"sync"
	"testing"
)

func TestRace_33305(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-rotate-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)
	defer os.Remove(tmpPath + ".1")

	const N = 20
	const ITERS = 500

	for trial := 0; trial < 10; trial++ {
		// Small capacity to force frequent rotations
		w, err := NewRotateFileWriter(tmpPath, 100, 2)
		if err != nil {
			t.Fatal(err)
		}

		var wg sync.WaitGroup
		wg.Add(N * 2)

		for i := 0; i < N; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < ITERS; j++ {
					w.Write([]byte("test log data that will fill up the small buffer quickly to trigger rotation\n"))
				}
			}(i)
		}

		for i := 0; i < N; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < ITERS; j++ {
					_ = w.LogPath()
				}
			}()
		}

		wg.Wait()
		w.Close()
	}
}
