# cockroach-61854

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/61854 |
| Bug commit | `a86cf540219a` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/util/timeutil/stopwatch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d9020 by goroutine 8:
  ase/cockroach-61854.(*Stopwatch).Start()
      /work/stopwatch.go:17 +0xae
  ase/cockroach-61854.TestRace_61854.func1()
      /work/verified_test.go:38 +0x86

Previous write at 0x00c0000d9020 by goroutine 12:
  ase/cockroach-61854.(*Stopwatch).Start()
      /work/stopwatch.go:17 +0xae
  ase/cockroach-61854.TestRace_61854.func1()
      /work/verified_test.go:38 +0x86

Goroutine 8 (running) created at:
  ase/cockroach-61854.TestRace_61854()
      /work/verified_test.go:36 +0xd8
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 12 (finished) created at:
  ase/cockroach-61854.TestRace_61854()
      /work/verified_test.go:36 +0xd8
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000d9040 by goroutine 8:
  ase/cockroach-61854.(*Stopwatch).Start()
      /work/stopwatch.go:18 +0x108
  ase/cockroach-61854.TestRace_61854.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-61854-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-61854-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-61854-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-61854-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-61854-bug .
# (then run as above, no --ssh flag)
```
