# kubernetes-137804

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/137804 |
| Bug commit | `7d56731021ec` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/scheduler/backend/cache/podgroupstate.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Write at 0x00c0002e4158 by goroutine 53:
  k8s.io/kubernetes/pkg/scheduler/backend/cache.(*podGroupStateData).addPod()
      /k8s/pkg/scheduler/backend/cache/podgroupstate.go:94 +0x5b
  k8s.io/kubernetes/pkg/scheduler/backend/cache.TestPodGroupState_RaceOnEmptyAndForgetPod_137804.func3()
      /k8s/pkg/scheduler/backend/cache/137804_handcrafted_race_test.go:90 +0x25a

Previous write at 0x00c0002e4158 by goroutine 52:
  k8s.io/kubernetes/pkg/scheduler/backend/cache.(*podGroupStateData).assumePod()
      /k8s/pkg/scheduler/backend/cache/podgroupstate.go:134 +0xf7
  k8s.io/kubernetes/pkg/scheduler/backend/cache.TestPodGroupState_RaceOnEmptyAndForgetPod_137804.func2()
      /k8s/pkg/scheduler/backend/cache/137804_handcrafted_race_test.go:74 +0xde

Goroutine 53 (running) created at:
  k8s.io/kubernetes/pkg/scheduler/backend/cache.TestPodGroupState_RaceOnEmptyAndForgetPod_137804()
      /k8s/pkg/scheduler/backend/cache/137804_handcrafted_race_test.go:86 +0x619
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:2036 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:2101 +0x38

Goroutine 52 (running) created at:
  k8s.io/kubernetes/pkg/scheduler/backend/cache.TestPodGroupState_RaceOnEmptyAndForgetPod_137804()
      /k8s/pkg/scheduler/backend/cache/137804_handcrafted_race_test.go:69 +0x55c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:2036 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:2101 +0x38
==================
==================
WARNING: DATA RACE
Write at 0x00c00037a6f0 by goroutine 53:
  runtime.mapassign_faststr()
      /usr/local/go/src/internal/runtime/maps/runtime_faststr.go:263 +0x0
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-137804-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-137804-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-137804-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-137804-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-137804-bug .
# (then run as above, no --ssh flag)
```
