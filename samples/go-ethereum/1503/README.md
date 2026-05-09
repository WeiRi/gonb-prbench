# go-ethereum-1503

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/1503 |
| Bug commit | `02c5022742e2` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `accounts/account_manager.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000200d40 by goroutine 9:
  ase/go-ethereum-1503.(*Manager).expire()
      /work/account_manager.go:40 +0x174
  ase/go-ethereum-1503.TestRace_1503_account_manager_Sign_vs_expire.func2()
      /work/verified_test.go:27 +0xc7

Previous read at 0x00c000200d40 by goroutine 8:
  runtime.slicecopy()
      /usr/local/go/src/runtime/slice.go:325 +0x0
  ase/go-ethereum-1503.(*Manager).Sign()
      /work/account_manager.go:57 +0x108
  ase/go-ethereum-1503.TestRace_1503_account_manager_Sign_vs_expire.func1()
      /work/verified_test.go:21 +0x193

Goroutine 9 (running) created at:
  ase/go-ethereum-1503.TestRace_1503_account_manager_Sign_vs_expire()
      /work/verified_test.go:24 +0x284
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/go-ethereum-1503.TestRace_1503_account_manager_Sign_vs_expire()
      /work/verified_test.go:17 +0x1b9
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRace_1503_account_manager_Sign_vs_expire (0.00s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	ase/go-ethereum-1503	0.019s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-1503-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-1503-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-1503-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-1503-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-1503-bug .
# (then run as above, no --ssh flag)
```
