# kubernetes-107452

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/107452 |
| Bug commit | `d1f559711de7` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apiserver/pkg/server/filters/timeout.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000188000 by goroutine 22:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/kubernetes-107452.(*baseTimeoutWriter).WriteHeaderTimeout()
      /work/timeout.go:23 +0x89
  ase/kubernetes-107452.TestRace_107452_HeaderMutation.func1.2()
      /work/verified_test.go:52 +0x84

Previous write at 0x00c000188000 by goroutine 20:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  net/textproto.MIMEHeader.Set()
      /usr/local/go/src/net/textproto/header.go:22 +0x171
  net/http.Header.Set()
      /usr/local/go/src/net/http/header.go:40 +0xf4
  ase/kubernetes-107452.TestRace_107452_HeaderMutation.func1.1()
      /work/verified_test.go:44 +0xf3

Goroutine 22 (running) created at:
  ase/kubernetes-107452.TestRace_107452_HeaderMutation.func1()
      /work/verified_test.go:50 +0xa9

Goroutine 20 (finished) created at:
  ase/kubernetes-107452.TestRace_107452_HeaderMutation.func1()
      /work/verified_test.go:41 +0x30a
==================
==================
WARNING: DATA RACE
Read at 0x00c000196088 by goroutine 22:
  ase/kubernetes-107452.(*baseTimeoutWriter).WriteHeaderTimeout()
      /work/timeout.go:23 +0xd5
  ase/kubernetes-107452.TestRace_107452_HeaderMutation.func1.2()
      /work/verified_test.go:52 +0x84
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-107452-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-107452-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-107452-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-107452-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-107452-bug .
# (then run as above, no --ssh flag)
```
