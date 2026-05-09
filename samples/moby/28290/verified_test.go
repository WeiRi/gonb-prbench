// Race-trigger test for moby-28290; see README.md for usage.

package store

import (
	"sync"
	"testing"
	"time"
)

// VolumeStore mirrors moby/volume/store/store.go BUG state pre-PR #28290.
// globalLock is a plain Mutex; list() reads bypass the lock.
type _BugVolumeStore struct {
	globalLock sync.Mutex
	names      map[string]string
	labels     map[string]map[string]string
	options    map[string]map[string]string
}

// list (BUG): reads labels[name] and options[name] WITHOUT holding globalLock.
func (s *_BugVolumeStore) list(name string) (map[string]string, map[string]string) {
	return s.labels[name], s.options[name]
}

func (s *_BugVolumeStore) addVolume(name string) {
	s.globalLock.Lock()
	s.names[name] = name
	s.labels[name] = map[string]string{"k": "v"}
	s.options[name] = map[string]string{"o": "p"}
	s.globalLock.Unlock()
}

func (s *_BugVolumeStore) removeVolume(name string) {
	s.globalLock.Lock()
	delete(s.names, name)
	delete(s.labels, name)
	delete(s.options, name)
	s.globalLock.Unlock()
}

func TestRace_PR28290_VolumeStoreMapAccess(t *testing.T) {
	s := &_BugVolumeStore{
		names:   make(map[string]string),
		labels:  make(map[string]map[string]string),
		options: make(map[string]map[string]string),
	}
	for i := 0; i < 32; i++ {
		s.addVolume("v")
	}
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				s.removeVolume("v")
				s.addVolume("v")
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				_, _ = s.list("v")
			}
		}()
	}
	time.Sleep(150 * time.Millisecond)
	close(stop)
	wg.Wait()
}
