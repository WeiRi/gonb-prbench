// Minimal reproduction stub for kubernetes-128495
// Bug: InitPlugins() does NOT call refreshProbedPlugins() at the end.
// FindPluginByName calls refreshProbedPlugins() while only holding RLock.
// refreshProbedPlugins writes to probedPlugins map WITHOUT acquiring Lock.
// Concurrent FindPluginByName calls -> concurrent map writes -> DATA RACE.
//
// Pre-fix version: InitPlugins missing pm.refreshProbedPlugins().

package main

import (
	"fmt"
	"sync"
)

// Simplified interfaces
type VolumePlugin interface {
	GetPluginName() string
	CanSupport(spec *Spec) bool
	Init(host VolumeHost) error
}

type DynamicPluginProber interface {
	Init() error
	Probe() ([]ProbeEvent, error)
}

type ProbeEvent struct {
	Op          ProbeOperation
	Plugin      VolumePlugin
	PluginName  string
}

type ProbeOperation int

const (
	ProbeAddOrUpdate ProbeOperation = 1
	ProbeRemove      ProbeOperation = 2
)

type VolumeHost interface{}

type Spec struct {
	Name string
}

type VolumePluginMgr struct {
	mutex         sync.RWMutex
	plugins       map[string]VolumePlugin
	prober        DynamicPluginProber
	probedPlugins map[string]VolumePlugin
	Host          VolumeHost
}

// Simple plugin implementation
type testPlugin struct {
	name string
}

func (p *testPlugin) GetPluginName() string     { return p.name }
func (p *testPlugin) CanSupport(spec *Spec) bool { return false }
func (p *testPlugin) Init(host VolumeHost) error { return nil }

// Simple prober that always returns events (makes race reproducible)
type fakeProber struct {
	plugin VolumePlugin
}

func (p *fakeProber) Init() error { return nil }
func (p *fakeProber) Probe() ([]ProbeEvent, error) {
	// Always return events so every concurrent FindPluginByName
	// call will try to write to probedPlugins map.
	return []ProbeEvent{{
		Op:     ProbeAddOrUpdate,
		Plugin: p.plugin,
	}}, nil
}

func NewVolumePluginMgr(prober DynamicPluginProber, host VolumeHost) *VolumePluginMgr {
	return &VolumePluginMgr{
		plugins:       make(map[string]VolumePlugin),
		probedPlugins: make(map[string]VolumePlugin),
		prober:        prober,
		Host:          host,
	}
}

// InitPlugins - PRE-FIX: missing pm.refreshProbedPlugins() at end.
func (pm *VolumePluginMgr) InitPlugins(plugins []VolumePlugin, prober DynamicPluginProber, host VolumeHost) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.Host = host
	pm.prober = prober

	if err := pm.prober.Init(); err != nil {
		return err
	}

	for _, plugin := range plugins {
		name := plugin.GetPluginName()
		if _, found := pm.plugins[name]; found {
			return fmt.Errorf("volume plugin %q was registered more than once", name)
		}
		if err := plugin.Init(host); err != nil {
			return err
		}
		pm.plugins[name] = plugin
	}

	// PRE-FIX BUG: refreshProbedPlugins NOT called here!
	// pm.refreshProbedPlugins()
	return nil
}

// FindPluginByName fetches a plugin by name from both plugins and probedPlugins.
// It calls refreshProbedPlugins() while only holding RLock, which writes to
// probedPlugins map. Concurrent calls -> DATA RACE on map write.
func (pm *VolumePluginMgr) FindPluginByName(name string) (VolumePlugin, error) {
	pm.mutex.RLock()

	// refreshProbedPlugins called under RLock — writes to map!
	pm.refreshProbedPlugins()

	if plugin, found := pm.probedPlugins[name]; found {
		pm.mutex.RUnlock()
		return plugin, nil
	}

	if v, found := pm.plugins[name]; found {
		pm.mutex.RUnlock()
		return v, nil
	}

	pm.mutex.RUnlock()
	return nil, fmt.Errorf("no volume plugin matched name: %s", name)
}

// FindPluginBySpec — same pattern, calls refreshProbedPlugins under RLock.
func (pm *VolumePluginMgr) FindPluginBySpec(spec *Spec) (VolumePlugin, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	pm.refreshProbedPlugins()

	for _, plugin := range pm.probedPlugins {
		if plugin.CanSupport(spec) {
			return plugin, nil
		}
	}
	for _, v := range pm.plugins {
		if v.CanSupport(spec) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("no plugin found for spec")
}

// refreshProbedPlugins writes to probedPlugins map.
// Called from FindPluginByName/FindPluginBySpec under ONLY RLock.
// Multiple concurrent calls -> concurrent map writes -> DATA RACE.
func (pm *VolumePluginMgr) refreshProbedPlugins() {
	events, err := pm.prober.Probe()
	if err != nil {
		return
	}
	for _, event := range events {
		if event.Op == ProbeAddOrUpdate {
			// WRITE to probedPlugins map under only RLock!
			pm.probedPlugins[event.Plugin.GetPluginName()] = event.Plugin
		} else if event.Op == ProbeRemove {
			delete(pm.probedPlugins, event.PluginName)
		}
	}
}
