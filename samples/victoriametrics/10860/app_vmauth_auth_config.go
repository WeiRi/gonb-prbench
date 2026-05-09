package vm10860repro

import (
	"context"
	"sync"
	"sync/atomic"
)

type backendURLs struct {
	healthChecksContext context.Context
	healthChecksCancel  func()
	healthChecksWG      sync.WaitGroup
	bus                 []*backendURL
}

type backendURL struct {
	broken             atomic.Bool
	healthCheckContext context.Context
	healthCheckWG      *sync.WaitGroup
}

func newBackendURLs() *backendURLs {
	ctx, cancel := context.WithCancel(context.Background())
	return &backendURLs{
		healthChecksContext: ctx,
		healthChecksCancel:  cancel,
	}
}

func (bus *backendURLs) add() *backendURL {
	bu := &backendURL{
		healthCheckContext: bus.healthChecksContext,
		healthCheckWG:      &bus.healthChecksWG,
	}
	bus.bus = append(bus.bus, bu)
	return bu
}

// app/vmauth/auth_config.go:434
func (bu *backendURL) setBroken() {
	if bu.broken.CompareAndSwap(false, true) {
		bu.healthCheckWG.Add(1)
		go func() {
			defer bu.healthCheckWG.Done()
			<-bu.healthCheckContext.Done()
			bu.broken.Store(false)
		}()
	}
}

// app/vmauth/auth_config.go:393
func (bus *backendURLs) stopHealthChecks() {
	bus.healthChecksCancel()
	bus.healthChecksWG.Wait()
}
