# kubernetes-109114

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/109114 |
| Bug commit | `c4fdf3ded61f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/rest/with_retry.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0001ea008 by goroutine 11:
  ase/kubernetes-109114.(*withRetry).After()
      /work/with_retry.go:35 +0x190
  ase/kubernetes-109114.TestRace_109114.func1()
      /work/verified_test.go:32 +0x158

Previous write at 0x00c0001ea008 by goroutine 8:
  ase/kubernetes-109114.(*withRetry).Before()
      /work/with_retry.go:29 +0x11c
  ase/kubernetes-109114.TestRace_109114.func1()
      /work/verified_test.go:33 +0xc1

Goroutine 11 (running) created at:
  ase/kubernetes-109114.TestRace_109114()
      /work/verified_test.go:24 +0x238
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/kubernetes-109114.TestRace_109114()
      /work/verified_test.go:24 +0x238
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0001ea010 by goroutine 10:
  ase/kubernetes-109114.(*withRetry).Before()
      /work/with_retry.go:30 +0x137
  ase/kubernetes-109114.TestRace_109114.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-109114-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-109114-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-109114-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-109114-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-109114-bug .
# (then run as above, no --ssh flag)
```
