package logging

import (
	"sync"
	"testing"
)

func TestRace_OldestLogFileIdx(t *testing.T) {
	dir := t.TempDir()
	fr, err := NewFileRotator(dir, "racetest", 10, 1024*1024, nil)
	if err != nil {
		t.Fatalf("NewFileRotator: %v", err)
	}
	defer fr.Close()

	var wg sync.WaitGroup
	nWriters := 60
	nReaders := 60
	nIters := 500
	wg.Add(nWriters + nReaders)
	for i := 0; i < nWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				fr.purgeOldFiles(id + j)
			}
		}(i)
	}
	for i := 0; i < nReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				_ = fr.nextFile()
			}
		}()
	}
	wg.Wait()
}
