# cockroach-156275

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/156275 |
| Bug commit | `e7de41dae5f9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/sql/catalog/lease/lease.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00012e16c by goroutine 9:
  ase/cockroach-156275.(*Manager).WatchHandler()
      /work/lease.go:21 +0xae
  ase/cockroach-156275.Test156275Race.func1()
      /work/verified_test.go:18 +0xa9
  ase/cockroach-156275.Test156275Race.gowrap1()
      /work/verified_test.go:20 +0x41

Previous write at 0x00c00012e16c by goroutine 10:
  ase/cockroach-156275.(*Manager).TestingSetDisable()
      /work/lease.go:31 +0x8d
  ase/cockroach-156275.Test156275Race.func1()
      /work/verified_test.go:16 +0x9b
  ase/cockroach-156275.Test156275Race.gowrap1()
      /work/verified_test.go:20 +0x41

Goroutine 9 (running) created at:
  ase/cockroach-156275.Test156275Race()
      /work/verified_test.go:13 +0x9e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 10 (finished) created at:
  ase/cockroach-156275.Test156275Race()
      /work/verified_test.go:13 +0x9e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: Test156275Race (0.01s)
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-156275-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-156275-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-156275-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-156275-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-156275-bug .
# (then run as above, no --ssh flag)
```
