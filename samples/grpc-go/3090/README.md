# grpc-go-3090

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/3090 |
| Bug commit | `f07f2cffa0a1` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `resolver_conn_wrapper.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000296168 by goroutine 250:
  ase/grpc-go-3090.NewCCResolverWrapper.func1()
      /work/resolver_conn_wrapper.go:16 +0x45

Previous write at 0x00c000296168 by goroutine 116:
  ase/grpc-go-3090.NewCCResolverWrapper()
      /work/resolver_conn_wrapper.go:18 +0x14f
  ase/grpc-go-3090.TestRace_PR3090_ResolverWrapperBuild.func1()
      /work/verified_test.go:37 +0xc4

Goroutine 250 (running) created at:
  ase/grpc-go-3090.NewCCResolverWrapper()
      /work/resolver_conn_wrapper.go:14 +0x11c
  ase/grpc-go-3090.TestRace_PR3090_ResolverWrapperBuild.func1()
      /work/verified_test.go:37 +0xc4

Goroutine 116 (finished) created at:
  ase/grpc-go-3090.TestRace_PR3090_ResolverWrapperBuild()
      /work/verified_test.go:35 +0x89
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR3090_ResolverWrapperBuild (0.00s)
FAIL
FAIL	ase/grpc-go-3090	0.021s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-3090-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3090-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-3090-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3090-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-3090-bug .
# (then run as above, no --ssh flag)
```
