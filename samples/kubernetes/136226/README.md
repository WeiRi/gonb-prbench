# kubernetes-136226

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/136226 |
| Bug commit | `c086a712b1f1` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/kubelet_pods.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Read at 0x00c000713b08 by goroutine 103:
  k8s.io/component-helpers/resource.IsPodLevelRequestsSet()
      /k8s/staging/src/k8s.io/component-helpers/resource/helpers.go:111 +0x94
  k8s.io/component-helpers/resource.PodRequests()
      /k8s/staging/src/k8s.io/component-helpers/resource/helpers.go:155 +0x124
  k8s.io/kubernetes/pkg/kubelet.getEffectiveAllocatedResources()
      /k8s/pkg/kubelet/kubelet_pods.go:2164 +0x134
  k8s.io/kubernetes/pkg/kubelet.TestGetEffectiveAllocatedResources_NoSharedAlias_136226.func1()
      /k8s/pkg/kubelet/136226_handcrafted_race_test.go:61 +0x99

Previous write at 0x00c000713b08 by goroutine 102:
  k8s.io/kubernetes/pkg/kubelet.getEffectiveAllocatedResources()
      /k8s/pkg/kubelet/kubelet_pods.go:2164 +0x149
  k8s.io/kubernetes/pkg/kubelet.TestGetEffectiveAllocatedResources_NoSharedAlias_136226.func1()
      /k8s/pkg/kubelet/136226_handcrafted_race_test.go:61 +0x99

Goroutine 103 (running) created at:
  k8s.io/kubernetes/pkg/kubelet.TestGetEffectiveAllocatedResources_NoSharedAlias_136226()
      /k8s/pkg/kubelet/136226_handcrafted_race_test.go:58 +0xb6a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 102 (running) created at:
  k8s.io/kubernetes/pkg/kubelet.TestGetEffectiveAllocatedResources_NoSharedAlias_136226()
      /k8s/pkg/kubelet/136226_handcrafted_race_test.go:58 +0xb6a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-136226-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136226-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-136226-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136226-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-136226-bug .
# (then run as above, no --ssh flag)
```
