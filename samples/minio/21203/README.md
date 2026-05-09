# minio-21203

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/21203 |
| Bug commit | `7ee75368e055` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/grid/connection.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0002d4000 by goroutine 56:
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.func1.1()
      /work/internal/grid/race_21203_capture_test.go:35 +0x44
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.func1.gowrap2()
      /work/internal/grid/race_21203_capture_test.go:37 +0x61

Previous write at 0x00c0002d4000 by goroutine 12:
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.func1()
      /work/internal/grid/race_21203_capture_test.go:44 +0x3ed
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.gowrap1()
      /work/internal/grid/race_21203_capture_test.go:54 +0x41

Goroutine 56 (running) created at:
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.func1()
      /work/internal/grid/race_21203_capture_test.go:33 +0x364
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.gowrap1()
      /work/internal/grid/race_21203_capture_test.go:54 +0x41

Goroutine 12 (running) created at:
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle()
      /work/internal/grid/race_21203_capture_test.go:24 +0x6d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0004a4001 by goroutine 33:
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.func1()
      /work/internal/grid/race_21203_capture_test.go:29 +0x1cb
  github.com/minio/minio/internal/grid.TestRace_BufferLifecycle.gowrap1()
      /work/internal/grid/race_21203_capture_test.go:54 +0x41
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-minio-21203-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-minio-21203-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-minio-21203-fix .
docker run --rm --memory=2g --cpus=1 gonb-minio-21203-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-minio-21203-bug .
# (then run as above, no --ssh flag)
```
