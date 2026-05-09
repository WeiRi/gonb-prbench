# go-ethereum-17173

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/17173 |
| Bug commit | `f7d3678c28c4` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `miner/worker.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000124510 by goroutine 9:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/go-ethereum-17173.(*StateDB).Copy()
      /work/worker.go:25 +0x1a4
  ase/go-ethereum-17173.(*worker).Pending()
      /work/worker.go:58 +0xb4
  ase/go-ethereum-17173.TestRace_17173_state_concurrent_map.func2()
      /work/verified_test.go:24 +0x99

Previous write at 0x00c000124510 by goroutine 8:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/go-ethereum-17173.(*StateDB).Update()
      /work/worker.go:21 +0x14d
  ase/go-ethereum-17173.(*worker).Wait()
      /work/worker.go:48 +0xdd
  ase/go-ethereum-17173.TestRace_17173_state_concurrent_map.func1()
      /work/verified_test.go:18 +0xbc

Goroutine 9 (running) created at:
  ase/go-ethereum-17173.TestRace_17173_state_concurrent_map()
      /work/verified_test.go:21 +0x2dc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/go-ethereum-17173.TestRace_17173_state_concurrent_map()
      /work/verified_test.go:15 +0x22c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-17173-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-17173-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-17173-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-17173-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-17173-bug .
# (then run as above, no --ssh flag)
```
