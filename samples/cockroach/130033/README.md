# cockroach-130033

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/130033 |
| Bug commit | `b120a8a9fe6c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/server/diagnostics/reporter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00011c120 by goroutine 13:
  ase/cockroach-130033.(*Reporter).Read()
      /work/reporter.go:21 +0x113
  ase/cockroach-130033.Test130031Race.func1()
      /work/verified_test.go:18 +0x109
  ase/cockroach-130033.Test130031Race.gowrap1()
      /work/verified_test.go:20 +0x41

Previous write at 0x00c00011c120 by goroutine 12:
  ase/cockroach-130033.(*Reporter).ReportDiagnostics()
      /work/reporter.go:16 +0xbe
  ase/cockroach-130033.Test130031Race.func1()
      /work/verified_test.go:16 +0x96
  ase/cockroach-130033.Test130031Race.gowrap1()
      /work/verified_test.go:20 +0x41

Goroutine 13 (running) created at:
  ase/cockroach-130033.Test130031Race()
      /work/verified_test.go:13 +0xb1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 12 (finished) created at:
  ase/cockroach-130033.Test130031Race()
      /work/verified_test.go:13 +0xb1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: Test130031Race (0.00s)
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-130033-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-130033-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-130033-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-130033-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-130033-bug .
# (then run as above, no --ssh flag)
```
