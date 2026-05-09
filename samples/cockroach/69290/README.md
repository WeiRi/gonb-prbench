# cockroach-69290

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/69290 |
| Bug commit | `6c05f997bb82` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/sql/sqlliveness/slinstance/slinstance.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000112908 by goroutine 9:
  ase/cockroach-69290.(*session).Expiration()
      /work/slinstance.go:21 +0xb8
  ase/cockroach-69290.Test69290Race.func1()
      /work/verified_test.go:20 +0xb3
  ase/cockroach-69290.Test69290Race.gowrap1()
      /work/verified_test.go:22 +0x41

Previous write at 0x00c000112908 by goroutine 12:
  ase/cockroach-69290.(*Instance).ExtendSession()
      /work/slinstance.go:35 +0xa9
  ase/cockroach-69290.Test69290Race.func1()
      /work/verified_test.go:18 +0xa5
  ase/cockroach-69290.Test69290Race.gowrap1()
      /work/verified_test.go:22 +0x41

Goroutine 9 (running) created at:
  ase/cockroach-69290.Test69290Race()
      /work/verified_test.go:15 +0x177
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 12 (finished) created at:
  ase/cockroach-69290.Test69290Race()
      /work/verified_test.go:15 +0x177
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: Test69290Race (0.00s)
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-69290-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-69290-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-69290-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-69290-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-69290-bug .
# (then run as above, no --ssh flag)
```
