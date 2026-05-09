# moby-47961

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/47961 |
| Bug commit | `ff5cc18482be` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000c23f0 by goroutine 14:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-47961.(*Client).AddHeader()
      /work/client.go:22 +0xfc
  ase/moby-47961.TestRace_47961.func1()
      /work/verified_test.go:22 +0xd7
  ase/moby-47961.TestRace_47961.func3()
      /work/verified_test.go:24 +0x41

Previous write at 0x00c0000c23f0 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-47961.(*Client).AddHeader()
      /work/client.go:22 +0xfc
  ase/moby-47961.TestRace_47961.func1()
      /work/verified_test.go:22 +0xd7
  ase/moby-47961.TestRace_47961.func3()
      /work/verified_test.go:24 +0x41

Goroutine 14 (running) created at:
  ase/moby-47961.TestRace_47961()
      /work/verified_test.go:19 +0x104
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 8 (finished) created at:
  ase/moby-47961.TestRace_47961()
      /work/verified_test.go:19 +0x104
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-47961-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-47961-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-47961-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-47961-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-47961-bug .
# (then run as above, no --ssh flag)
```
