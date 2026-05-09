package volume

import "sync"

// Stripped reproduction of pkg/volume/plugins.go pre-PR #96826.
// BUG: the Get path reads the plugins map without holding the mutex,
// while Add holds the mutex. Concurrent map access fires the race.

type VolumePlugin interface{}

type VolumePluginMgr struct {
	mutex   sync.Mutex
	plugins map[string]VolumePlugin
}

// AddPlugin — locked write.
func (pm *VolumePluginMgr) AddPlugin(name string, p VolumePlugin) {
	pm.mutex.Lock()
	pm.plugins[name] = p
	pm.mutex.Unlock()
}

// GetPlugin — BUG: reads map without lock.
func (pm *VolumePluginMgr) GetPlugin(name string) (VolumePlugin, bool) {
	v, ok := pm.plugins[name]
	return v, ok
}
