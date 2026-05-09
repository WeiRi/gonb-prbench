# tidb-62900

| Field | Value |
|---|---|
| Project | tidb |
| Reference | https://github.com/pingcap/tidb/pull/62900 |
| Bug commit | `857a162ea4a5` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/ddl/executor.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00001423f by goroutine 9:
  ase/tidb-62900.ReportMinStartTS_BUG()
      /work/verified_test.go:66 +0xb7
  ase/tidb-62900.TestRace_tidb62900.func2()
      /work/verified_test.go:84 +0x90

Previous write at 0x00c00001423f by goroutine 8:
  ase/tidb-62900.(*Executor).DoDDLJobWrapper_BUG()
      /work/verified_test.go:61 +0xcc
  ase/tidb-62900.TestRace_tidb62900.func1()
      /work/verified_test.go:78 +0x9b

Goroutine 9 (running) created at:
  ase/tidb-62900.TestRace_tidb62900()
      /work/verified_test.go:81 +0x264
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/tidb-62900.TestRace_tidb62900()
      /work/verified_test.go:75 +0x1bc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRace_tidb62900 (0.00s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	ase/tidb-62900	0.017s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-tidb-62900-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-tidb-62900-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-tidb-62900-fix .
docker run --rm --memory=2g --cpus=1 gonb-tidb-62900-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-tidb-62900-bug .
# (then run as above, no --ssh flag)
```
