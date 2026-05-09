package common

import (
	"sync"
	"testing"

	"github.com/google/cel-go/cel"

	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/apiserver/pkg/cel/environment"
)

func TestRace_132061(t *testing.T) {
	resolver := &mockTypeResolver{}

	// Buggy version: NewResolverTypeProviderAndEnvOption creates ONE shared tp
	// that is written to (tp.underlyingTypeProvider = ...) from EVERY call to envOption.
	// When multiple goroutines call envOption concurrently, the shared tp has a data race.
	tp, envOption := NewResolverTypeProviderAndEnvOption(resolver)

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Create a new env using the envOption, which writes to shared tp.underlyingTypeProvider
				envSet, err := environment.MustBaseEnvSet(environment.DefaultCompatibilityVersion(), true).
					Extend(environment.VersionedOptions{
						IntroducedVersion: version.MajorMinor(1, 30),
						EnvOptions:        []cel.EnvOption{envOption},
					})
				if err != nil {
					continue
				}
				env, err := envSet.Env(environment.StoredExpressions)
				if err != nil {
					continue
				}
				// Also read through tp concurrently to race with writes to underlyingTypeProvider
				_, _ = tp.FindStructType("Test")
				_ = env
			}
		}()
	}

	wg.Wait()
}
