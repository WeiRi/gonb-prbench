// Pre-fix authenticated_group_adder.go from PR #109969.
// BUG: r.User = &DefaultInfo{Groups: append(r.User.GetGroups(), "system:authenticated")}
// — append may reuse the original (shared) slice's backing array, so two
// concurrent AuthenticateRequest calls on the SAME response race on the
// underlying array element.
package group

import "sync"

type DefaultInfo struct {
	Name   string
	UID    string
	Groups []string
}

func (d *DefaultInfo) GetName() string     { return d.Name }
func (d *DefaultInfo) GetUID() string      { return d.UID }
func (d *DefaultInfo) GetGroups() []string { return d.Groups }
func (d *DefaultInfo) GetExtra() map[string][]string {
	return nil
}

type Response struct {
	User *DefaultInfo
}

const AllAuthenticated = "system:authenticated"

type AuthenticatedGroupAdder struct {
	mu sync.Mutex
}

// AuthenticateRequest -- pre-fix path:
// r.User = &DefaultInfo{Groups: append(r.User.GetGroups(), AllAuthenticated)}
// Mutates shared backing array of caller-owned slice.
func (g *AuthenticatedGroupAdder) AuthenticateRequest(r *Response) (*Response, bool) {
	for _, g := range r.User.GetGroups() {
		if g == AllAuthenticated {
			return r, true
		}
	}
	r.User = &DefaultInfo{
		Name:   r.User.GetName(),
		UID:    r.User.GetUID(),
		Groups: append(r.User.GetGroups(), AllAuthenticated), // PRE-FIX: aliases backing array
	}
	return r, true
}
