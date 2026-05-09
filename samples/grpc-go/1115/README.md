# grpc-go-1115

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/1115 |
| Bug commit | `68a5d50f4517` |
| Category | channel_misuse |
| Oracle | RACE |
| Primary diff file | `transport/handler_server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000182310 by goroutine 9:
  runtime.chansend()
      /usr/local/go/src/runtime/chan.go:160 +0x0
  ase/grpc-go-1115.(*serverHandlerTransport).do()
      /work/handler_server.go:20 +0x1bc
  ase/grpc-go-1115.TestRace_PR1115.func2()
      /work/verified_test.go:143 +0x20a

Previous write at 0x00c000182310 by goroutine 8:
  runtime.closechan()
      /usr/local/go/src/runtime/chan.go:357 +0x0
  ase/grpc-go-1115.TestRace_PR1115.func1()
      /work/verified_test.go:130 +0xc9

Goroutine 9 (running) created at:
  ase/grpc-go-1115.TestRace_PR1115()
      /work/verified_test.go:137 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/grpc-go-1115.TestRace_PR1115()
      /work/verified_test.go:127 +0x24c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR1115 (0.01s)
FAIL
FAIL	ase/grpc-go-1115	0.023s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-1115-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-1115-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-1115-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-1115-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-1115-bug .
# (then run as above, no --ssh flag)
```
