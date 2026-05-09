// Production stub for kubernetes-84487.
// Pre-PR: WatchNodes reads r.Err without holding r.Mutex.
package registrytest

import (
	"context"
	"sync"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
)

type NodeRegistry struct {
	sync.Mutex
	Err error // RACE: read by WatchNodes (no lock), written by tests
}

// WatchNodes reads r.Err WITHOUT lock — RACE with concurrent writes to Err.
func (r *NodeRegistry) WatchNodes(_ context.Context, _ *metainternalversion.ListOptions) (interface{}, error) {
	if r.Err != nil { // RACE
		return nil, r.Err
	}
	return nil, nil
}
