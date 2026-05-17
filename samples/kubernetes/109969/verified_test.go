// Whitebox PoC for kubernetes-109969: data race on Response.User when
// AuthenticatedGroupAdder modifies a shared user.Info backing array via append.
// Production code in authenticated_group_adder.go.
package group

import (
	"net/http"
	"sync"
	"testing"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
)

func TestRace_109969_GroupAdderSharedBackingArray(t *testing.T) {
	const numGoroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				capacity := make([]string, 0, 16)
				capacity = append(capacity, "user1", "user2")

				response := &authenticator.Response{
					User: &user.DefaultInfo{
						Name:   "a",
						Groups: capacity,
					},
				}

				ag := &AuthenticatedGroupAdder{
					Authenticator: authenticator.RequestFunc(func(req *http.Request) (*authenticator.Response, bool, error) {
						return response, true, nil
					}),
				}

				req1, _ := http.NewRequest("GET", "/", nil)
				req2, _ := http.NewRequest("GET", "/", nil)

				var innerWg sync.WaitGroup
				innerWg.Add(2)
				go func() {
					defer innerWg.Done()
					ag.AuthenticateRequest(req1)
				}()
				go func() {
					defer innerWg.Done()
					ag.AuthenticateRequest(req2)
				}()
				innerWg.Wait()
			}
		}()
	}
	wg.Wait()
}
