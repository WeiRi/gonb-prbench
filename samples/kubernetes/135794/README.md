# kubernetes-135794

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/135794 |
| Bug commit | `dd838ccf07a5` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/volumemanager/volume_manager.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000012018 by goroutine 32:
  ase/kubernetes-135794.TestRace_135794.func1.1()
      /work/verified_test.go:28 +0x104

Previous write at 0x00c000012018 by goroutine 34:
  ase/kubernetes-135794.TestRace_135794.func1.2()
      /work/verified_test.go:32 +0x1cc

Goroutine 32 (running) created at:
  ase/kubernetes-135794.TestRace_135794.func1()
      /work/verified_test.go:26 +0x22a

Goroutine 34 (finished) created at:
  ase/kubernetes-135794.TestRace_135794.func1()
      /work/verified_test.go:30 +0xb6
==================
==================
WARNING: DATA RACE
Read at 0x00c00001e030 by goroutine 32:
  runtime.growslice()
      /usr/local/go/src/runtime/slice.go:157 +0x0
  ase/kubernetes-135794.TestRace_135794.func1.1()
      /work/verified_test.go:28 +0x144

Previous write at 0x00c00001e030 by goroutine 34:
  ase/kubernetes-135794.TestRace_135794.func1.2()
      /work/verified_test.go:32 +0x176

Goroutine 32 (running) created at:
  ase/kubernetes-135794.TestRace_135794.func1()
      /work/verified_test.go:26 +0x22a

Goroutine 34 (finished) created at:
  ase/kubernetes-135794.TestRace_135794.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-135794-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-135794-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-135794-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-135794-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-135794-bug .
# (then run as above, no --ssh flag)
```
