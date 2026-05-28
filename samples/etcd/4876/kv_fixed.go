package clientv3

import "sync"

type kvConn struct{ id int }

type kV struct {
	mu   sync.Mutex
	conn *kvConn
}

func newKV() *kV {
	return &kV{conn: &kvConn{id: 0}}
}

// Do — unchanged (was already correct).
func (k *kV) Do() *kvConn {
	k.mu.Lock()
	defer k.mu.Unlock()
	c := k.conn
	return c
}

// switchRemote — FIX (PR4876): acquire mutex BEFORE reading k.conn.
// Original BUG read k.conn before locking; fix moves Lock to the top.
func (k *kV) switchRemote() {
	k.mu.Lock()
	defer k.mu.Unlock()
	old := k.conn // safe under lock
	k.conn = &kvConn{id: old.id + 1}
}
