# istio-54477

| Field | Value |
|---|---|
| Project | istio |
| Reference | https://github.com/istio/istio/pull/54477 |
| Bug commit | `733be18a30c9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pilot/pkg/bootstrap/ads.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000014228 by goroutine 28:
  ase/istio-54477.(*T_54477).write()
      /work/verified_test.go:6 +0xa7
  ase/istio-54477.TestRace_54477.func1()
      /work/verified_test.go:13 +0x8b

Previous write at 0x00c000014228 by goroutine 8:
  ase/istio-54477.(*T_54477).write()
      /work/verified_test.go:6 +0xa7
  ase/istio-54477.TestRace_54477.func1()
      /work/verified_test.go:13 +0x8b

Goroutine 28 (running) created at:
  ase/istio-54477.TestRace_54477()
      /work/verified_test.go:12 +0x84
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/istio-54477.TestRace_54477()
      /work/verified_test.go:12 +0x84
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000014228 by goroutine 64:
  ase/istio-54477.(*T_54477).read()
      /work/verified_test.go:7 +0xa4
  ase/istio-54477.TestRace_54477.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-istio-54477-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-istio-54477-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-istio-54477-fix .
docker run --rm --memory=2g --cpus=1 gonb-istio-54477-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-istio-54477-bug .
# (then run as above, no --ssh flag)
```
