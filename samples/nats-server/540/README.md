# nats-server-540

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/540 |
| Bug commit | `1ca5e57b2d65` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/monitor.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00011c188 by goroutine 126:
  ase/nats-server-540.TestRaceHandleRoutezNc.func2()
      /work/verified_test.go:65 +0xd7

Previous read at 0x00c00011c188 by goroutine 91:
  ase/nats-server-540.TestRaceHandleRoutezNc.func1()
      /work/verified_test.go:51 +0xee

Goroutine 126 (running) created at:
  ase/nats-server-540.TestRaceHandleRoutezNc()
      /work/verified_test.go:60 +0x1d7
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 91 (running) created at:
  ase/nats-server-540.TestRaceHandleRoutezNc()
      /work/verified_test.go:39 +0xa4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c00011c188 by goroutine 92:
  ase/nats-server-540.TestRaceHandleRoutezNc.func1()
      /work/verified_test.go:51 +0xee

Previous write at 0x00c00011c188 by goroutine 127:
  ase/nats-server-540.TestRaceHandleRoutezNc.func2()
      /work/verified_test.go:68 +0xfd
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-540-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-540-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-540-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-540-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-540-bug .
# (then run as above, no --ssh flag)
```
