package networkdb

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/memberlist"
)

// BUG: NetworkDB.SetPrimaryKey holds RLock; inside it calls
// memberlist.Keyring.UseKey which reads/writes k.keys with internal locking
// only on the write path (installKeys), not on the read iter. Concurrent
// SetPrimaryKey under RLock therefore races on k.keys.
// FIX (PR #49985): replace RLock with Lock so only one SetPrimaryKey runs at a time.
func TestRace_49985_setprimarykey(t *testing.T) {
	key1 := []byte("0123456789abcdef")
	key2 := []byte("fedcba9876543210")
	kr, err := memberlist.NewKeyring([][]byte{key1, key2}, key1)
	if err != nil {
		t.Fatal(err)
	}
	nDB := &NetworkDB{
		config:  &Config{Keys: [][]byte{key1, key2}},
		keyring: kr,
	}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			nDB.SetPrimaryKey(key1)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			nDB.SetPrimaryKey(key2)
		}
	}()
	wg.Wait()
}
