package pkg

// Stripped reproduction of staging/src/k8s.io/cri-client/pkg/remote_image.go pre-PR #137949.
// BUG: useStreaming is a plain bool read/written across goroutines.
// FIX: the PR replaces it with atomic.Bool.

type remoteImageService struct {
	useStreaming bool
}

// ListImages reads useStreaming.
func (s *remoteImageService) ListImages() bool {
	return s.useStreaming
}

// streamImagesFallback writes useStreaming = false on streaming error.
func (s *remoteImageService) streamImagesFallback() {
	s.useStreaming = false
}
