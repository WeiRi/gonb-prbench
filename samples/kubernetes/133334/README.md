# kubernetes-133334

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/133334 |
| Bug commit | `032142c53e54` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/scheduler/framework/api_calls/pod_status_patch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Write at 0x00c00030b0c0 by goroutine 63:
  k8s.io/kubernetes/pkg/api/v1/pod.UpdatePodCondition()
      /k8s/pkg/api/v1/pod/util.go:364 +0x70
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.syncStatus()
      /k8s/pkg/scheduler/framework/api_calls/pod_status_patch.go:76 +0xe4
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.(*PodStatusPatchCall).Sync()
      /k8s/pkg/scheduler/framework/api_calls/pod_status_patch.go:134 +0x244
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.TestPodStatusPatchCall_Sync_ConcurrentNewConditionRace_133334.func1()
      /k8s/pkg/scheduler/framework/api_calls/133334_handcrafted_race_test.go:47 +0xcf

Previous write at 0x00c00030b0c0 by goroutine 59:
  k8s.io/kubernetes/pkg/api/v1/pod.UpdatePodCondition()
      /k8s/pkg/api/v1/pod/util.go:364 +0x70
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.syncStatus()
      /k8s/pkg/scheduler/framework/api_calls/pod_status_patch.go:76 +0xe4
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.(*PodStatusPatchCall).Sync()
      /k8s/pkg/scheduler/framework/api_calls/pod_status_patch.go:134 +0x244
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.TestPodStatusPatchCall_Sync_ConcurrentNewConditionRace_133334.func1()
      /k8s/pkg/scheduler/framework/api_calls/133334_handcrafted_race_test.go:47 +0xcf

Goroutine 63 (running) created at:
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.TestPodStatusPatchCall_Sync_ConcurrentNewConditionRace_133334()
      /k8s/pkg/scheduler/framework/api_calls/133334_handcrafted_race_test.go:42 +0x344
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44

Goroutine 59 (finished) created at:
  k8s.io/kubernetes/pkg/scheduler/framework/api_calls.TestPodStatusPatchCall_Sync_ConcurrentNewConditionRace_133334()
      /k8s/pkg/scheduler/framework/api_calls/133334_handcrafted_race_test.go:42 +0x344
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-133334-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133334-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-133334-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133334-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-133334-bug .
# (then run as above, no --ssh flag)
```
