# etcd-13203

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/13203 |
| Bug commit | `53d234f1fe2b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client/v3/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d8518 by goroutine 9:
  ase/etcd-13203.(*Client).SetEndpoints()
      /work/client.go:21 +0xc6
  ase/etcd-13203.TestRace_PR13203_Endpoints.func2()
      /work/verified_test.go:32 +0xe4

Previous read at 0x00c0000d8518 by goroutine 8:
  ase/etcd-13203.TestRace_PR13203_Endpoints.func1()
      /work/verified_test.go:25 +0xa4

Goroutine 9 (running) created at:
  ase/etcd-13203.TestRace_PR13203_Endpoints()
      /work/verified_test.go:28 +0x1d1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-13203.TestRace_PR13203_Endpoints()
      /work/verified_test.go:21 +0x29c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR13203_Endpoints (0.01s)
FAIL
FAIL	ase/etcd-13203	0.026s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-13203-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-13203-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-13203-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-13203-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-13203-bug .
# (then run as above, no --ssh flag)
```
