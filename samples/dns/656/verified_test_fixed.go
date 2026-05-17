package dns

import "testing"

func TestRace_dns_656_conn_t(t *testing.T) { _ = &Conn{} }
