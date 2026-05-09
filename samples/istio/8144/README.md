# istio-8144

| Field | Value |
|---|---|
| Project | istio |
| Reference | https://github.com/istio/istio/pull/8144 |
| Bug commit | `1f3f1780d701` |
| Category | order_violation |
| Oracle | RACE |
| Primary diff file | `pkg/cache/ttlCache.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c4200e4060 by goroutine 58:
  _/work.TestRace.func2()
      /work/verified_test.go:26 +0x67

Previous write at 0x00c4200e4060 by goroutine 57:
  sync/atomic.AddInt64()
      /usr/local/go/src/runtime/race_amd64.s:276 +0xb
  _/work.(*ttlCache).IncrementEvictions()
      /work/verified_test.go:38 +0x43
  _/work.TestRace.func1()
      /work/verified_test.go:51 +0x61

Goroutine 58 (running) created at:
  _/work.TestRace()
      /work/verified_test.go:59 +0x108
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:777 +0x16d

Goroutine 57 (finished) created at:
  _/work.TestRace()
      /work/verified_test.go:47 +0xc4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:777 +0x16d
==================
==================
WARNING: DATA RACE
Read at 0x00c4200e4068 by goroutine 58:
  _/work.TestRace.func2()
      /work/verified_test.go:26 +0x67

Previous write at 0x00c4200e4068 by goroutine 57:
  sync/atomic.AddInt64()
      /usr/local/go/src/runtime/race_amd64.s:276 +0xb
  _/work.(*ttlCache).IncrementHits()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-istio-8144-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-istio-8144-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-istio-8144-fix .
docker run --rm --memory=2g --cpus=1 gonb-istio-8144-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-istio-8144-bug .
# (then run as above, no --ssh flag)
```
