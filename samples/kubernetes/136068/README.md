# kubernetes-136068

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/136068 |
| Bug commit | `fe36b79c2ab5` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/tools/leaderelection/leaderelection.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0003c7db0 by goroutine 111:
  k8s.io/client-go/tools/leaderelection.(*LeaderElector).setObservedRecord()
      /work/staging/src/k8s.io/client-go/tools/leaderelection/leaderelection.go:533 +0x171
  k8s.io/client-go/tools/leaderelection.TestRace_136068.func2()
      /work/staging/src/k8s.io/client-go/tools/leaderelection/race_136068_capture_test.go:43 +0x14f

Previous read at 0x00c0003c7db0 by goroutine 40:
  k8s.io/client-go/tools/leaderelection.TestRace_136068.func1()
      /work/staging/src/k8s.io/client-go/tools/leaderelection/race_136068_capture_test.go:28 +0xa6

Goroutine 111 (running) created at:
  k8s.io/client-go/tools/leaderelection.TestRace_136068()
      /work/staging/src/k8s.io/client-go/tools/leaderelection/race_136068_capture_test.go:36 +0x311
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 40 (finished) created at:
  k8s.io/client-go/tools/leaderelection.TestRace_136068()
      /work/staging/src/k8s.io/client-go/tools/leaderelection/race_136068_capture_test.go:25 +0x22f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
--- FAIL: TestRace_136068 (0.02s)
    testing.go:1617: race detected during execution of test
FAIL
FAIL	k8s.io/client-go/tools/leaderelection	0.103s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-136068-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136068-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-136068-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136068-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-136068-bug .
# (then run as above, no --ssh flag)
```
