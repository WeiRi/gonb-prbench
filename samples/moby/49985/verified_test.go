// Race-trigger test for moby-49985; see README.md for usage.

package networkdb

import (
	"sync"
	"testing"
)

type _bugKeyring struct {
	primaryKey []byte
	keys       [][]byte
}

func (k *_bugKeyring) UseKey(key []byte) {
	for _, kk := range k.keys {
		if len(kk) == len(key) {
			k.primaryKey = kk // BUG: write while concurrent reader at PrimaryKey()
			return
		}
	}
}

func (k *_bugKeyring) PrimaryKey() []byte { return k.primaryKey }

type _BugNetworkDB struct {
	sync.RWMutex
	keys    [][]byte
	keyring *_bugKeyring
}

// SetPrimaryKey (BUG): RLock instead of Lock — multiple writers allowed.
func (nDB *_BugNetworkDB) SetPrimaryKey(key []byte) {
	nDB.RLock()
	defer nDB.RUnlock()
	for _, dbKey := range nDB.keys {
		if len(key) == len(dbKey) {
			if nDB.keyring != nil {
				nDB.keyring.UseKey(key)
			}
			break
		}
	}
}

func (nDB *_BugNetworkDB) readPrimary() []byte {
	nDB.RLock()
	defer nDB.RUnlock()
	if nDB.keyring != nil {
		return nDB.keyring.PrimaryKey()
	}
	return nil
}

func TestRace_PR49985_SetPrimaryKey(t *testing.T) {
	nDB := &_BugNetworkDB{keyring: &_bugKeyring{}}
	for i := 0; i < 8; i++ {
		k := []byte{byte(i), byte(i), byte(i), byte(i)}
		nDB.keys = append(nDB.keys, k)
		nDB.keyring.keys = append(nDB.keyring.keys, k)
	}
	nDB.keyring.primaryKey = nDB.keys[0]

	keyA := []byte{1, 1, 1, 1}
	keyB := []byte{2, 2, 2, 2}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				if (i+j)%2 == 0 {
					nDB.SetPrimaryKey(keyA)
				} else {
					nDB.SetPrimaryKey(keyB)
				}
				_ = nDB.readPrimary()
			}
		}(i)
	}
	wg.Wait()
}
