// Minimal extraction of pre-fix go-ethereum accounts/account_manager.go for
// PR #1503. Pre-fix race: Sign() does am.mutex.RUnlock() *before* reading
// unlockedKey.PrivateKey, allowing a concurrent timed-lock goroutine to
// zero out the key under am.mutex.Lock(). PR fix: hold am.mutex.RLock for
// the full duration of Sign() via `defer am.mutex.RUnlock()`.
// Production-code path: accounts/account_manager.go
package accounts

import "sync"

type Address [4]byte

type unlocked struct {
	PrivateKey []byte
}

type Manager struct {
	mutex    sync.RWMutex
	unlocked map[Address]*unlocked
}

func NewManager() *Manager {
	return &Manager{unlocked: make(map[Address]*unlocked)}
}

func (am *Manager) seedUnlocked(addr Address, key []byte) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.unlocked[addr] = &unlocked{PrivateKey: key}
}

// expire mimics the timeout goroutine that zeroes the key + removes the entry
// while holding am.mutex.Lock.
func (am *Manager) expire(addr Address) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	if u, ok := am.unlocked[addr]; ok {
		// Pre-fix code path: write into u.PrivateKey under am.mutex.
		for i := range u.PrivateKey {
			u.PrivateKey[i] = 0
		}
		delete(am.unlocked, addr)
	}
}

// Sign — pre-fix version: releases the read lock too early.
// Upstream path: accounts/account_manager.go (pre-fix line ~80-87).
func (am *Manager) Sign(addr Address) []byte {
	am.mutex.RLock()
	unlockedKey, found := am.unlocked[addr]
	am.mutex.RUnlock() // pre-fix: released before key is used.
	if !found {
		return nil
	}
	// Read unlockedKey.PrivateKey OUTSIDE the lock — races with expire()'s zero/delete.
	out := make([]byte, len(unlockedKey.PrivateKey))
	copy(out, unlockedKey.PrivateKey)
	return out
}
