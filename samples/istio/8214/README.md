# istio-8214

| Field | Value |
|---|---|
| Project | istio |
| Reference | https://github.com/istio/istio/pull/8214 |
| Bug commit | `96649f851ded` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/cache/lruCache.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c420014060 by goroutine 58:
  _/work.TestRace.func2()
      /work/verified_test.go:29 +0x4a

Previous write at 0x00c420014060 by goroutine 57:
  sync/atomic.AddInt64()
      /usr/local/go/src/runtime/race_amd64.s:276 +0xb
  _/work.(*lruCache).SetWithExpiration()
      /work/verified_test.go:25 +0x43
  _/work.TestRace.func1()
      /work/verified_test.go:40 +0x41

Goroutine 58 (running) created at:
  _/work.TestRace()
      /work/verified_test.go:48 +0x104
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:777 +0x16d

Goroutine 57 (running) created at:
  _/work.TestRace()
      /work/verified_test.go:38 +0xc0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:777 +0x16d
==================
--- FAIL: TestRace (0.02s)
	testing.go:730: race detected during execution of test
FAIL
FAIL	_/work	0.036s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-istio-8214-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-istio-8214-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-istio-8214-fix .
docker run --rm --memory=2g --cpus=1 gonb-istio-8214-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-istio-8214-bug .
# (then run as above, no --ssh flag)
```
