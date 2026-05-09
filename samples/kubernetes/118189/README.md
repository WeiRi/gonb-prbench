# kubernetes-118189

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/118189 |
| Bug commit | `6d83e22ba48e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `topologycache.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00011e3f0 by goroutine 35:
  runtime.mapaccess2_faststr()
      /usr/local/go/src/runtime/map_faststr.go:108 +0x0
  ase/kubernetes-118189.(*ServiceSet).Has()
      /work/topologycache.go:18 +0x1cc
  ase/kubernetes-118189.(*TopologyCache).HasPopulatedHints()
      /work/topologycache.go:60 +0x16b
  ase/kubernetes-118189.TestRace_118189_TopologyCache.func2()
      /work/verified_test.go:24 +0xc8
  ase/kubernetes-118189.TestRace_118189_TopologyCache.func4()
      /work/verified_test.go:26 +0x41

Previous write at 0x00c00011e3f0 by goroutine 22:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-118189.(*ServiceSet).Insert()
      /work/topologycache.go:23 +0x2d7
  ase/kubernetes-118189.(*TopologyCache).SetHints()
      /work/topologycache.go:55 +0x288
  ase/kubernetes-118189.(*TopologyCache).AddHints()
      /work/topologycache.go:43 +0x9c
  ase/kubernetes-118189.TestRace_118189_TopologyCache.func1()
      /work/verified_test.go:18 +0x175
  ase/kubernetes-118189.TestRace_118189_TopologyCache.func3()
      /work/verified_test.go:20 +0x41

Goroutine 35 (running) created at:
  ase/kubernetes-118189.TestRace_118189_TopologyCache()
      /work/verified_test.go:21 +0x1ac
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-118189-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-118189-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-118189-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-118189-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-118189-bug .
# (then run as above, no --ssh flag)
```
