// Race test for kubernetes-132061
// Targets the race in ResolverEnvOption: BUG state creates ONE shared
// ResolverTypeProvider that is written-to (tp.underlyingTypeProvider =)
// inside the envOpt closure, every time envOpt is invoked.
// FIX state creates a NEW tp per invocation, no shared state.
package common

import (
	"sync"
	"testing"

	"github.com/google/cel-go/cel"
)

// minimalTypeResolver satisfies the TypeResolver interface with no-op behavior.
type minimalTypeResolver struct{}

func (m *minimalTypeResolver) Resolve(name string) (ResolvedType, bool) {
	return nil, false
}

func TestRace_132061_ResolverEnvOption(t *testing.T) {
	envOpt := ResolverEnvOption(&minimalTypeResolver{})

	const N = 80
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Invoke envOpt via cel.NewEnv. In BUG state this writes
			// the shared tp.underlyingTypeProvider from many goroutines
			// at typeprovider.go:117 (inside the closure).
			env, _ := cel.NewEnv(envOpt)
			_ = env
		}()
	}
	wg.Wait()
}
