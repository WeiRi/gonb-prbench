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

// Do — BUG (pre-PR4876): reads kv.conn under lock, but switchRemote also reads
// it WITHOUT holding the lock first (line 35-55 region).
func (k *kV) Do() *kvConn {
	k.mu.Lock()
	defer k.mu.Unlock()
	c := k.conn
	return c
}

// switchRemote — BUG: reads k.conn before locking (line 41-42, racy).
func (k *kV) switchRemote() {
	old := k.conn // line 41 BUG: unlocked read
	_ = old       // line 42
	k.mu.Lock()
	k.conn = &kvConn{id: old.id + 1} // line 50 write under lock
	k.mu.Unlock()
}
