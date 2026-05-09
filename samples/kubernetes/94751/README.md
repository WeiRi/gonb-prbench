# kubernetes-94751

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/94751 |
| Bug commit | `66334f02e8c5` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/time_cache.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d6480 by goroutine 15:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-94751.(*timeCache).Add()
      /work/time_cache.go:20 +0x12c
  ase/kubernetes-94751.TestRace_94751.func2()
      /work/verified_test.go:30 +0xe0

Previous read at 0x00c0000d6480 by goroutine 8:
  runtime.mapaccess2_faststr()
      /usr/local/go/src/runtime/map_faststr.go:108 +0x0
  ase/kubernetes-94751.(*timeCache).Get()
      /work/time_cache.go:25 +0xf7
  ase/kubernetes-94751.TestRace_94751.func1()
      /work/verified_test.go:24 +0xd1

Goroutine 15 (running) created at:
  ase/kubernetes-94751.TestRace_94751()
      /work/verified_test.go:27 +0x164
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1576 +0x216
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1629 +0x47

Goroutine 8 (finished) created at:
  ase/kubernetes-94751.TestRace_94751()
      /work/verified_test.go:21 +0x273
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1576 +0x216
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1629 +0x47
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-94751-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-94751-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-94751-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-94751-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-94751-bug .
# (then run as above, no --ssh flag)
```
