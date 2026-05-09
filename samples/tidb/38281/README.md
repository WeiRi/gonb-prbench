# tidb-38281

| Field | Value |
|---|---|
| Project | tidb |
| Reference | https://github.com/pingcap/tidb/pull/38281 |
| Bug commit | `556daf722ecb` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `expression/collation.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000bc8c8 by goroutine 9:
  ase/tidb-38281.(*collationInfo).HasCoercibility()
      /work/collation.go:17 +0xa7
  ase/tidb-38281.TestRace_tidb_38281.func2()
      /work/verified_test.go:24 +0x9f

Previous write at 0x00c0000bc8c8 by goroutine 8:
  ase/tidb-38281.(*collationInfo).SetCoercibility()
      /work/collation.go:26 +0xc8
  ase/tidb-38281.TestRace_tidb_38281.func1()
      /work/verified_test.go:18 +0xa9

Goroutine 9 (running) created at:
  ase/tidb-38281.TestRace_tidb_38281()
      /work/verified_test.go:21 +0x56
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/tidb-38281.TestRace_tidb_38281()
      /work/verified_test.go:15 +0x17c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_tidb_38281 (0.00s)
FAIL
FAIL	ase/tidb-38281	0.020s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-tidb-38281-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-tidb-38281-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-tidb-38281-fix .
docker run --rm --memory=2g --cpus=1 gonb-tidb-38281-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-tidb-38281-bug .
# (then run as above, no --ssh flag)
```
