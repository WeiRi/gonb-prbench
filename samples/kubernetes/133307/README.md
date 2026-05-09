# kubernetes-133307

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/133307 |
| Bug commit | `8e6d78888703` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/component-base/metrics/counter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Write at 0x00c0000e0dc0 by goroutine 13:
  k8s.io/component-base/metrics.(*Counter).WithContext()
      /k8s/staging/src/k8s.io/component-base/metrics/counter.go:110 +0x127
  k8s.io/component-base/metrics.TestCounterWithContextRace_133307.func1()
      /k8s/staging/src/k8s.io/component-base/metrics/133307_handcrafted_race_test.go:51 +0x122

Previous write at 0x00c0000e0dc0 by goroutine 11:
  k8s.io/component-base/metrics.(*Counter).WithContext()
      /k8s/staging/src/k8s.io/component-base/metrics/counter.go:110 +0x1af
  k8s.io/component-base/metrics.TestCounterWithContextRace_133307.func1()
      /k8s/staging/src/k8s.io/component-base/metrics/133307_handcrafted_race_test.go:53 +0x1a5

Goroutine 13 (running) created at:
  k8s.io/component-base/metrics.TestCounterWithContextRace_133307()
      /k8s/staging/src/k8s.io/component-base/metrics/133307_handcrafted_race_test.go:48 +0x144
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44

Goroutine 11 (running) created at:
  k8s.io/component-base/metrics.TestCounterWithContextRace_133307()
      /k8s/staging/src/k8s.io/component-base/metrics/133307_handcrafted_race_test.go:48 +0x144
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0000e0dc0 by goroutine 13:
  k8s.io/component-base/metrics.(*exemplarCounterMetric).withExemplar()
      /k8s/staging/src/k8s.io/component-base/metrics/counter.go:130 +0xd7
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-133307-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133307-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-133307-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133307-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-133307-bug .
# (then run as above, no --ssh flag)
```
