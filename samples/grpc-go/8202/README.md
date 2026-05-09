# grpc-go-8202

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/8202 |
| Bug commit | `208e03b3bae2` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/resolver/delegatingresolver/delegatingresolver.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000012140 by goroutine 9:
  ase/grpc-go-8202.(*Resolver).Close()
      /work/delegatingresolver.go:43 +0x8d
  ase/grpc-go-8202.TestRace_PR8202_DelegatingResolverChildMu.func2()
      /work/verified_test.go:37 +0x85

Previous read at 0x00c000012140 by goroutine 8:
  ase/grpc-go-8202.(*Resolver).updateProxyResolverState()
      /work/delegatingresolver.go:26 +0x92
  ase/grpc-go-8202.TestRace_PR8202_DelegatingResolverChildMu.func1()
      /work/verified_test.go:32 +0x8a

Goroutine 9 (running) created at:
  ase/grpc-go-8202.TestRace_PR8202_DelegatingResolverChildMu()
      /work/verified_test.go:35 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/grpc-go-8202.TestRace_PR8202_DelegatingResolverChildMu()
      /work/verified_test.go:30 +0x24c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000012148 by goroutine 9:
  ase/grpc-go-8202.(*Resolver).Close()
      /work/delegatingresolver.go:44 +0xc4
  ase/grpc-go-8202.TestRace_PR8202_DelegatingResolverChildMu.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-8202-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-8202-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-8202-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-8202-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-8202-bug .
# (then run as above, no --ssh flag)
```
