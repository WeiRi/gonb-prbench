# kubernetes-136205

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/136205 |
| Bug commit | `62277ef5d29d` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/status/status_manager.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Read at 0x00c0001fa198 by goroutine 22:
  k8s.io/kubernetes/pkg/kubelet/status.TestUpdateStatusInternal_StartTimeSharedPointerRace_136205.func2()
      /k8s/pkg/kubelet/status/136205_handcrafted_race_test.go:72 +0xa4

Previous write at 0x00c0001fa198 by goroutine 21:
  k8s.io/kubernetes/pkg/kubelet/status.normalizeStatus.func1()
      /k8s/pkg/kubelet/status/status_manager.go:1193 +0x1ee
  k8s.io/kubernetes/pkg/kubelet/status.normalizeStatus()
      /k8s/pkg/kubelet/status/status_manager.go:1209 +0x175
  k8s.io/kubernetes/pkg/kubelet/status.(*manager).updateStatusInternal()
      /k8s/pkg/kubelet/status/status_manager.go:838 +0x9dc
  k8s.io/kubernetes/pkg/kubelet/status.TestUpdateStatusInternal_StartTimeSharedPointerRace_136205.func1()
      /k8s/pkg/kubelet/status/136205_handcrafted_race_test.go:61 +0x1c5

Goroutine 22 (running) created at:
  k8s.io/kubernetes/pkg/kubelet/status.TestUpdateStatusInternal_StartTimeSharedPointerRace_136205()
      /k8s/pkg/kubelet/status/136205_handcrafted_race_test.go:69 +0x72c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 21 (running) created at:
  k8s.io/kubernetes/pkg/kubelet/status.TestUpdateStatusInternal_StartTimeSharedPointerRace_136205()
      /k8s/pkg/kubelet/status/136205_handcrafted_race_test.go:53 +0x684
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-136205-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136205-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-136205-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136205-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-136205-bug .
# (then run as above, no --ssh flag)
```
