# moby-22279

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/22279 |
| Bug commit | `0147164cfd06` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `container/state.go:131` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000b2038 by goroutine 67:
  ase/moby-22279.(*State).SetRunning()
      /work/state.go:27 +0x157
  ase/moby-22279.TestRace_22279.func2()
      /work/verified_test.go:36 +0xd5

Previous read at 0x00c0000b2038 by goroutine 24:
  ase/moby-22279.(*State).WaitRunning()
      /work/state.go:33 +0x39
  ase/moby-22279.TestRace_22279.func1()
      /work/verified_test.go:25 +0xab

Goroutine 67 (running) created at:
  ase/moby-22279.TestRace_22279()
      /work/verified_test.go:33 +0x1f1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 24 (running) created at:
  ase/moby-22279.TestRace_22279()
      /work/verified_test.go:22 +0xfa
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000b2028 by goroutine 60:
  ase/moby-22279.(*State).SetRunning()
      /work/state.go:24 +0xfe
  ase/moby-22279.TestRace_22279.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-22279-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-22279-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-22279-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-22279-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-22279-bug .
# (then run as above, no --ssh flag)
```
