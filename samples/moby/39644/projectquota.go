package quota

// SetQuota assigns a project id to the given target.
// Mirrors daemon/graphdriver/quota/projectquota.go from moby PR #39644.
// Non-synchronized: races with GetQuota.
func (c *Control) SetQuota(target string) {
	c.quotas[target] = c.nextProjectID // RACE write
	c.nextProjectID++                  // RACE write
}

// GetQuota looks up an existing project id.
// Non-synchronized: races with SetQuota.
func (c *Control) GetQuota(target string) (uint32, bool) {
	id, ok := c.quotas[target] // RACE read
	return id, ok
}

// NewControl constructs a Control without locking primitives,
// matching the original PR's pre-fix layout.
func NewControl() *Control {
	return &Control{
		quotas:        make(map[string]uint32),
		nextProjectID: 1,
	}
}
