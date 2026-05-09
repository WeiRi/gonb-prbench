# kubernetes-92139

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/92139 |
| Bug commit | `1fa20301a017` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/tools/clientcmd/merged_client_builder.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000025fb8 by goroutine 24:
  k8s.io/client-go/tools/clientcmd.TestRace_92139.func1()
      /work/staging/src/k8s.io/client-go/tools/clientcmd/race_92139_capture_test.go:63 +0xc6

Previous read at 0x00c000025fb8 by goroutine 13:
  k8s.io/client-go/tools/clientcmd.(*DeferredLoadingClientConfig).createClientConfig()
      /work/staging/src/k8s.io/client-go/tools/clientcmd/merged_client_builder.go:67 +0x104
  k8s.io/client-go/tools/clientcmd.TestRace_92139.func1()
      /work/staging/src/k8s.io/client-go/tools/clientcmd/race_92139_capture_test.go:64 +0x9e

Goroutine 24 (running) created at:
  k8s.io/client-go/tools/clientcmd.TestRace_92139()
      /work/staging/src/k8s.io/client-go/tools/clientcmd/race_92139_capture_test.go:57 +0x185
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 13 (running) created at:
  k8s.io/client-go/tools/clientcmd.TestRace_92139()
      /work/staging/src/k8s.io/client-go/tools/clientcmd/race_92139_capture_test.go:57 +0x185
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000025fb8 by goroutine 31:
  k8s.io/client-go/tools/clientcmd.TestRace_92139.func1()
      /work/staging/src/k8s.io/client-go/tools/clientcmd/race_92139_capture_test.go:63 +0xc6

Previous write at 0x00c000025fb8 by goroutine 13:
  k8s.io/client-go/tools/clientcmd.(*DeferredLoadingClientConfig).createClientConfig()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-92139-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-92139-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-92139-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-92139-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-92139-bug .
# (then run as above, no --ssh flag)
```
