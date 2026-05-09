# go-ethereum-19162

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/19162 |
| Bug commit | `81babe15090e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `swarm/pss/handshake.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d8510 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/go-ethereum-19162.(*HandshakeController).Insert()
      /work/handshake.go:31 +0x116
  ase/go-ethereum-19162.TestRace_19162_HandshakeController_symKeyIndex.func1()
      /work/verified_test.go:20 +0xa5

Previous read at 0x00c0000d8510 by goroutine 9:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/go-ethereum-19162.(*HandshakeController).Handler()
      /work/handshake.go:37 +0xf6
  ase/go-ethereum-19162.TestRace_19162_HandshakeController_symKeyIndex.func2()
      /work/verified_test.go:26 +0xcd

Goroutine 8 (running) created at:
  ase/go-ethereum-19162.TestRace_19162_HandshakeController_symKeyIndex()
      /work/verified_test.go:17 +0x1c4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/go-ethereum-19162.TestRace_19162_HandshakeController_symKeyIndex()
      /work/verified_test.go:23 +0x26c
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-19162-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-19162-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-19162-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-19162-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-19162-bug .
# (then run as above, no --ssh flag)
```
