# grpc-go-7128

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/7128 |
| Bug commit | `e22436abb809` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `credentials/tls/certprovider/store.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00012017f by goroutine 11:
  ase/grpc-go-7128.(*provider_7128).Close()
      /work/verified_test.go:10 +0x91
  ase/grpc-go-7128.TestRace_7128.func1()
      /work/verified_test.go:18 +0x72

Previous write at 0x00c00012017f by goroutine 12:
  ase/grpc-go-7128.(*provider_7128).Close()
      /work/verified_test.go:10 +0x91
  ase/grpc-go-7128.TestRace_7128.func1()
      /work/verified_test.go:18 +0x72

Goroutine 11 (running) created at:
  ase/grpc-go-7128.TestRace_7128()
      /work/verified_test.go:17 +0x109
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 12 (finished) created at:
  ase/grpc-go-7128.TestRace_7128()
      /work/verified_test.go:17 +0x109
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
--- FAIL: TestRace_7128 (0.04s)
    testing.go:1617: race detected during execution of test
FAIL
FAIL	ase/grpc-go-7128	0.054s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-7128-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-7128-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-7128-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-7128-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-7128-bug .
# (then run as above, no --ssh flag)
```
