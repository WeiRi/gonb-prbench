# grpc-go-2074

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/2074 |
| Bug commit | `f66923519395` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `transport/http2_server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000166010 by goroutine 8:
  ase/grpc-go-2074.(*Stream).WriteHeader()
      /work/http2_server.go:23 +0x464
  ase/grpc-go-2074.TestRace_PR2074_WriteStatusVsWriteHeader.func1()
      /work/verified_test.go:35 +0xa8

Previous write at 0x00c000166010 by goroutine 9:
  ase/grpc-go-2074.(*Stream).WriteStatus()
      /work/http2_server.go:30 +0xa4
  ase/grpc-go-2074.TestRace_PR2074_WriteStatusVsWriteHeader.func2()
      /work/verified_test.go:39 +0x86

Goroutine 8 (running) created at:
  ase/grpc-go-2074.TestRace_PR2074_WriteStatusVsWriteHeader()
      /work/verified_test.go:33 +0x2ec
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/grpc-go-2074.TestRace_PR2074_WriteStatusVsWriteHeader()
      /work/verified_test.go:37 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0000b2050 by goroutine 17:
  ase/grpc-go-2074.(*Stream).WriteStatus()
      /work/http2_server.go:29 +0x8e
  ase/grpc-go-2074.TestRace_PR2074_WriteStatusVsWriteHeader.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-2074-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-2074-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-2074-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-2074-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-2074-bug .
# (then run as above, no --ssh flag)
```
