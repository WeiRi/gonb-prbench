# kubernetes-104604

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/104604 |
| Bug commit | `3e10db97d07e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/util/manager/watch_based_manager.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00012ca50 by goroutine 21:
  ase/kubernetes-104604.(*objectCacheItem).stopThreadUnsafe()
      /work/watch_based_manager.go:43 +0x1e9
  ase/kubernetes-104604.(*objectCacheItem).stopIfIdle()
      /work/watch_based_manager.go:53 +0x219
  ase/kubernetes-104604.TestRace_104604.func2()
      /work/verified_test.go:28 +0xd4

Previous read at 0x00c00012ca50 by goroutine 20:
  ase/kubernetes-104604.(*objectCacheItem).Get()
      /work/watch_based_manager.go:64 +0x2f
  ase/kubernetes-104604.TestRace_104604.func1()
      /work/verified_test.go:22 +0xa4

Goroutine 21 (running) created at:
  ase/kubernetes-104604.TestRace_104604()
      /work/verified_test.go:25 +0x56
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 20 (running) created at:
  ase/kubernetes-104604.TestRace_104604()
      /work/verified_test.go:19 +0x2bc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
--- FAIL: TestRace_104604 (0.01s)
    testing.go:1465: race detected during execution of test
FAIL
FAIL	ase/kubernetes-104604	0.028s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-104604-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-104604-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-104604-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-104604-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-104604-bug .
# (then run as above, no --ssh flag)
```
