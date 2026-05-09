# etcd-6324

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6324 |
| Bug commit | `48706a9cd678` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `proxy/grpcproxy/watch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000196540 by goroutine 8:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/etcd-6324.(*serverWatchStream).addDedicatedWatcher()
      /work/watch.go:15 +0x84
  ase/etcd-6324.TestRace_etcd_6324.func1()
      /work/verified_test.go:40 +0xa4

Previous read at 0x00c000196540 by goroutine 9:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/etcd-6324.(*serverWatchStream).close()
      /work/watch.go:21 +0xe4
  ase/etcd-6324.TestRace_etcd_6324.func2()
      /work/verified_test.go:45 +0x95

Goroutine 8 (running) created at:
  ase/etcd-6324.TestRace_etcd_6324()
      /work/verified_test.go:37 +0x1cc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/etcd-6324.TestRace_etcd_6324()
      /work/verified_test.go:43 +0x2a
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6324-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6324-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6324-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6324-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6324-bug .
# (then run as above, no --ssh flag)
```
