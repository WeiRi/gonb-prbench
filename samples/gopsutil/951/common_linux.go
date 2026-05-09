// Production stub for gopsutil/internal/common/common_linux.go (PR #951).
// Models virtualizationCache + virtualizationSystemMu pattern. Pre-PR #951
// the cache was a package-level map written/read concurrently.
package common

import "context"

var (
	cachedVirtMap     map[string]string
	cachedVirtSystem  string
	cachedVirtRole    string
)

// VirtualizationWithContext mirrors the racy cache fill+lookup path
// (pre-PR #951). The cache write & read happen without synchronization.
func VirtualizationWithContext(ctx context.Context) (string, string, error) {
	// Read-then-write the cache without locks (matches pre-PR behavior)
	if cachedVirtMap == nil {
		cachedVirtMap = make(map[string]string)
	}
	cachedVirtMap["last"] = "kvm"
	cachedVirtSystem = "kvm"
	cachedVirtRole = "guest"
	_ = cachedVirtMap["last"]
	return cachedVirtSystem, cachedVirtRole, nil
}
