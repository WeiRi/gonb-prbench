# etcd-7361

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/7361 |
| Bug commit | `0c0fbbd7c500` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `proxy/tcpproxy/userspace.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000096008 by goroutine 27:
  ase/etcd-7361.(*TCPProxy).runMonitorOnce.func1()
      /work/userspace.go:25 +0x2e

Previous write at 0x00c000096008 by goroutine 9:
  ase/etcd-7361.(*TCPProxy).runMonitorOnce()
      /work/userspace.go:22 +0xb8
  ase/etcd-7361.TestRace_TCPProxy__LoopVarCapture.func1()
      /work/verified_test.go:67 +0x99

Goroutine 27 (running) created at:
  ase/etcd-7361.(*TCPProxy).runMonitorOnce()
      /work/userspace.go:24 +0x184
  ase/etcd-7361.TestRace_TCPProxy__LoopVarCapture.func1()
      /work/verified_test.go:67 +0x99

Goroutine 9 (running) created at:
  ase/etcd-7361.TestRace_TCPProxy__LoopVarCapture()
      /work/verified_test.go:64 +0x264
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_TCPProxy__LoopVarCapture (0.67s)
FAIL
FAIL	ase/etcd-7361	0.693s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-7361-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-7361-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-7361-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-7361-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-7361-bug .
# (then run as above, no --ssh flag)
```
