# kubernetes-123396

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/123396 |
| Bug commit | `7606cf7b3d78` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pathrecorder.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000214100 by goroutine 23:
  runtime.slicecopy()
      /usr/local/go/src/runtime/slice.go:310 +0x0
  ase/kubernetes-123396.(*PathRecorderMux).ListedPaths()
      /work/pathrecorder.go:26 +0xcc
  ase/kubernetes-123396.TestRace_123396_PathRecorder.func2()
      /work/verified_test.go:24 +0xee

Previous write at 0x00c000214100 by goroutine 12:
  ase/kubernetes-123396.(*PathRecorderMux).Register()
      /work/pathrecorder.go:21 +0x13a
  ase/kubernetes-123396.TestRace_123396_PathRecorder.func1()
      /work/verified_test.go:18 +0x149
  ase/kubernetes-123396.TestRace_123396_PathRecorder.func3()
      /work/verified_test.go:20 +0x41

Goroutine 23 (running) created at:
  ase/kubernetes-123396.TestRace_123396_PathRecorder()
      /work/verified_test.go:21 +0xa4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 12 (running) created at:
  ase/kubernetes-123396.TestRace_123396_PathRecorder()
      /work/verified_test.go:15 +0x1fc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-123396-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-123396-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-123396-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-123396-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-123396-bug .
# (then run as above, no --ssh flag)
```
