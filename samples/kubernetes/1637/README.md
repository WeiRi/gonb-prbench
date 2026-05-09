# kubernetes-1637

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/1637 |
| Bug commit | `9c0fafee39d2` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/runonce.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00019a000 by goroutine 52:
  ase/kubernetes-1637.(*Kubelet).runOnce.func1()
      /work/runonce.go:39 +0x4d

Previous write at 0x00c00019a000 by goroutine 12:
  ase/kubernetes-1637.(*Kubelet).runOnce()
      /work/runonce.go:37 +0xf6
  ase/kubernetes-1637.TestRace_PR1637_RunOnceLoopVarCapture.func1()
      /work/verified_test.go:73 +0xc5

Goroutine 52 (running) created at:
  ase/kubernetes-1637.(*Kubelet).runOnce()
      /work/runonce.go:38 +0x87
  ase/kubernetes-1637.TestRace_PR1637_RunOnceLoopVarCapture.func1()
      /work/verified_test.go:73 +0xc5

Goroutine 12 (running) created at:
  ase/kubernetes-1637.TestRace_PR1637_RunOnceLoopVarCapture()
      /work/verified_test.go:70 +0x2aa
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR1637_RunOnceLoopVarCapture (0.03s)
FAIL
FAIL	ase/kubernetes-1637	0.135s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-1637-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-1637-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-1637-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-1637-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-1637-bug .
# (then run as above, no --ssh flag)
```
