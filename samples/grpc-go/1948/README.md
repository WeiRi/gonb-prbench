# grpc-go-1948

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/1948 |
| Bug commit | `f72b28a6d170` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `call.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000110000 by goroutine 18:
  runtime.slicecopy()
      /usr/local/go/src/runtime/slice.go:325 +0x0
  ase/grpc-go-1948.(*ClientConn).Invoke()
      /work/call.go:19 +0xef
  ase/grpc-go-1948.TestRace_1948_SliceBackingArray.func1.1()
      /work/verified_test.go:34 +0x111
  ase/grpc-go-1948.TestRace_1948_SliceBackingArray.func1.gowrap2()
      /work/verified_test.go:36 +0x41

Previous write at 0x00c000110000 by goroutine 19:
  runtime.slicecopy()
      /usr/local/go/src/runtime/slice.go:325 +0x0
  ase/grpc-go-1948.(*ClientConn).Invoke()
      /work/call.go:19 +0xef
  ase/grpc-go-1948.TestRace_1948_SliceBackingArray.func1.1()
      /work/verified_test.go:34 +0x111
  ase/grpc-go-1948.TestRace_1948_SliceBackingArray.func1.gowrap2()
      /work/verified_test.go:36 +0x41

Goroutine 18 (running) created at:
  ase/grpc-go-1948.TestRace_1948_SliceBackingArray.func1()
      /work/verified_test.go:31 +0x177

Goroutine 19 (running) created at:
  ase/grpc-go-1948.TestRace_1948_SliceBackingArray.func1()
      /work/verified_test.go:31 +0x177
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_1948_SliceBackingArray (0.03s)
FAIL
FAIL	ase/grpc-go-1948	0.050s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-1948-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-1948-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-1948-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-1948-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-1948-bug .
# (then run as above, no --ssh flag)
```
