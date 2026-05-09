# kubernetes-96826

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/96826 |
| Bug commit | `7df379d75199` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/volume/plugins.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000122540 by goroutine 20:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-96826.(*VolumePluginMgr).AddPlugin()
      /work/plugins.go:19 +0x74
  ase/kubernetes-96826.TestRace_96826.func1()
      /work/verified_test.go:21 +0xb0

Previous read at 0x00c000122540 by goroutine 23:
  runtime.mapaccess2_faststr()
      /usr/local/go/src/runtime/map_faststr.go:108 +0x0
  ase/kubernetes-96826.(*VolumePluginMgr).GetPlugin()
      /work/plugins.go:25 +0xcf
  ase/kubernetes-96826.TestRace_96826.func2()
      /work/verified_test.go:27 +0xa6

Goroutine 20 (running) created at:
  ase/kubernetes-96826.TestRace_96826()
      /work/verified_test.go:18 +0x1bc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 23 (running) created at:
  ase/kubernetes-96826.TestRace_96826()
      /work/verified_test.go:24 +0xf2
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-96826-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-96826-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-96826-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-96826-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-96826-bug .
# (then run as above, no --ssh flag)
```
