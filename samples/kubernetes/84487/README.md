# kubernetes-84487

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/84487 |
| Bug commit | `b6c8f4916dc3` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/registry/registrytest/node.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000282c80 by goroutine 203:
  k8s.io/kubernetes/pkg/registry/registrytest.(*NodeRegistry).WatchNodes()
      /work/pkg/registry/registrytest/node.go:117 +0x14b
  k8s.io/kubernetes/pkg/registry/registrytest.TestRace_84487.func1()
      /work/pkg/registry/registrytest/race_84487_capture_v2_test.go:28 +0xc6

Previous write at 0x00c000282c80 by goroutine 204:
  k8s.io/kubernetes/pkg/registry/registrytest.TestRace_84487.func2()
      /work/pkg/registry/registrytest/race_84487_capture_v2_test.go:38 +0xbe

Goroutine 203 (running) created at:
  k8s.io/kubernetes/pkg/registry/registrytest.TestRace_84487()
      /work/pkg/registry/registrytest/race_84487_capture_v2_test.go:25 +0x1a9
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 204 (finished) created at:
  k8s.io/kubernetes/pkg/registry/registrytest.TestRace_84487()
      /work/pkg/registry/registrytest/race_84487_capture_v2_test.go:34 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRace_84487 (0.00s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	k8s.io/kubernetes/pkg/registry/registrytest	0.285s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-84487-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-84487-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-84487-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-84487-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-84487-bug .
# (then run as above, no --ssh flag)
```
