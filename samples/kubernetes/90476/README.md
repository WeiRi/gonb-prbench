# kubernetes-90476

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/90476 |
| Bug commit | `4efb43e1f574` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0001b4108 by goroutine 15:
  k8s.io/apimachinery/pkg/util/wait.(*Backoff).Step()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go:244 +0x3c
  k8s.io/apimachinery/pkg/util/wait.(*exponentialBackoffManagerImpl).getNextBackoff()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go:334 +0x1ca
  k8s.io/apimachinery/pkg/util/wait.(*exponentialBackoffManagerImpl).Backoff()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go:339 +0x2e
  k8s.io/apimachinery/pkg/util/wait.TestRace_90476.func1()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/race_90476_capture_test.go:35 +0xb1

Previous write at 0x00c0001b4108 by goroutine 10:
  k8s.io/apimachinery/pkg/util/wait.(*Backoff).Step()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go:250 +0x13c
  k8s.io/apimachinery/pkg/util/wait.(*exponentialBackoffManagerImpl).getNextBackoff()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go:334 +0x1ca
  k8s.io/apimachinery/pkg/util/wait.(*exponentialBackoffManagerImpl).Backoff()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/wait.go:339 +0x2e
  k8s.io/apimachinery/pkg/util/wait.TestRace_90476.func1()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/race_90476_capture_test.go:35 +0xb1

Goroutine 15 (running) created at:
  k8s.io/apimachinery/pkg/util/wait.TestRace_90476()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/race_90476_capture_test.go:32 +0x9d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 10 (running) created at:
  k8s.io/apimachinery/pkg/util/wait.TestRace_90476()
      /work/staging/src/k8s.io/apimachinery/pkg/util/wait/race_90476_capture_test.go:32 +0x9d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-90476-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-90476-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-90476-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-90476-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-90476-bug .
# (then run as above, no --ssh flag)
```
