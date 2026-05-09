# grpc-go-5927

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/5927 |
| Bug commit | `bf3ad352405d` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `xds/internal/xdsclient/clientimpl_loadreport.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00010e2f0 by goroutine 10:
  ase/grpc-go-5927.(*clientImpl).refLocked()
      /work/clientimpl_loadreport.go:23 +0xbb
  ase/grpc-go-5927.(*clientImpl).ReportLoad()
      /work/clientimpl_loadreport.go:34 +0xb6
  ase/grpc-go-5927.TestRace_5927.func1()
      /work/verified_test.go:42 +0x109

Previous write at 0x00c00010e2f0 by goroutine 8:
  ase/grpc-go-5927.(*clientImpl).refLocked()
      /work/clientimpl_loadreport.go:23 +0xcd
  ase/grpc-go-5927.(*clientImpl).ReportLoad()
      /work/clientimpl_loadreport.go:34 +0xb6
  ase/grpc-go-5927.TestRace_5927.func1()
      /work/verified_test.go:42 +0x109

Goroutine 10 (running) created at:
  ase/grpc-go-5927.TestRace_5927()
      /work/verified_test.go:38 +0x28a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/grpc-go-5927.TestRace_5927()
      /work/verified_test.go:38 +0x28a
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-5927-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-5927-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-5927-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-5927-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-5927-bug .
# (then run as above, no --ssh flag)
```
