# grpc-go-8519

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/8519 |
| Bug commit | `fa0d65832080` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/transport/handler_server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000114010 by goroutine 34:
  ase/grpc-go-8519.(*Stream).SetTrailer()
      /work/handler_server.go:27 +0x492
  ase/grpc-go-8519.TestRace_HandlerServer_TrailerCopy_8519.func1.1()
      /work/verified_test.go:25 +0xc9

Previous read at 0x00c000114010 by goroutine 37:
  ase/grpc-go-8519.(*serverHandlerTransport).writeStatus()
      /work/handler_server.go:38 +0xe4
  ase/grpc-go-8519.TestRace_HandlerServer_TrailerCopy_8519.func1.2()
      /work/verified_test.go:31 +0xd8

Goroutine 34 (running) created at:
  ase/grpc-go-8519.TestRace_HandlerServer_TrailerCopy_8519.func1()
      /work/verified_test.go:22 +0x364

Goroutine 37 (finished) created at:
  ase/grpc-go-8519.TestRace_HandlerServer_TrailerCopy_8519.func1()
      /work/verified_test.go:28 +0x451
==================
==================
WARNING: DATA RACE
Read at 0x00c0002099b0 by goroutine 30:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/grpc-go-8519.(*serverHandlerTransport).writeStatus()
      /work/handler_server.go:38 +0x124
  ase/grpc-go-8519.TestRace_HandlerServer_TrailerCopy_8519.func1.2()
      /work/verified_test.go:31 +0xd8

Previous write at 0x00c0002099b0 by goroutine 27:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/grpc-go-8519.MD.Join()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-8519-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-8519-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-8519-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-8519-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-8519-bug .
# (then run as above, no --ssh flag)
```
