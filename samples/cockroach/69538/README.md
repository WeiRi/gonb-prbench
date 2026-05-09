# cockroach-69538

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/69538 |
| Bug commit | `3d202590e356` |
| Category | special_library |
| Oracle | RACE |
| Primary diff file | `pkg/util/metric/registry.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d9020 by goroutine 39:
  runtime.mapassign_faststr()
      /usr/local/go/src/internal/runtime/maps/runtime_faststr_swiss.go:263 +0x0
  ase/cockroach-69538.(*Registry).AddMetric()
      /work/registry.go:16 +0xa9
  ase/cockroach-69538.TestRace_69538.func2()
      /work/verified_test.go:37 +0x84
  ase/cockroach-69538.TestRace_69538.gowrap1()
      /work/verified_test.go:38 +0x41

Previous read at 0x00c0000d9020 by goroutine 8:
  runtime.mapIterStart()
      /usr/local/go/src/runtime/map_swiss.go:160 +0x0
  ase/cockroach-69538.(*Registry).Each()
      /work/registry.go:21 +0xe4
  ase/cockroach-69538.TestRace_69538.func1()
      /work/verified_test.go:28 +0x9c

Goroutine 39 (running) created at:
  ase/cockroach-69538.TestRace_69538()
      /work/verified_test.go:35 +0x1c8
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 8 (finished) created at:
  ase/cockroach-69538.TestRace_69538()
      /work/verified_test.go:26 +0xe5
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-69538-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-69538-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-69538-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-69538-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-69538-bug .
# (then run as above, no --ssh flag)
```
