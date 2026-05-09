// Production stub for syncthing lib/model/model.go (PR #9235).
// Pre-PR: ensureIndexHandler reads m.folderCfgs and m.folderRunners under pmut
// but they're protected by fmut. Missing fmut.RLock().
package model

import "sync"

type FolderConfiguration struct{ ID string }
type folderRunner struct{}

type model struct {
	fmut          sync.RWMutex // protects folderCfgs / folderRunners
	pmut          sync.RWMutex // protects something else
	folderCfgs    map[string]FolderConfiguration
	folderRunners map[string]*folderRunner
	indexHandlers map[string]struct{}
}

func newModel() *model {
	return &model{
		folderCfgs:    make(map[string]FolderConfiguration),
		folderRunners: make(map[string]*folderRunner),
		indexHandlers: make(map[string]struct{}),
	}
}

// ensureIndexHandler reads folderCfgs/folderRunners WITHOUT acquiring fmut (pre-PR bug).
func (m *model) ensureIndexHandler() error {
	m.pmut.Lock()
	defer m.pmut.Unlock()
	for id := range m.folderCfgs { // RACE: read map without fmut
		_ = m.folderRunners[id] // RACE
	}
	return nil
}

func (m *model) addFolder(id string) {
	m.fmut.Lock()
	m.folderCfgs[id] = FolderConfiguration{ID: id}
	m.folderRunners[id] = &folderRunner{}
	m.fmut.Unlock()
}

func (m *model) removeFolder(id string) {
	m.fmut.Lock()
	delete(m.folderCfgs, id)
	delete(m.folderRunners, id)
	m.fmut.Unlock()
}
