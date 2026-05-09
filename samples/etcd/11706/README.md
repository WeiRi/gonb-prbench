# etcd-11706

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/11706 |
| Bug commit | `bf883bd15b73` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `clientv3/ctx.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00019e540 by goroutine 15:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/etcd-11706.WithRequireLeader()
      /work/ctx.go:13 +0xa4
  ase/etcd-11706.TestRace_11706.func2()
      /work/verified_test.go:26 +0x69

Previous write at 0x00c00019e540 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/etcd-11706.MDLocal.Set()
      /work/metadata_stub.go:17 +0xaf
  ase/etcd-11706.TestRace_11706.func1()
      /work/verified_test.go:19 +0x51

Goroutine 15 (running) created at:
  ase/etcd-11706.TestRace_11706()
      /work/verified_test.go:23 +0x26a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-11706.TestRace_11706()
      /work/verified_test.go:17 +0x336
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-11706-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-11706-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-11706-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-11706-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-11706-bug .
# (then run as above, no --ssh flag)
```
