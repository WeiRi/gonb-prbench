# etcd-6279

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6279 |
| Bug commit | `c388b2f22f12` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `mvcc/kvstore.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000012148 by goroutine 9:
  ase/etcd-6279.(*store).Hash()
      /work/kvstore.go:34 +0xb0
  ase/etcd-6279.TestRace_PR6279_ForceCommitOutsideLock.func2()
      /work/verified_test.go:27 +0xa8

Previous write at 0x00c000012148 by goroutine 8:
  ase/etcd-6279.(*store).Compact()
      /work/kvstore.go:27 +0x44
  ase/etcd-6279.TestRace_PR6279_ForceCommitOutsideLock.func1()
      /work/verified_test.go:21 +0xa9
  ase/etcd-6279.TestRace_PR6279_ForceCommitOutsideLock.gowrap1()
      /work/verified_test.go:23 +0x41

Goroutine 9 (running) created at:
  ase/etcd-6279.TestRace_PR6279_ForceCommitOutsideLock()
      /work/verified_test.go:24 +0x11d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/etcd-6279.TestRace_PR6279_ForceCommitOutsideLock()
      /work/verified_test.go:18 +0x267
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR6279_ForceCommitOutsideLock (0.05s)
FAIL
FAIL	ase/etcd-6279	0.065s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6279-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6279-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6279-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6279-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6279-bug .
# (then run as above, no --ssh flag)
```
