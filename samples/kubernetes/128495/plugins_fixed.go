// FIX version for kubernetes-128495 (PR #128495):
// InitPlugins now calls refreshProbedPlugins() under Lock, and
// FindPluginByName/FindPluginBySpec no longer call refreshProbedPlugins
// (which writes to map) — they only read.

package main

import (
	"fmt"
	"sync"
)

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
	Op         ProbeOperation
	Plugin     VolumePlugin
	PluginName string
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

type testPlugin struct {
	name string
}

func (p *testPlugin) GetPluginName() string      { return p.name }
func (p *testPlugin) CanSupport(spec *Spec) bool { return false }
func (p *testPlugin) Init(host VolumeHost) error { return nil }

type fakeProber struct {
	plugin VolumePlugin
}

func (p *fakeProber) Init() error { return nil }
func (p *fakeProber) Probe() ([]ProbeEvent, error) {
	return []ProbeEvent{{Op: ProbeAddOrUpdate, Plugin: p.plugin}}, nil
}

func NewVolumePluginMgr(prober DynamicPluginProber, host VolumeHost) *VolumePluginMgr {
	return &VolumePluginMgr{
		plugins:       make(map[string]VolumePlugin),
		probedPlugins: make(map[string]VolumePlugin),
		prober:        prober,
		Host:          host,
	}
}

// InitPlugins — FIX: calls refreshProbedPluginsLocked under Lock to populate map at init time.
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

	pm.refreshProbedPluginsLocked()
	return nil
}

// FindPluginByName — FIX: only read under RLock; never write.
func (pm *VolumePluginMgr) FindPluginByName(name string) (VolumePlugin, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	if plugin, found := pm.probedPlugins[name]; found {
		return plugin, nil
	}
	if v, found := pm.plugins[name]; found {
		return v, nil
	}
	return nil, fmt.Errorf("no volume plugin matched name: %s", name)
}

func (pm *VolumePluginMgr) FindPluginBySpec(spec *Spec) (VolumePlugin, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
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

// refreshProbedPluginsLocked — caller MUST hold pm.mutex.Lock() (write lock).
func (pm *VolumePluginMgr) refreshProbedPluginsLocked() {
	events, err := pm.prober.Probe()
	if err != nil {
		return
	}
	for _, event := range events {
		if event.Op == ProbeAddOrUpdate {
			pm.probedPlugins[event.Plugin.GetPluginName()] = event.Plugin
		} else if event.Op == ProbeRemove {
			delete(pm.probedPlugins, event.PluginName)
		}
	}
}
