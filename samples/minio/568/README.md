# minio-568

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/568 |
| Bug commit | `3fc9b4554f99` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `memory/intelligent.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00011e570 by goroutine 11:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/minio-568.(*Intelligent).ExpireObjects()
      /work/intelligent.go:37 +0xa4
  ase/minio-568.TestRaceExpireObjects.func2()
      /work/verified_test.go:29 +0xa9

Previous write at 0x00c00011e570 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/minio-568.(*Intelligent).Set()
      /work/intelligent.go:30 +0x109
  ase/minio-568.TestRaceExpireObjects.func1()
      /work/verified_test.go:23 +0x1e8
  ase/minio-568.TestRaceExpireObjects.gowrap1()
      /work/verified_test.go:25 +0x41

Goroutine 11 (running) created at:
  ase/minio-568.TestRaceExpireObjects()
      /work/verified_test.go:26 +0x274
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/minio-568.TestRaceExpireObjects()
      /work/verified_test.go:20 +0x3e4
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-minio-568-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-minio-568-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-minio-568-fix .
docker run --rm --memory=2g --cpus=1 gonb-minio-568-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-minio-568-bug .
# (then run as above, no --ssh flag)
```
