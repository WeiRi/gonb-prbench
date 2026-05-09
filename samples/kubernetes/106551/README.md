# kubernetes-106551

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/106551 |
| Bug commit | `7f920da4420b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/config/config.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00057d6b0 by goroutine 54:
  k8s.io/apimachinery/pkg/util/sets.String.List()
      /work/staging/src/k8s.io/apimachinery/pkg/util/sets/string.go:171 +0x44
  k8s.io/kubernetes/pkg/kubelet/config.(*PodConfig).SeenAllSources()
      /work/pkg/kubelet/config/config.go:98 +0x9d
  k8s.io/kubernetes/pkg/kubelet/config.TestRace106551SeenAllSources.func2()
      /work/pkg/kubelet/config/race_106551_capture_test.go:39 +0x97

Previous write at 0x00c00057d6b0 by goroutine 149:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  k8s.io/apimachinery/pkg/util/sets.String.Insert()
      /work/staging/src/k8s.io/apimachinery/pkg/util/sets/string.go:51 +0x1b1
  k8s.io/kubernetes/pkg/kubelet/config.(*PodConfig).Channel()
      /work/pkg/kubelet/config/config.go:88 +0xd3
  k8s.io/kubernetes/pkg/kubelet/config.TestRace106551SeenAllSources.func1()
      /work/pkg/kubelet/config/race_106551_capture_test.go:33 +0x11a
  k8s.io/kubernetes/pkg/kubelet/config.TestRace106551SeenAllSources.func3()
      /work/pkg/kubelet/config/race_106551_capture_test.go:34 +0x41

Goroutine 54 (running) created at:
  k8s.io/kubernetes/pkg/kubelet/config.TestRace106551SeenAllSources()
      /work/pkg/kubelet/config/race_106551_capture_test.go:35 +0x638
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 149 (finished) created at:
  k8s.io/kubernetes/pkg/kubelet/config.TestRace106551SeenAllSources()
      /work/pkg/kubelet/config/race_106551_capture_test.go:28 +0x7a7
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-106551-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-106551-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-106551-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-106551-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-106551-bug .
# (then run as above, no --ssh flag)
```
