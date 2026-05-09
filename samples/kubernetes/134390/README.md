# kubernetes-134390

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/134390 |
| Bug commit | `a412a1510952` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/component-base/metrics/metric.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00008a058 by goroutine 8:
  ase/kubernetes-134390.(*Counter).ClearState()
      /work/metric.go:24 +0x5d
  ase/kubernetes-134390.TestRace_134390.func1()
      /work/verified_test.go:27 +0xa4

Previous read at 0x00c00008a058 by goroutine 9:
  ase/kubernetes-134390.(*Counter).IsHidden()
      /work/metric.go:31 +0xce
  ase/kubernetes-134390.TestRace_134390.func1()
      /work/verified_test.go:28 +0xc6

Goroutine 8 (running) created at:
  ase/kubernetes-134390.TestRace_134390()
      /work/verified_test.go:22 +0x113
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/kubernetes-134390.TestRace_134390()
      /work/verified_test.go:22 +0x113
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c00008a059 by goroutine 25:
  ase/kubernetes-134390.(*Counter).IsDeprecated()
      /work/metric.go:35 +0xd9
  ase/kubernetes-134390.TestRace_134390.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-134390-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-134390-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-134390-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-134390-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-134390-bug .
# (then run as above, no --ssh flag)
```
