# etcd-21375

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/21375 |
| Bug commit | `5e19d9a04` |
| Category | data_race |
| Oracle | PANIC |
| Primary diff file | `server/etcdserver/v3_server.go` |


## PANIC oracle excerpt

The following output is captured when running the bug build (`go test -race`):

```
--- FAIL: TestRequestCurrentIndex_LeaderChangedRace_21375 (0.00s)
    verified_test.go:36: PANIC oracle: 87/200 iters returned nil err instead of ErrLeaderChanged (v3_server.go)
goroutine dump (PANIC oracle fired):
goroutine 7 [running]:
ase/etcd-21375.TestRequestCurrentIndex_LeaderChangedRace_21375(0xc0000fa820)
	/work/verified_test.go:31 +0x1df

	at /work/v3_server.go:21 EtcdServer.requestCurrentIndex
FAIL	ase/etcd-21375	0.027s
```

`requestCurrentIndex`'s `select` picks between `readStateC` and `leaderChangedNotifier`
non-deterministically. When both are ready, the BUG state returns the read state
(`nil` err) instead of `ErrLeaderChanged` — a stale read after leadership change.
Fix re-checks `leaderChangedNotifier` in an inner `select` after `readStateC` fires.

## Reproduce

```bash
# BUG state (expect FAIL: PANIC oracle fires)
docker build -f bug.Dockerfile -t gonb-etcd-21375-bug .
docker run --rm --memory=2g --cpus=1 gonb-etcd-21375-bug

# FIX state (expect ok)
docker build -f fix.Dockerfile -t gonb-etcd-21375-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-21375-fix
```
