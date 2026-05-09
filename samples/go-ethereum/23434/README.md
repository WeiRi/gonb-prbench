# go-ethereum-23434

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/23434 |
| Bug commit | `c368f728c19e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `p2p/dial.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00012816c by goroutine 9:
  ase/go-ethereum-23434.(*dialScheduler).Register()
      /work/dial.go:30 +0x44
  ase/go-ethereum-23434.TestRace_23434_conn_flags.func2()
      /work/verified_test.go:24 +0xb3

Previous write at 0x00c00012816c by goroutine 8:
  ase/go-ethereum-23434.(*conn).SetFlag()
      /work/dial.go:36 +0xa8
  ase/go-ethereum-23434.TestRace_23434_conn_flags.func1()
      /work/verified_test.go:18 +0xa0

Goroutine 9 (running) created at:
  ase/go-ethereum-23434.TestRace_23434_conn_flags()
      /work/verified_test.go:21 +0x29c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/go-ethereum-23434.TestRace_23434_conn_flags()
      /work/verified_test.go:15 +0x1b0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRace_23434_conn_flags (0.00s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	ase/go-ethereum-23434	0.019s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-23434-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-23434-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-23434-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-23434-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-23434-bug .
# (then run as above, no --ssh flag)
```
