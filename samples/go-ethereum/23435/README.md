# go-ethereum-23435

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/23435 |
| Bug commit | `c368f728c19e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `miner/worker.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000022310 by goroutine 7:
  ase/go-ethereum-23435.(*worker).Close()
      /work/worker.go:35 +0x1b9
  ase/go-ethereum-23435.TestRace_23435_worker_close_vs_mainLoop_current()
      /work/verified_test.go:15 +0x1af
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Previous write at 0x00c000022310 by goroutine 8:
  ase/go-ethereum-23435.(*worker).MainLoop()
      /work/worker.go:49 +0xe4
  ase/go-ethereum-23435.TestRace_23435_worker_close_vs_mainLoop_current.gowrap1()
      /work/verified_test.go:13 +0x1f

Goroutine 7 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1742 +0x825
  testing.runTests.func1()
      /usr/local/go/src/testing/testing.go:2161 +0x85
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.runTests()
      /usr/local/go/src/testing/testing.go:2159 +0x8be
  testing.(*M).Run()
      /usr/local/go/src/testing/testing.go:2027 +0xf17
  main.main()
      _testmain.go:47 +0x2bd

Goroutine 8 (running) created at:
  ase/go-ethereum-23435.TestRace_23435_worker_close_vs_mainLoop_current()
      /work/verified_test.go:13 +0x1a4
  testing.tRunner()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-23435-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-23435-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-23435-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-23435-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-23435-bug .
# (then run as above, no --ssh flag)
```
