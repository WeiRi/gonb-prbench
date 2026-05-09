# kubernetes-103487

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/103487 |
| Bug commit | `b289fbb03dee` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `fixture.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000ca450 by goroutine 11:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-103487.(*Watcher).Modify()
      /work/fixture.go:30 +0x351
  ase/kubernetes-103487.(*tracker).add()
      /work/fixture.go:56 +0x325
  ase/kubernetes-103487.TestRace_103487_FixtureWatcher.func2()
      /work/verified_test.go:26 +0x412

Previous read at 0x00c0000ca450 by goroutine 10:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/kubernetes-103487.TestRace_103487_FixtureWatcher.func1()
      /work/verified_test.go:20 +0xe6

Goroutine 11 (running) created at:
  ase/kubernetes-103487.TestRace_103487_FixtureWatcher()
      /work/verified_test.go:23 +0x67
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 10 (finished) created at:
  ase/kubernetes-103487.TestRace_103487_FixtureWatcher()
      /work/verified_test.go:17 +0x496
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
fatal error: concurrent map read and map write
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-103487-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-103487-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-103487-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-103487-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-103487-bug .
# (then run as above, no --ssh flag)
```
