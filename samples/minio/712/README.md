# minio-712

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/712 |
| Bug commit | `5132fd84db9a` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/donut/disk/disk.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4540 by goroutine 17:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/minio-712.Disk.GetFSInfo()
      /work/disk.go:11 +0x104
  ase/minio-712.TestRaceGetFSInfo.func1()
      /work/verified_test.go:30 +0xe4

Previous write at 0x00c0000d4540 by goroutine 13:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/minio-712.Disk.GetFSInfo()
      /work/disk.go:11 +0x104
  ase/minio-712.TestRaceGetFSInfo.func1()
      /work/verified_test.go:30 +0xe4

Goroutine 17 (running) created at:
  ase/minio-712.TestRaceGetFSInfo()
      /work/verified_test.go:25 +0x7c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 13 (finished) created at:
  ase/minio-712.TestRaceGetFSInfo()
      /work/verified_test.go:25 +0x7c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-minio-712-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-minio-712-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-minio-712-fix .
docker run --rm --memory=2g --cpus=1 gonb-minio-712-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-minio-712-bug .
# (then run as above, no --ssh flag)
```
