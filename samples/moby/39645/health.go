package container

import "sync"

type HealthStatus struct {
	Status string
}

type Health struct {
	mu     sync.Mutex
	Health HealthStatus
}

// SetHealth writes status under lock (fix applied), but also exposes racy path
func (h *Health) SetHealth(s string) {
	h.mu.Lock()
	h.Health.Status = s
	h.mu.Unlock()
}

// RacyReadHealth reads status WITHOUT lock (BUG)
func (h *Health) RacyReadHealth() string {
	return h.Health.Status // RACE: read without lock
}

// RacyWriteHealth writes status WITHOUT lock (BUG)
func (h *Health) RacyWriteHealth(s string) {
	h.Health.Status = s // RACE: write without lock
}
