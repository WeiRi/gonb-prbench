// Same-package whitebox race test for moby#29893
// Race target: storage.plugins map in pkg/plugins/plugins.go
//   Pre-fix: GetAll at line ~301 reads storage.plugins[name] WITHOUT storage lock.
//   Test triggers data race detector with frame in plugins.go (PR diff file).
package plugins

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestRace_29893(t *testing.T) {
	tmp, err := ioutil.TempDir("", "race_29893_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	backupSP := socketsPath
	backupSpecs := specsPaths
	socketsPath = tmp
	specsPaths = []string{tmp}
	defer func() { socketsPath = backupSP; specsPaths = backupSpecs }()

	const NAMES = 8
	for i := 0; i < NAMES; i++ {
		sockPath := filepath.Join(tmp, "p"+string(rune('a'+i))+".sock")
		l, err := net.Listen("unix", sockPath)
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
	}

	// Pre-populate storage so unlocked read at plugins.go:301 sees an entry.
	for i := 0; i < NAMES; i++ {
		nm := "p" + string(rune('a'+i))
		p := &Plugin{name: nm, Addr: "unix://" + filepath.Join(tmp, nm+".sock")}
		storage.Lock()
		storage.plugins[nm] = p
		storage.Unlock()
	}

	const N = 20
	const ITERS = 10
	var wg sync.WaitGroup
	wg.Add(N + N)

	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				nm := "p" + string(rune('a'+(id%NAMES)))
				p := &Plugin{name: nm}
				storage.Lock()
				storage.plugins[nm] = p
				storage.Unlock()
			}
		}(i)
	}

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			for j := 0; j < ITERS; j++ {
				_, _ = GetAll("anything")
			}
		}()
	}
	wg.Wait()
}
