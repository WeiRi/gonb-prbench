# grpc-go-2695

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/2695 |
| Bug commit | `3c84def89307` |
| Category | channel_misuse |
| Oracle | RACE|PANIC |
| Primary diff file | `internal/transport/handler_server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000be190 by goroutine 15:
  runtime.closechan()
      /usr/local/go/src/runtime/chan.go:357 +0x0
  ase/grpc-go-2695.(*serverHandlerTransport).WriteStatus()
      /work/handler_server.go:31 +0x47
  ase/grpc-go-2695.TestRace_grpcgo2695.func3()
      /work/verified_test.go:49 +0x35

Previous read at 0x00c0000be190 by goroutine 12:
  runtime.chansend()
      /usr/local/go/src/runtime/chan.go:160 +0x0
  ase/grpc-go-2695.(*serverHandlerTransport).do()
      /work/handler_server.go:22 +0x286
  ase/grpc-go-2695.TestRace_grpcgo2695.func2()
      /work/verified_test.go:43 +0x1f5

Goroutine 15 (running) created at:
  ase/grpc-go-2695.TestRace_grpcgo2695()
      /work/verified_test.go:48 +0x30d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 12 (running) created at:
  ase/grpc-go-2695.TestRace_grpcgo2695()
      /work/verified_test.go:28 +0x444
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    verified_test.go:55: iter 0: caught panic stack:
        goroutine 8 [running]:
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-2695-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-2695-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-2695-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-2695-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-2695-bug .
# (then run as above, no --ssh flag)
```
