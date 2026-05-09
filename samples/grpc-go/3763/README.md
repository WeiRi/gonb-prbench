# grpc-go-3763

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/3763 |
| Bug commit | `dfc0c05b2da9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `interop/xds/client/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00012814f by goroutine 8:
  ase/grpc-go-3763.TestRace_3763.func1()
      /work/verified_test.go:14 +0xa4

Previous read at 0x00c00012814f by goroutine 102:
  ase/grpc-go-3763.TestRace_3763.func2()
      /work/verified_test.go:17 +0xa4

Goroutine 8 (running) created at:
  ase/grpc-go-3763.TestRace_3763()
      /work/verified_test.go:13 +0x7e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 102 (finished) created at:
  ase/grpc-go-3763.TestRace_3763()
      /work/verified_test.go:16 +0x151
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00012814f by goroutine 26:
  ase/grpc-go-3763.TestRace_3763.func1()
      /work/verified_test.go:14 +0xa4

Previous write at 0x00c00012814f by goroutine 48:
  ase/grpc-go-3763.TestRace_3763.func1()
      /work/verified_test.go:14 +0xa4
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-3763-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3763-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-3763-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-3763-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-3763-bug .
# (then run as above, no --ssh flag)
```
