package auth

type authStore struct {
	revision uint64
}

// Revision — BUG (pre-PR7664): returns as.revision directly without sync (line 959).
func (as *authStore) Revision() uint64 {
	return as.revision // BUG line 959
}
