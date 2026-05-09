// PR miekg/dns#780 - server.go - data race / order violation between
// readTCP/readUDP and ShutdownContext on Server.started + connection
// read deadline. Pre-fix (and pre-#774): readTCP/readUDP read srv.started
// without srv.lock and call SetReadDeadline, while ShutdownContext sets
// srv.started=false and SetReadDeadline(aLongTimeAgo) under srv.lock.
// PR fix #780 takes srv.lock.RLock around srv.started check + SetReadDeadline
// to make the operation atomic w.r.t. shutdown.
// Production-code path: server.go (pre-fix line ~668, ~711, ~417).
package dns

import (
	"sync"
	"time"
)

type Server struct {
	lock     sync.RWMutex
	started  bool
	deadline time.Time
}

func NewServer() *Server { return &Server{started: true} }

// ReadTCP — pre-fix: reads srv.started WITHOUT srv.lock and sets deadline.
// Upstream: server.go (pre-fix line ~668).
func (srv *Server) ReadTCP() {
	if srv.started { // <- racy read
		srv.deadline = time.Now().Add(time.Hour) // <- not atomic with ShutdownContext
	}
}

// ShutdownContext — sets started=false and aLongTimeAgo deadline under srv.lock.
// Pre-fix race: between ReadTCP's check of srv.started and its SetReadDeadline,
// ShutdownContext can flip started=false and set aLongTimeAgo. Then ReadTCP
// overrides the aLongTimeAgo with a future deadline, blocking shutdown.
// Upstream: server.go (pre-fix line ~417).
func (srv *Server) ShutdownContext() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	srv.started = false
	srv.deadline = time.Unix(1, 0) // aLongTimeAgo
}
