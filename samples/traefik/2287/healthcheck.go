// Production stub for traefik healthcheck/healthcheck.go (PR #2287).
// Pre-PR: HealthCheck has no mutex; SetBackendsConfiguration assigns
// Backends map and iterates without sync.
package healthcheck

import (
	"context"
	"time"
)

type Options struct {
	Path     string
	Port     int
	Interval time.Duration
}

type BackendHealthCheck struct {
	Options Options
	URLs    []string
}

func NewBackendHealthCheck(opts Options) *BackendHealthCheck {
	return &BackendHealthCheck{Options: opts}
}

type HealthCheck struct {
	Backends map[string]*BackendHealthCheck
}

func newHealthCheck() *HealthCheck {
	return &HealthCheck{Backends: make(map[string]*BackendHealthCheck)}
}

// SetBackendsConfiguration writes hc.Backends and iterates without mutex.
func (hc *HealthCheck) SetBackendsConfiguration(ctx context.Context, backends map[string]*BackendHealthCheck) {
	hc.Backends = backends // RACE: map assignment without lock
	for name, b := range hc.Backends { // RACE: iteration without lock
		_ = name
		_ = b
	}
}
