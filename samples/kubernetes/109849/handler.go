// Pre-fix handler.go from PR #109849: disconnectClient reads s.clients[name]
// WITHOUT holding s.mutex (line 86), racing with deregisterClient
// (line 103) which writes the map under lock. Production-code stand-in so
// race-detector frames hit handler.go (PR diff path).
package v1beta1

import "sync"

// Client is the plugin client interface.
type Client interface {
	Run()
	Disconnect() error
	Name() string
}

type mockClient struct {
	name string
}

func (m *mockClient) Run()             {}
func (m *mockClient) Disconnect() error { return nil }
func (m *mockClient) Name() string      { return m.name }

type server struct {
	mutex   sync.Mutex
	clients map[string]Client
}

func newServer() *server {
	return &server{clients: make(map[string]Client)}
}

// registerClient: writes s.clients[name] under lock (handler.go pre-fix path).
func (s *server) registerClient(name string, c Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[name] = c
}

// deregisterClient: handler.go:103 in pre-fix — writes (deletes) s.clients[name] under lock.
func (s *server) deregisterClient(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.clients, name)
}

// disconnectClient: handler.go:86 in pre-fix — UNPROTECTED read of s.clients[name].
func (s *server) disconnectClient(name string) error {
	c := s.clients[name] // line 86 — racy map read (no lock)
	s.deregisterClient(name)
	if c != nil {
		return c.Disconnect()
	}
	return nil
}

// DeRegisterPlugin: handler.go pre-fix — also reads s.clients while locking,
// but then dispatches disconnectClient which races. (Locked variant kept here
// for completeness; the racy site is in disconnectClient.)
func (s *server) DeRegisterPlugin(name string) {
	s.mutex.Lock()
	_, exists := s.clients[name]
	s.mutex.Unlock()
	if exists {
		_ = s.disconnectClient(name)
	}
}
