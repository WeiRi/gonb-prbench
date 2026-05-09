# etcd-6587

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6587 |
| Bug commit | `98897b760345` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `clientv3/watch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00019c128 by goroutine 10:
  ase/etcd-6587.(*watcherStream).resume()
      /work/watch.go:33 +0xa7
  ase/etcd-6587.TestRace_PR6587_InitReqRevRace.func3()
      /work/verified_test.go:31 +0x9f

Previous write at 0x00c00019c128 by goroutine 8:
  ase/etcd-6587.(*watcherStream).serveSubstream()
      /work/watch.go:27 +0xbb
  ase/etcd-6587.TestRace_PR6587_InitReqRevRace.func1()
      /work/verified_test.go:19 +0x84

Goroutine 10 (running) created at:
  ase/etcd-6587.TestRace_PR6587_InitReqRevRace()
      /work/verified_test.go:28 +0x64
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/etcd-6587.TestRace_PR6587_InitReqRevRace()
      /work/verified_test.go:17 +0x224
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR6587_InitReqRevRace (0.00s)
FAIL
FAIL	ase/etcd-6587	0.016s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6587-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6587-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6587-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6587-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6587-bug .
# (then run as above, no --ssh flag)
```
