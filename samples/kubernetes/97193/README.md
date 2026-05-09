# kubernetes-97193

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/97193 |
| Bug commit | `fb02a59a6a10` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/pluginmanager/reconciler/reconciler.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4540 by goroutine 10:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-97193.(*reconciler).AddHandler()
      /work/reconciler.go:16 +0x74
  ase/kubernetes-97193.TestRace_97193.func1()
      /work/verified_test.go:23 +0x14f
  ase/kubernetes-97193.TestRace_97193.gowrap1()
      /work/verified_test.go:25 +0x41

Previous read at 0x00c0000d4540 by goroutine 47:
  runtime.mapiternext()
      /usr/local/go/src/runtime/map.go:862 +0x0
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:859 +0x22e

Goroutine 10 (running) created at:
  ase/kubernetes-97193.TestRace_97193()
      /work/verified_test.go:19 +0x10d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 47 (running) created at:
  ase/kubernetes-97193.TestRace_97193()
      /work/verified_test.go:29 +0x251
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-97193-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-97193-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-97193-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-97193-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-97193-bug .
# (then run as above, no --ssh flag)
```
