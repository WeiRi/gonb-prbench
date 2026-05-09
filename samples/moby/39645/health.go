// Production stub for moby container/health.go (PR #39645).
// Pre-PR: callers read h.Health.Status directly (no mutex on read side).
// PR moves all reads through h.mu locking.
package container

import (
	"sync"
)

type HealthStatus struct {
	Status string
}

type Health struct {
	mu sync.Mutex
	Health HealthStatus
}
