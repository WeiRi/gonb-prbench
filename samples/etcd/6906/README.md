# etcd-6906

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6906 |
| Bug commit | `a076510cc192` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `proxy/grpcproxy/watch_broadcasts.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00011e540 by goroutine 12:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/etcd-6906.TestRace_6906.func1()
      /work/verified_test.go:29 +0xfe

Previous read at 0x00c00011e540 by goroutine 21:
  ase/etcd-6906.(*watchBroadcasts).empty()
      /work/watch_broadcasts.go:21 +0xcc
  ase/etcd-6906.TestRace_6906.func2()
      /work/verified_test.go:42 +0xae

Goroutine 12 (running) created at:
  ase/etcd-6906.TestRace_6906()
      /work/verified_test.go:24 +0x204
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 21 (finished) created at:
  ase/etcd-6906.TestRace_6906()
      /work/verified_test.go:39 +0x2d1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00011e540 by goroutine 12:
  runtime.mapdelete_fast64()
      /usr/local/go/src/runtime/map_fast64.go:273 +0x0
  ase/etcd-6906.TestRace_6906.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6906-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6906-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6906-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6906-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6906-bug .
# (then run as above, no --ssh flag)
```
