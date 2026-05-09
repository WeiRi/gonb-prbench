# kubernetes-105734

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/105734 |
| Bug commit | `9248f27e2368` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `httplog.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0001a5320 by goroutine 133:
  ase/kubernetes-105734.(*respLogger).Log()
      /work/httplog.go:34 +0x1cf
  ase/kubernetes-105734.TestRace_105734_HttplogConcurrent.func3()
      /work/verified_test.go:29 +0x151

Previous write at 0x00c0001a5320 by goroutine 131:
  ase/kubernetes-105734.(*respLogger).AddKeyValue()
      /work/httplog.go:28 +0x23a
  ase/kubernetes-105734.TestRace_105734_HttplogConcurrent.func1()
      /work/verified_test.go:17 +0xd0

Goroutine 133 (running) created at:
  ase/kubernetes-105734.TestRace_105734_HttplogConcurrent()
      /work/verified_test.go:26 +0x56
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 131 (finished) created at:
  ase/kubernetes-105734.TestRace_105734_HttplogConcurrent()
      /work/verified_test.go:14 +0x1e4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000113400 by goroutine 13:
  runtime.slicecopy()
      /usr/local/go/src/runtime/slice.go:310 +0x0
  ase/kubernetes-105734.(*respLogger).Log()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-105734-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-105734-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-105734-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-105734-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-105734-bug .
# (then run as above, no --ssh flag)
```
