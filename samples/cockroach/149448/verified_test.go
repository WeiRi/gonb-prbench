package ui

import (
	"runtime"
	"testing"
)

func Test149448Race(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := &Spinner{}
		stop := s.Start()
		runtime.Gosched()
		stop()
	}
}
