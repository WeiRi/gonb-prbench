# cockroach-123215

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/123215 |
| Bug commit | `7cd5b546000b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kv/kvserver/store.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000190540 by goroutine 17:
  ase/cockroach-123215.TestRace_123215.func1()
      /work/verified_test.go:20 +0x98

Previous write at 0x00c000190540 by goroutine 8:
  ase/cockroach-123215.(*monitorTracer).tracker()
      /work/monitor.go:35 +0x5c
  ase/cockroach-123215.NewMonitor.gowrap1()
      /work/monitor.go:23 +0x33

Goroutine 17 (running) created at:
  ase/cockroach-123215.TestRace_123215()
      /work/verified_test.go:18 +0x93
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/cockroach-123215.NewMonitor()
      /work/monitor.go:23 +0x130
  ase/cockroach-123215.TestRace_123215()
      /work/verified_test.go:13 +0x44
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000510008 by goroutine 69:
  ase/cockroach-123215.(*monitorTracer).tracker()
      /work/monitor.go:36 +0x75
  ase/cockroach-123215.NewMonitor.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-123215-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-123215-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-123215-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-123215-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-123215-bug .
# (then run as above, no --ssh flag)
```
