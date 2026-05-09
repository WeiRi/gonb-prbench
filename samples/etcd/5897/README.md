# etcd-5897

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/5897 |
| Bug commit | `dc2dced129be` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `etcdserver/api/v3rpc/watch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000124570 by goroutine 11:
  runtime.mapaccess1_fast64()
      /usr/local/go/src/runtime/map_fast64.go:13 +0x0
  ase/etcd-5897.(*serverWatchStream).sendLoop()
      /work/watch.go:29 +0x110
  ase/etcd-5897.TestRace_PR5897_ProgressPrevKVUnlocked.func2()
      /work/verified_test.go:27 +0xb8
  ase/etcd-5897.TestRace_PR5897_ProgressPrevKVUnlocked.gowrap2()
      /work/verified_test.go:29 +0x41

Previous write at 0x00c000124570 by goroutine 8:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/etcd-5897.(*serverWatchStream).recvLoop()
      /work/watch.go:23 +0x12b
  ase/etcd-5897.TestRace_PR5897_ProgressPrevKVUnlocked.func1()
      /work/verified_test.go:21 +0xcb
  ase/etcd-5897.TestRace_PR5897_ProgressPrevKVUnlocked.gowrap1()
      /work/verified_test.go:23 +0x41

Goroutine 11 (running) created at:
  ase/etcd-5897.TestRace_PR5897_ProgressPrevKVUnlocked()
      /work/verified_test.go:24 +0x15e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-5897.TestRace_PR5897_ProgressPrevKVUnlocked()
      /work/verified_test.go:18 +0x2a4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-5897-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-5897-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-5897-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-5897-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-5897-bug .
# (then run as above, no --ssh flag)
```
