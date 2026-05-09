package main

// connTrackKey identifies a UDP connection by client address. Mirrors
// cmd/docker-proxy/udp_proxy_linux.go from moby PR #49649.
// The original implementation accessed the connTrackTable map without a
// mutex, racing between the read on the receive path and the write on the
// connection-add path.

type connTrackTable struct {
	conns map[string]int
}

func newConnTrackTable() *connTrackTable {
	return &connTrackTable{conns: make(map[string]int)}
}

// addConn records a new client backend index. Non-synchronized: races with
// lookupConn on the receive path.
func (t *connTrackTable) addConn(key string, fd int) {
	t.conns[key] = fd // RACE write
}

// lookupConn returns the backend fd registered for key. Non-synchronized:
// races with addConn.
func (t *connTrackTable) lookupConn(key string) int {
	return t.conns[key] // RACE read
}
