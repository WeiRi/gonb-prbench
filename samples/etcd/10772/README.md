# etcd-10772

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/10772 |
| Bug commit | `dc6885d73f71` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `mvcc/index.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000014228 by goroutine 23:
  ase/etcd-10772.(*treeIndex).Get()
      /work/index.go:33 +0x10a
  ase/etcd-10772.TestRace_10772.func2()
      /work/verified_test.go:31 +0x10b
  ase/etcd-10772.TestRace_10772.gowrap2()
      /work/verified_test.go:33 +0x41

Previous write at 0x00c000014228 by goroutine 12:
  ase/etcd-10772.(*treeIndex).Put()
      /work/index.go:24 +0xc4
  ase/etcd-10772.TestRace_10772.func1()
      /work/verified_test.go:23 +0xc5
  ase/etcd-10772.TestRace_10772.gowrap1()
      /work/verified_test.go:25 +0x41

Goroutine 23 (running) created at:
  ase/etcd-10772.TestRace_10772()
      /work/verified_test.go:28 +0x32e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 12 (finished) created at:
  ase/etcd-10772.TestRace_10772()
      /work/verified_test.go:20 +0x1d7
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-10772-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-10772-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-10772-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-10772-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-10772-bug .
# (then run as above, no --ssh flag)
```
