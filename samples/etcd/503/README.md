# etcd-503

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/503 |
| Bug commit | `3264b51a745a` |
| Category | channel_misuse |
| Oracle | RACE |
| Primary diff file | `mod/lock/v2/handler.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000188370 by goroutine 39:
  runtime.chansend()
      /usr/local/go/src/runtime/chan.go:160 +0x0
  ase/etcd-503.TestRace_503.func1.1()
      /work/verified_test.go:30 +0xbe

Previous write at 0x00c000188370 by goroutine 8:
  runtime.closechan()
      /usr/local/go/src/runtime/chan.go:357 +0x0
  ase/etcd-503.TestRace_503.func1()
      /work/verified_test.go:38 +0xce

Goroutine 39 (running) created at:
  ase/etcd-503.TestRace_503.func1()
      /work/verified_test.go:25 +0xac

Goroutine 8 (running) created at:
  ase/etcd-503.TestRace_503()
      /work/verified_test.go:15 +0x6f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
panic: send on closed channel

goroutine 53 [running]:
ase/etcd-503.TestRace_503.func1.1()
	/work/verified_test.go:30 +0xbf
created by ase/etcd-503.TestRace_503.func1 in goroutine 22
	/work/verified_test.go:25 +0xad
panic: send on closed channel

goroutine 50 [running]:
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-503-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-503-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-503-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-503-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-503-bug .
# (then run as above, no --ssh flag)
```
