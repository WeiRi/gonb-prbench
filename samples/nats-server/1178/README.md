# nats-server-1178

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/1178 |
| Bug commit | `636ff9562755` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/events.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4540 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/nats-server-1178.(*Server).debugSubscribers.func2()
      /work/events.go:27 +0x6e
  runtime.deferreturn()
      /usr/local/go/src/runtime/panic.go:602 +0x5d
  ase/nats-server-1178.TestRaceDebugSubscribers.func1()
      /work/verified_test.go:31 +0x104
  ase/nats-server-1178.TestRaceDebugSubscribers.gowrap1()
      /work/verified_test.go:32 +0x41

Previous write at 0x00c0000d4540 by goroutine 9:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/nats-server-1178.(*Server).debugSubscribers.func2()
      /work/events.go:27 +0x6e
  runtime.deferreturn()
      /usr/local/go/src/runtime/panic.go:602 +0x5d
  ase/nats-server-1178.TestRaceDebugSubscribers.func1()
      /work/verified_test.go:31 +0x104
  ase/nats-server-1178.TestRaceDebugSubscribers.gowrap1()
      /work/verified_test.go:32 +0x41

Goroutine 8 (running) created at:
  ase/nats-server-1178.TestRaceDebugSubscribers()
      /work/verified_test.go:29 +0xcd
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/nats-server-1178.TestRaceDebugSubscribers()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-1178-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-1178-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-1178-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-1178-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-1178-bug .
# (then run as above, no --ssh flag)
```
