# grpc-go-4641

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/4641 |
| Bug commit | `edb9b3bc2266` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/transport/http2_client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000b2030 by goroutine 9:
  ase/grpc-go-4641.operateHeaders()
      /work/http2_client.go:17 +0x1d7
  ase/grpc-go-4641.TestRace_PR4641_RecvCompressBeforeChanClose.func2()
      /work/verified_test.go:38 +0x99

Previous write at 0x00c0000b2030 by goroutine 10:
  ase/grpc-go-4641.raceWriter()
      /work/http2_client.go:30 +0x89
  ase/grpc-go-4641.TestRace_PR4641_RecvCompressBeforeChanClose.func3()
      /work/verified_test.go:42 +0x81

Goroutine 9 (running) created at:
  ase/grpc-go-4641.TestRace_PR4641_RecvCompressBeforeChanClose()
      /work/verified_test.go:36 +0x2cc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 10 (finished) created at:
  ase/grpc-go-4641.TestRace_PR4641_RecvCompressBeforeChanClose()
      /work/verified_test.go:40 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR4641_RecvCompressBeforeChanClose (0.00s)
FAIL
FAIL	ase/grpc-go-4641	0.018s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-4641-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-4641-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-4641-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-4641-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-4641-bug .
# (then run as above, no --ssh flag)
```
