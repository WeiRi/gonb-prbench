# minio-994

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/994 |
| Bug commit | `c327c56a1696` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `contentdb/contentdb.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00000078120d by goroutine 11:
  ase/minio-994.Init()
      /work/contentdb.go:13 +0xc5
  ase/minio-994.TestRaceConcurrentInit.func1()
      /work/verified_test.go:34 +0x1a0

Previous write at 0x00000078120d by goroutine 8:
  ase/minio-994.Init()
      /work/contentdb.go:17 +0x195
  ase/minio-994.TestRaceConcurrentInit.func1()
      /work/verified_test.go:34 +0x1a0

Goroutine 11 (running) created at:
  ase/minio-994.TestRaceConcurrentInit()
      /work/verified_test.go:25 +0x284
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/minio-994.TestRaceConcurrentInit()
      /work/verified_test.go:25 +0x284
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x000000720fe0 by goroutine 11:
  ase/minio-994.Init()
      /work/contentdb.go:15 +0xeb
  ase/minio-994.TestRaceConcurrentInit.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-minio-994-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-minio-994-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-minio-994-fix .
docker run --rm --memory=2g --cpus=1 gonb-minio-994-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-minio-994-bug .
# (then run as above, no --ssh flag)
```
