# serving-4973

| Field | Value |
|---|---|
| Project | serving |
| Reference | https://github.com/knative/serving/pull/4973 |
| Bug commit | `6110a2d4418a` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/activator/handler/concurrency_reporter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000d6510 by goroutine 8:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/serving-4973.(*ConcurrencyReporter).HandleRequest()
      /work/concurrency_reporter.go:18 +0xdc
  ase/serving-4973.TestRace_4973.func1()
      /work/verified_test.go:22 +0xb4
  ase/serving-4973.TestRace_4973.gowrap1()
      /work/verified_test.go:24 +0x41

Previous write at 0x00c0000d6510 by goroutine 11:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/serving-4973.(*ConcurrencyReporter).HandleRequest()
      /work/concurrency_reporter.go:18 +0x124
  ase/serving-4973.TestRace_4973.func1()
      /work/verified_test.go:22 +0xb4
  ase/serving-4973.TestRace_4973.gowrap1()
      /work/verified_test.go:24 +0x41

Goroutine 8 (running) created at:
  ase/serving-4973.TestRace_4973()
      /work/verified_test.go:19 +0x104
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 11 (finished) created at:
  ase/serving-4973.TestRace_4973()
      /work/verified_test.go:19 +0x104
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-serving-4973-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-serving-4973-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-serving-4973-fix .
docker run --rm --memory=2g --cpus=1 gonb-serving-4973-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-serving-4973-bug .
# (then run as above, no --ssh flag)
```
