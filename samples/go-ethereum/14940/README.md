# go-ethereum-14940

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/14940 |
| Bug commit | `e063d538b80d` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `core/tx_pool.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d8540 by goroutine 8:
  runtime.mapassign_fast64()
      /usr/local/go/src/runtime/map_fast64.go:93 +0x0
  ase/go-ethereum-14940.(*TxPool).Reset()
      /work/tx_pool.go:38 +0x104
  ase/go-ethereum-14940.TestRace_14940_journal_rotate_vs_reset.func1()
      /work/verified_test.go:18 +0xa9

Previous read at 0x00c0000d8540 by goroutine 9:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/go-ethereum-14940.(*TxPool).RotateJournal()
      /work/tx_pool.go:25 +0x18d
  ase/go-ethereum-14940.TestRace_14940_journal_rotate_vs_reset.func2()
      /work/verified_test.go:24 +0xea

Goroutine 8 (running) created at:
  ase/go-ethereum-14940.TestRace_14940_journal_rotate_vs_reset()
      /work/verified_test.go:15 +0x1e4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (running) created at:
  ase/go-ethereum-14940.TestRace_14940_journal_rotate_vs_reset()
      /work/verified_test.go:21 +0x28c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRace_14940_journal_rotate_vs_reset (0.00s)
    testing.go:1398: race detected during execution of test
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-14940-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-14940-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-14940-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-14940-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-14940-bug .
# (then run as above, no --ssh flag)
```
