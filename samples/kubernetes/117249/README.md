# kubernetes-117249

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/117249 |
| Bug commit | `83a1774df2bc` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/controller/endpointslice/topologycache/topologycache.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0007879f8 by goroutine 121:
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.(*TopologyCache).SetNodes()
      /work/pkg/controller/endpointslice/topologycache/topologycache.go:257 +0xd56
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.TestRace117249TopologyCache.func1()
      /work/pkg/controller/endpointslice/topologycache/race_117249_capture_v2_test.go:61 +0xbc

Previous read at 0x00c0007879f8 by goroutine 122:
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.(*TopologyCache).getAllocations()
      /work/pkg/controller/endpointslice/topologycache/topologycache.go:274 +0x6f
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.(*TopologyCache).AddHints()
      /work/pkg/controller/endpointslice/topologycache/topologycache.go:93 +0x56
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.TestRace117249TopologyCache.func2()
      /work/pkg/controller/endpointslice/topologycache/race_117249_capture_v2_test.go:66 +0x8d

Goroutine 121 (running) created at:
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.TestRace117249TopologyCache()
      /work/pkg/controller/endpointslice/topologycache/race_117249_capture_v2_test.go:58 +0x10b8
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 122 (running) created at:
  k8s.io/kubernetes/pkg/controller/endpointslice/topologycache.TestRace117249TopologyCache()
      /work/pkg/controller/endpointslice/topologycache/race_117249_capture_v2_test.go:63 +0xf86
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
I0507 16:45:05.164071   11707 topologycache.go:96] "Nodes only ready in one zone, removing hints" serviceKey="ns/svc" addressType=IPv4
I0507 16:45:05.168028   11707 topologycache.go:96] "Nodes only ready in one zone, removing hints" serviceKey="ns/svc" addressType=IPv4
==================
WARNING: DATA RACE
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-117249-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-117249-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-117249-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-117249-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-117249-bug .
# (then run as above, no --ssh flag)
```
