# kubernetes-114236

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/114236 |
| Bug commit | `3e26e104bdf9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/tools/events/event_broadcaster.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000206440 by goroutine 9:
  ase/kubernetes-114236.(*eventBroadcasterImpl).recordToSink()
      /work/event_broadcaster.go:47 +0xd1
  ase/kubernetes-114236.TestRace_114236_EventBroadcasterSharedCache.func2()
      /work/verified_test.go:29 +0x1c4

Previous write at 0x00c000206440 by goroutine 8:
  ase/kubernetes-114236.attemptRecording()
      /work/event_broadcaster.go:60 +0xee
  ase/kubernetes-114236.TestRace_114236_EventBroadcasterSharedCache.func1()
      /work/verified_test.go:22 +0x1d0

Goroutine 9 (running) created at:
  ase/kubernetes-114236.TestRace_114236_EventBroadcasterSharedCache()
      /work/verified_test.go:26 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 8 (running) created at:
  ase/kubernetes-114236.TestRace_114236_EventBroadcasterSharedCache()
      /work/verified_test.go:18 +0x2bc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000206430 by goroutine 9:
  ase/kubernetes-114236.attemptRecording()
      /work/event_broadcaster.go:59 +0x1d8
  ase/kubernetes-114236.TestRace_114236_EventBroadcasterSharedCache.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-114236-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-114236-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-114236-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-114236-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-114236-bug .
# (then run as above, no --ssh flag)
```
