package integration

import (
	"ase/etcd-6197/clientv3"
)

type grpcAPI struct{}

// BUG (pre-PR6197): package-level map without lock; concurrent writes/reads race.
var proxies = make(map[*clientv3.Client]grpcAPI)

func toGRPC(c *clientv3.Client) grpcAPI {
	if v, ok := proxies[c]; ok { // BUG line 28: read map without lock
		return v
	}
	return grpcAPI{}
}
