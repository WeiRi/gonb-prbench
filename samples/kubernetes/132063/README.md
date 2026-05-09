# kubernetes-132063

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/132063 |
| Bug commit | `45d267ca164e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/component-base/metrics/counter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Read at 0x00c0001a2db0 by goroutine 43:
  k8s.io/component-base/metrics.(*CounterVec).WithLabelValues()
      /k8s/staging/src/k8s.io/component-base/metrics/counter.go:215 +0xa5
  k8s.io/component-base/metrics.TestCounterVecWithLabelValues_LabelAllowListRace_132063.func1()
      /k8s/staging/src/k8s.io/component-base/metrics/132063_handcrafted_race_test.go:47 +0xfc

Previous write at 0x00c0001a2db0 by goroutine 69:
  k8s.io/component-base/metrics.(*CounterVec).WithLabelValues.func1()
      /k8s/staging/src/k8s.io/component-base/metrics/counter.go:221 +0x104
  sync.(*Once).doSlow()
      /usr/local/go/src/sync/once.go:78 +0xe1
  sync.(*Once).Do()
      /usr/local/go/src/sync/once.go:69 +0x44
  k8s.io/component-base/metrics.(*CounterVec).WithLabelValues()
      /k8s/staging/src/k8s.io/component-base/metrics/counter.go:218 +0x1eb
  k8s.io/component-base/metrics.TestCounterVecWithLabelValues_LabelAllowListRace_132063.func1()
      /k8s/staging/src/k8s.io/component-base/metrics/132063_handcrafted_race_test.go:47 +0xfc

Goroutine 43 (running) created at:
  k8s.io/component-base/metrics.TestCounterVecWithLabelValues_LabelAllowListRace_132063()
      /k8s/staging/src/k8s.io/component-base/metrics/132063_handcrafted_race_test.go:43 +0x3fb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44

Goroutine 69 (running) created at:
  k8s.io/component-base/metrics.TestCounterVecWithLabelValues_LabelAllowListRace_132063()
      /k8s/staging/src/k8s.io/component-base/metrics/132063_handcrafted_race_test.go:43 +0x3fb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-132063-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-132063-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-132063-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-132063-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-132063-bug .
# (then run as above, no --ssh flag)
```
