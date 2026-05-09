package quota

// Control holds the project-id quota state.
// Mirrors daemon/graphdriver/quota/types.go from moby PR #39644.
// The original lacked a sync.Mutex around the maps; we keep the same
// non-synchronized layout to reproduce the data race in prod code.
type Control struct {
	quotas        map[string]uint32
	nextProjectID uint32
}
