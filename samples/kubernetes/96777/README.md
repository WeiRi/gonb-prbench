# kubernetes-96777

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/96777 |
| Bug commit | `b2ecd1b3a319` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/scheduler/framework/plugins/volumebinding/volume_binding.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0001897a0 by goroutine 63:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  k8s.io/kubernetes/pkg/scheduler/framework/plugins/volumebinding.TestRace_96777.func1()
      /work/pkg/scheduler/framework/plugins/volumebinding/race_96777_capture_v2_test.go:25 +0x139

Previous write at 0x00c0001897a0 by goroutine 64:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  k8s.io/kubernetes/pkg/scheduler/framework/plugins/volumebinding.TestRace_96777.func2()
      /work/pkg/scheduler/framework/plugins/volumebinding/race_96777_capture_v2_test.go:31 +0x139

Goroutine 63 (running) created at:
  k8s.io/kubernetes/pkg/scheduler/framework/plugins/volumebinding.TestRace_96777()
      /work/pkg/scheduler/framework/plugins/volumebinding/race_96777_capture_v2_test.go:22 +0x1c4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 64 (running) created at:
  k8s.io/kubernetes/pkg/scheduler/framework/plugins/volumebinding.TestRace_96777()
      /work/pkg/scheduler/framework/plugins/volumebinding/race_96777_capture_v2_test.go:28 +0x104
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0001897a0 by goroutine 66:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  k8s.io/kubernetes/pkg/scheduler/framework/plugins/volumebinding.TestRace_96777.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-96777-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-96777-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-96777-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-96777-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-96777-bug .
# (then run as above, no --ssh flag)
```
