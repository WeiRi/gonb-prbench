package clientcmd

import (
	"sync"
	"testing"

	restclient "k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// mockLoader implements ClientConfigLoader for testing.
// Load() returns a minimal valid config.
type mockLoader92139 struct{}

func (m *mockLoader92139) Load() (*clientcmdapi.Config, error) {
	return clientcmdapi.NewConfig(), nil
}

func (m *mockLoader92139) IsDefaultConfig(*restclient.Config) bool {
	return false
}

func (m *mockLoader92139) GetLoadingPrecedence() []string {
	return nil
}

func (m *mockLoader92139) GetStartingConfig() (*clientcmdapi.Config, error) {
	return clientcmdapi.NewConfig(), nil
}

func (m *mockLoader92139) GetDefaultFilename() string {
	return ""
}

func (m *mockLoader92139) IsExplicitFile() bool {
	return false
}

func (m *mockLoader92139) GetExplicitFile() string {
	return ""
}

func TestRace_92139(t *testing.T) {
	loader := &mockLoader92139{}
	config := &DeferredLoadingClientConfig{
		loader:    loader,
		overrides: &ConfigOverrides{},
	}

	const N = 30
	const ITERS = 200

	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// Reset clientConfig to nil to force the racy read path.
				// This write to clientConfig outside the lock races with
				// the read at line 63 and write at line 80 in createClientConfig().
				config.clientConfig = nil
				_, _ = config.createClientConfig()
			}
		}()
	}

	wg.Wait()
}
