# kubernetes-109849

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/109849 |
| Bug commit | `af4dceeac2d9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/cm/devicemanager/plugin/v1beta1/handler.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000ca3f0 by goroutine 21:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/kubernetes-109849.(*server).disconnectClient()
      /work/handler.go:49 +0x5c
  ase/kubernetes-109849.TestRace_109849_DisconnectMapRace.func1.2()
      /work/verified_test.go:48 +0xbb

Previous write at 0x00c0000ca3f0 by goroutine 20:
  runtime.mapdelete_faststr()
      /usr/local/go/src/runtime/map_faststr.go:301 +0x0
  ase/kubernetes-109849.(*server).deregisterClient()
      /work/handler.go:44 +0xcc
  ase/kubernetes-109849.TestRace_109849_DisconnectMapRace.func1.1()
      /work/verified_test.go:40 +0xe5

Goroutine 21 (running) created at:
  ase/kubernetes-109849.TestRace_109849_DisconnectMapRace.func1()
      /work/verified_test.go:45 +0x355

Goroutine 20 (finished) created at:
  ase/kubernetes-109849.TestRace_109849_DisconnectMapRace.func1()
      /work/verified_test.go:35 +0x28a
==================
==================
WARNING: DATA RACE
Read at 0x00c0000ca420 by goroutine 26:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/kubernetes-109849.(*server).disconnectClient()
      /work/handler.go:49 +0x5c
  ase/kubernetes-109849.TestRace_109849_DisconnectMapRace.func1.2()
      /work/verified_test.go:48 +0xbb
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-109849-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-109849-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-109849-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-109849-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-109849-bug .
# (then run as above, no --ssh flag)
```
