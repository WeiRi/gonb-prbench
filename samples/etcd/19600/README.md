# etcd-19600

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/19600 |
| Bug commit | `8b4c2cc11e6c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/storage/mvcc/watchable_store.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
=== RUN   TestCancelCloseRace
==================
WARNING: DATA RACE
Read at 0x0000007803e8 by goroutine 10:
  ase/etcd-19600.(*watchableStore).newWatcher()
      /work/watchable_store.go:44 +0x152
  ase/etcd-19600.TestCancelCloseRace.func1()
      /work/verified_test.go:24 +0xc1

Previous write at 0x0000007803e8 by goroutine 12:
  ase/etcd-19600.(*watchableStore).cancelWatcher()
      /work/watchable_store.go:69 +0x1bc
  ase/etcd-19600.TestCancelCloseRace.func1.(*watchableStore).newWatcher.1()
      /work/watchable_store.go:45 +0x38
  ase/etcd-19600.TestCancelCloseRace.func1()
      /work/verified_test.go:25 +0x1d2

Goroutine 10 (running) created at:
  ase/etcd-19600.TestCancelCloseRace()
      /work/verified_test.go:21 +0x13d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 12 (finished) created at:
  ase/etcd-19600.TestCancelCloseRace()
      /work/verified_test.go:21 +0x13d
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-19600-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-19600-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-19600-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-19600-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-19600-bug .
# (then run as above, no --ssh flag)
```
