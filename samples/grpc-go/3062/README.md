# grpc-go-3062

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/3062 |
| Bug commit | `7aa94b7eefde` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/transport/transport.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000166010 by goroutine 9:
  ase/grpc-go-3062.(*Stream).operateHeaders()
      /work/transport.go:22 +0x89
  ase/grpc-go-3062.TestRace_PR3062_WaitOnHeaderRecvCompress.func2()
      /work/verified_test.go:42 +0x81

Previous read at 0x00c000166010 by goroutine 8:
  ase/grpc-go-3062.(*Stream).RecvCompress()
      /work/transport.go:15 +0x17b
  ase/grpc-go-3062.TestRace_PR3062_WaitOnHeaderRecvCompress.func1()
      /work/verified_test.go:38 +0x17c

Goroutine 9 (running) created at:
  ase/grpc-go-3062.TestRace_PR3062_WaitOnHeaderRecvCompress()
      /work/verified_test.go:40 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/grpc-go-3062.TestRace_PR3062_WaitOnHeaderRecvCompress()
      /work/verified_test.go:35 +0x22c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR3062_WaitOnHeaderRecvCompress (0.01s)
FAIL
FAIL	ase/grpc-go-3062	0.024s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-3062-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3062-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-3062-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3062-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-3062-bug .
# (then run as above, no --ssh flag)
```
