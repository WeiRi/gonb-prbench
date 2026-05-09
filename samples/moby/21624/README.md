# moby-21624

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/21624 |
| Bug commit | `e42c164763ed` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `volume/store/store.go:271` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d6540 by goroutine 10:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-21624.(*VolumeStore).create()
      /work/store.go:22 +0x1ed
  ase/moby-21624.TestRace_21624.func1()
      /work/verified_test.go:21 +0x109
  ase/moby-21624.TestRace_21624.gowrap1()
      /work/verified_test.go:23 +0x41

Previous write at 0x00c0000d6540 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-21624.(*VolumeStore).create()
      /work/store.go:22 +0x1ed
  ase/moby-21624.TestRace_21624.func1()
      /work/verified_test.go:21 +0x109
  ase/moby-21624.TestRace_21624.gowrap1()
      /work/verified_test.go:23 +0x41

Goroutine 10 (running) created at:
  ase/moby-21624.TestRace_21624()
      /work/verified_test.go:18 +0x15e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/moby-21624.TestRace_21624()
      /work/verified_test.go:18 +0x15e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-21624-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-21624-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-21624-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-21624-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-21624-bug .
# (then run as above, no --ssh flag)
```
