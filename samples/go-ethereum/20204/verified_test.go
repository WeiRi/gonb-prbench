package downloader

import (
	"testing"
)

// TestRace_20204_downloader_capture: drives the variable-capture race that
// PR #20204 fixes by hoisting the closure body into a function-typed
// closeOnErr taking *stateSync as a parameter.
func TestRace_20204_downloader_capture(t *testing.T) {
	d := &Downloader{}
	d.processFastSyncContent(50)
}
