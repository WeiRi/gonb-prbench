# grpc-go-1897

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/1897 |
| Bug commit | `7c5299d71e2b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `clientconn.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x0000006fc428 by goroutine 9:
  ase/grpc-go-1897.TestRace_PR1897_MinConnectTimeout.func2()
      /work/verified_test.go:80 +0x91

Previous read at 0x0000006fc428 by goroutine 8:
  ase/grpc-go-1897.resetTransport()
      /work/clientconn.go:20 +0x96
  ase/grpc-go-1897.TestRace_PR1897_MinConnectTimeout.func1()
      /work/verified_test.go:73 +0xa9

Goroutine 9 (running) created at:
  ase/grpc-go-1897.TestRace_PR1897_MinConnectTimeout()
      /work/verified_test.go:77 +0x18d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/grpc-go-1897.TestRace_PR1897_MinConnectTimeout()
      /work/verified_test.go:71 +0x110
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x0000006fc428 by goroutine 9:
  ase/grpc-go-1897.TestRace_PR1897_MinConnectTimeout.func2()
      /work/verified_test.go:82 +0xb2

Previous read at 0x0000006fc428 by goroutine 8:
  ase/grpc-go-1897.resetTransport()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-1897-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-1897-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-1897-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-1897-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-1897-bug .
# (then run as above, no --ssh flag)
```
