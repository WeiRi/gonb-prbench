// Race-trigger test for syncthing-9235; see README.md for usage.

package model

import (
	"sync"
	"testing"
)

func TestEnsureIndexHandlerRace(t *testing.T) {
	for iter := 0; iter < 100; iter++ {
		m := newModel()
		// Populate data
		m.fmut.Lock()
		m.folderCfgs["folder1"] = FolderConfiguration{ID: "folder1"}
		m.folderRunners["folder1"] = &folderRunner{}
		m.folderCfgs["folder2"] = FolderConfiguration{ID: "folder2"}
		m.folderRunners["folder2"] = &folderRunner{}
		m.fmut.Unlock()

		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				_ = m.ensureIndexHandler()
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				m.addFolder("newfolder")
				m.removeFolder("newfolder")
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				_ = m.ensureIndexHandler()
			}
		}()
		wg.Wait()
	}
}
