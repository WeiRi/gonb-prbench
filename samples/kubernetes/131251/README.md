# kubernetes-131251

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/131251 |
| Bug commit | `b15dfce6cbd0` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/images/image_gc_manager.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Write at 0x00c000844f60 by goroutine 140:
  runtime.mapassign()
      /usr/local/go/src/internal/runtime/maps/runtime_swiss.go:191 +0x0
  k8s.io/kubernetes/pkg/kubelet/images.(*realImageGCManager).freeImage()
      /k8s/pkg/kubelet/images/image_gc_manager.go:524 +0x554
  k8s.io/kubernetes/pkg/kubelet/images.TestImageGCManager_FreeImage_LockOnDelete_131251.func2()
      /k8s/pkg/kubelet/images/131251_handcrafted_race_test.go:51 +0x159

Previous read at 0x00c000844f60 by goroutine 139:
  k8s.io/kubernetes/pkg/kubelet/images.(*realImageGCManager).imageRecordsLen()
      /k8s/pkg/kubelet/images/image_gc_manager_test.go:63 +0xb2
  k8s.io/kubernetes/pkg/kubelet/images.TestImageGCManager_FreeImage_LockOnDelete_131251.func1()
      /k8s/pkg/kubelet/images/131251_handcrafted_race_test.go:38 +0x99

Goroutine 140 (running) created at:
  k8s.io/kubernetes/pkg/kubelet/images.TestImageGCManager_FreeImage_LockOnDelete_131251()
      /k8s/pkg/kubelet/images/131251_handcrafted_race_test.go:45 +0x852
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44

Goroutine 139 (finished) created at:
  k8s.io/kubernetes/pkg/kubelet/images.TestImageGCManager_FreeImage_LockOnDelete_131251()
      /k8s/pkg/kubelet/images/131251_handcrafted_race_test.go:35 +0x792
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-131251-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-131251-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-131251-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-131251-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-131251-bug .
# (then run as above, no --ssh flag)
```
