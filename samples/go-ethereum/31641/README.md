# go-ethereum-31641

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/31641 |
| Bug commit | `476f117211d0` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `core/txpool/legacypool/legacypool.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000b6038 by goroutine 9:
  ase/go-ethereum-31641.(*LegacyPool).Clear()
      /work/legacypool.go:66 +0x157
  ase/go-ethereum-31641.TestRace_31641_legacypool_Add_vs_Clear.func2()
      /work/verified_test.go:27 +0xb1

Previous read at 0x00c0000b6038 by goroutine 8:
  ase/go-ethereum-31641.(*LegacyPool).Add()
      /work/legacypool.go:55 +0xb2
  ase/go-ethereum-31641.TestRace_31641_legacypool_Add_vs_Clear.func1()
      /work/verified_test.go:21 +0xe4

Goroutine 9 (running) created at:
  ase/go-ethereum-31641.TestRace_31641_legacypool_Add_vs_Clear()
      /work/verified_test.go:24 +0x450
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/go-ethereum-31641.TestRace_31641_legacypool_Add_vs_Clear()
      /work/verified_test.go:18 +0x3a7
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000201d30 by goroutine 8:
  ase/go-ethereum-31641.(*lookup).Get()
      /work/legacypool.go:31 +0x9d
  ase/go-ethereum-31641.(*LegacyPool).Add()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-31641-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-31641-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-31641-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-31641-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-31641-bug .
# (then run as above, no --ssh flag)
```
