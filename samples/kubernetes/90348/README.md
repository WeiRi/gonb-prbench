# kubernetes-90348

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/90348 |
| Bug commit | `85ee5fdd9016` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/legacy-cloud-providers/vsphere/vsphere.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00001c310 by goroutine 10:
  ase/kubernetes-90348.ProcessDisksAreAttached.func1()
      /work/vsphere.go:25 +0xd0

Previous write at 0x00c00001c310 by goroutine 7:
  ase/kubernetes-90348.ProcessDisksAreAttached()
      /work/vsphere.go:20 +0x1c4
  ase/kubernetes-90348.TestRace_90348_LoopVarCapture()
      /work/verified_test.go:14 +0x5c4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 10 (running) created at:
  ase/kubernetes-90348.ProcessDisksAreAttached()
      /work/vsphere.go:22 +0x104
  ase/kubernetes-90348.TestRace_90348_LoopVarCapture()
      /work/verified_test.go:14 +0x5c4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 7 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1742 +0x825
  testing.runTests.func1()
      /usr/local/go/src/testing/testing.go:2161 +0x85
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.runTests()
      /usr/local/go/src/testing/testing.go:2159 +0x8be
  testing.(*M).Run()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-90348-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-90348-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-90348-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-90348-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-90348-bug .
# (then run as above, no --ssh flag)
```
