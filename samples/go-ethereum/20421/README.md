# go-ethereum-20421

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/20421 |
| Bug commit | `c2d65d34d5c6` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `p2p/enode/iter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d6518 by goroutine 9:
  ase/go-ethereum-20421.(*sliceIter).Close()
      /work/iter.go:52 +0x8d
  ase/go-ethereum-20421.TestRace_20421_sliceIter_Node_vs_Close.func2()
      /work/verified_test.go:30 +0x99

Previous read at 0x00c0000d6518 by goroutine 8:
  ase/go-ethereum-20421.(*sliceIter).Node()
      /work/iter.go:42 +0xc8
  ase/go-ethereum-20421.TestRace_20421_sliceIter_Node_vs_Close.func1()
      /work/verified_test.go:23 +0x12b

Goroutine 9 (running) created at:
  ase/go-ethereum-20421.TestRace_20421_sliceIter_Node_vs_Close()
      /work/verified_test.go:26 +0x2fc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/go-ethereum-20421.TestRace_20421_sliceIter_Node_vs_Close()
      /work/verified_test.go:20 +0x250
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRace_20421_sliceIter_Node_vs_Close (0.00s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	ase/go-ethereum-20421	0.016s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-20421-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-20421-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-20421-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-20421-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-20421-bug .
# (then run as above, no --ssh flag)
```
