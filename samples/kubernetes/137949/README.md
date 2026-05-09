# kubernetes-137949

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/137949 |
| Bug commit | `ef247770b50e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/cri-client/pkg/remote_image.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00001428f by goroutine 22:
  ase/kubernetes-137949.(*remoteImageService).ListImages()
      /work/remote_image.go:13 +0xa4
  ase/kubernetes-137949.TestRace_137949_UseStreamingBoolRace.func1.1()
      /work/verified_test.go:30 +0x9b

Previous write at 0x00c00001428f by goroutine 25:
  ase/kubernetes-137949.(*remoteImageService).streamImagesFallback()
      /work/remote_image.go:18 +0xa4
  ase/kubernetes-137949.TestRace_137949_UseStreamingBoolRace.func1.2()
      /work/verified_test.go:38 +0x9b

Goroutine 22 (running) created at:
  ase/kubernetes-137949.TestRace_137949_UseStreamingBoolRace.func1()
      /work/verified_test.go:27 +0x17c

Goroutine 25 (finished) created at:
  ase/kubernetes-137949.TestRace_137949_UseStreamingBoolRace.func1()
      /work/verified_test.go:35 +0x224
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_137949_UseStreamingBoolRace (0.00s)
FAIL
FAIL	ase/kubernetes-137949	0.021s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-137949-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-137949-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-137949-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-137949-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-137949-bug .
# (then run as above, no --ssh flag)
```
