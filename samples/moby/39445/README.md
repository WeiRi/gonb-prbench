# moby-39445

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/39445 |
| Bug commit | `e105a74c5419` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/ioutils/bytespipe.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0001828a8 by goroutine 58:
  ase/moby-39445.(*BytesPipe).Close()
      /work/bytespipe.go:42 +0x5b
  ase/moby-39445.TestRace_39445.func2()
      /work/verified_test.go:44 +0x84

Previous read at 0x00c0001828a8 by goroutine 40:
  ase/moby-39445.(*BytesPipe).Read()
      /work/bytespipe.go:34 +0xaf
  ase/moby-39445.TestRace_39445.func1()
      /work/verified_test.go:24 +0xe6

Goroutine 58 (running) created at:
  ase/moby-39445.TestRace_39445()
      /work/verified_test.go:42 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 40 (finished) created at:
  ase/moby-39445.TestRace_39445()
      /work/verified_test.go:20 +0x258
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_39445 (0.28s)
FAIL
FAIL	ase/moby-39445	0.302s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-39445-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-39445-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-39445-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-39445-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-39445-bug .
# (then run as above, no --ssh flag)
```
