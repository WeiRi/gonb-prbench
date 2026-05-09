# kubernetes-128495

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/128495 |
| Bug commit | `dfd456c56741` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/volume/plugins.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
=== RUN   TestVolumePluginConcurrentRace
==================
WARNING: DATA RACE
Write at 0x00c0000d8570 by goroutine 9:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-128495.(*VolumePluginMgr).refreshProbedPlugins()
      /work/plugins.go:170 +0x1a6
  ase/kubernetes-128495.(*VolumePluginMgr).FindPluginByName()
      /work/plugins.go:123 +0x5a
  ase/kubernetes-128495.TestVolumePluginConcurrentRace.func1()
      /work/verified_test.go:28 +0xa5

Previous write at 0x00c0000d8570 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-128495.(*VolumePluginMgr).refreshProbedPlugins()
      /work/plugins.go:170 +0x1a6
  ase/kubernetes-128495.(*VolumePluginMgr).FindPluginByName()
      /work/plugins.go:123 +0x5a
  ase/kubernetes-128495.TestVolumePluginConcurrentRace.func1()
      /work/verified_test.go:28 +0xa5

Goroutine 9 (running) created at:
  ase/kubernetes-128495.TestVolumePluginConcurrentRace()
      /work/verified_test.go:25 +0x27c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/kubernetes-128495.TestVolumePluginConcurrentRace()
      /work/verified_test.go:25 +0x27c
  testing.tRunner()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-128495-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-128495-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-128495-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-128495-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-128495-bug .
# (then run as above, no --ssh flag)
```
