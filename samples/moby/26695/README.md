# moby-26695

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/26695 |
| Bug commit | `45a8f6802635` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `libcontainerd/pausemonitor_linux.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00020a000 by goroutine 13:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-26695.(*pauseMonitor).append()
      /work/pausemonitor.go:14 +0x144
  ase/moby-26695.TestRace26695PauseMonitor.func2()
      /work/verified_test.go:28 +0xb7

Previous write at 0x00c00020a000 by goroutine 17:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-26695.(*pauseMonitor).append()
      /work/pausemonitor.go:14 +0x144
  ase/moby-26695.TestRace26695PauseMonitor.func2()
      /work/verified_test.go:28 +0xb7

Goroutine 13 (running) created at:
  ase/moby-26695.TestRace26695PauseMonitor()
      /work/verified_test.go:25 +0x9d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 17 (running) created at:
  ase/moby-26695.TestRace26695PauseMonitor()
      /work/verified_test.go:25 +0x9d
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-26695-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-26695-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-26695-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-26695-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-26695-bug .
# (then run as above, no --ssh flag)
```
