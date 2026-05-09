# minio-544

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/544 |
| Bug commit | `710e732cf072` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/storage/drivers/memory/memory.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4570 by goroutine 9:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/minio-544.(*memoryDriver).CreateObject()
      /work/memory.go:45 +0xca
  ase/minio-544.TestRaceExpireObjects.func1()
      /work/verified_test.go:37 +0x167
  ase/minio-544.TestRaceExpireObjects.gowrap1()
      /work/verified_test.go:39 +0x41

Previous read at 0x00c0000d4570 by goroutine 8:
  ase/minio-544.(*memoryDriver).expireObjects()
      /work/memory.go:54 +0x5b
  ase/minio-544.Start.gowrap1()
      /work/memory.go:31 +0x17

Goroutine 9 (running) created at:
  ase/minio-544.TestRaceExpireObjects()
      /work/verified_test.go:30 +0x14b
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/minio-544.Start()
      /work/memory.go:31 +0x1a4
  ase/minio-544.TestRaceExpireObjects()
      /work/verified_test.go:19 +0x3d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-minio-544-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-minio-544-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-minio-544-fix .
docker run --rm --memory=2g --cpus=1 gonb-minio-544-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-minio-544-bug .
# (then run as above, no --ssh flag)
```
