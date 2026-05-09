# pq-1000

| Field | Value |
|---|---|
| Project | pq |
| Reference | https://github.com/lib/pq/pull/1000 |
| Bug commit | `be2b75c2254d` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `conn.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0001a817f by goroutine 8:
  ase/pq-1000.(*conn).setBad()
      /work/conn.go:9 +0xa4
  ase/pq-1000.TestRace_pq_1000.func1()
      /work/verified_test.go:18 +0x9b

Previous read at 0x00c0001a817f by goroutine 11:
  ase/pq-1000.(*conn).getBad()
      /work/conn.go:11 +0xa4
  ase/pq-1000.TestRace_pq_1000.func2()
      /work/verified_test.go:24 +0x9b

Goroutine 8 (running) created at:
  ase/pq-1000.TestRace_pq_1000()
      /work/verified_test.go:15 +0x144
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 11 (finished) created at:
  ase/pq-1000.TestRace_pq_1000()
      /work/verified_test.go:21 +0x7e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0001a817f by goroutine 14:
  ase/pq-1000.(*conn).setBad()
      /work/conn.go:9 +0xa4
  ase/pq-1000.TestRace_pq_1000.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-pq-1000-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-pq-1000-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-pq-1000-fix .
docker run --rm --memory=2g --cpus=1 gonb-pq-1000-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-pq-1000-bug .
# (then run as above, no --ssh flag)
```
