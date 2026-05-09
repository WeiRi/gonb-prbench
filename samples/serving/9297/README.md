# serving-9297

| Field | Value |
|---|---|
| Project | serving |
| Reference | https://github.com/knative/serving/pull/9297 |
| Bug commit | `45e39ea4a80c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/autoscaler/statforwarder/forwarder.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0005fe2a0 by goroutine 57:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  knative.dev/serving/pkg/autoscaler/statforwarder.TestForwarderCancelRace.func2()
      /work/pkg/autoscaler/statforwarder/race_9297_capture_v2_test.go:41 +0x3ac

Previous read at 0x00c0005fe2a0 by goroutine 56:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  knative.dev/serving/pkg/autoscaler/statforwarder.(*Forwarder).Cancel()
      /work/pkg/autoscaler/statforwarder/forwarder.go:347 +0xe4
  knative.dev/serving/pkg/autoscaler/statforwarder.TestForwarderCancelRace.func1()
      /work/pkg/autoscaler/statforwarder/race_9297_capture_v2_test.go:36 +0x95

Goroutine 57 (running) created at:
  knative.dev/serving/pkg/autoscaler/statforwarder.TestForwarderCancelRace()
      /work/pkg/autoscaler/statforwarder/race_9297_capture_v2_test.go:38 +0x36
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 56 (finished) created at:
  knative.dev/serving/pkg/autoscaler/statforwarder.TestForwarderCancelRace()
      /work/pkg/autoscaler/statforwarder/race_9297_capture_v2_test.go:34 +0x535
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0005fe2d0 by goroutine 60:
  runtime.mapiternext()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-serving-9297-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-serving-9297-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-serving-9297-fix .
docker run --rm --memory=2g --cpus=1 gonb-serving-9297-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-serving-9297-bug .
# (then run as above, no --ssh flag)
```
