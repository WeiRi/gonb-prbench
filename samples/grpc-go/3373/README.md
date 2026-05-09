# grpc-go-3373

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/3373 |
| Bug commit | `f3111a575aec` |
| Category | order_violation |
| Oracle | RACE |
| Primary diff file | `internal/grpctest/tlogger.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00011e540 by goroutine 14:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/grpc-go-3373.(*TLogger).AddExpect()
      /work/tlogger.go:14 +0xfd
  ase/grpc-go-3373.TestRace_PR3373_TLoggerMap.func2()
      /work/verified_test.go:46 +0xb8

Previous read at 0x00c00011e540 by goroutine 8:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/grpc-go-3373.(*TLogger).expected()
      /work/tlogger.go:18 +0x85
  ase/grpc-go-3373.TestRace_PR3373_TLoggerMap.func1()
      /work/verified_test.go:37 +0xa5

Goroutine 14 (running) created at:
  ase/grpc-go-3373.TestRace_PR3373_TLoggerMap()
      /work/verified_test.go:42 +0x2dc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/grpc-go-3373.TestRace_PR3373_TLoggerMap()
      /work/verified_test.go:34 +0x157
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-3373-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3373-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-3373-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3373-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-3373-bug .
# (then run as above, no --ssh flag)
```
