package accounts

import (
	"sync"
	"testing"
)

// TestRace_1503_account_manager_Sign_vs_expire: races Sign reading
// unlockedKey.PrivateKey after releasing am.mutex.RUnlock vs expire writing
// the same buffer under am.mutex.Lock.
func TestRace_1503_account_manager_Sign_vs_expire(t *testing.T) {
	am := NewManager()
	addr := Address{1}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			am.seedUnlocked(addr, []byte("12345678901234567890123456789012"))
			_ = am.Sign(addr)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			am.expire(addr)
		}
	}()
	wg.Wait()
}
