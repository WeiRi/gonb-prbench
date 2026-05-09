# etcd-6197

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6197 |
| Bug commit | `28b797b538a5` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `integration/cluster_proxy.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d0450 by goroutine 15:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/etcd-6197.TestRace_PR6197_ToGRPC.func2()
      /work/verified_test.go:31 +0xcb

Previous read at 0x00c0000d0450 by goroutine 8:
  runtime.mapaccess2_fast64()
      /usr/local/go/src/runtime/map_fast64.go:53 +0x0
  ase/etcd-6197.toGRPC()
      /work/cluster_proxy.go:13 +0xc4
  ase/etcd-6197.TestRace_PR6197_ToGRPC.func1()
      /work/verified_test.go:24 +0x94

Goroutine 15 (running) created at:
  ase/etcd-6197.TestRace_PR6197_ToGRPC()
      /work/verified_test.go:27 +0xb1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-6197.TestRace_PR6197_ToGRPC()
      /work/verified_test.go:20 +0x17c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000d0450 by goroutine 13:
  runtime.mapassign_fast64()
```

(Full trace in `race_report_bug.txt`.)

## How to reproduce

### 1. SSH agent setup (one-time)
```bash
eval $(ssh-agent -a /tmp/ssh-agent-gonb.sock)
ssh-add ~/.ssh/id_ed25519
export SSH_AUTH_SOCK=/tmp/ssh-agent-gonb.sock
```

### 2. Build bug image
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6197-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6197-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6197-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6197-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6197-bug .
# (then run as above, no --ssh flag)
```
