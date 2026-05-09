# kubernetes-121771

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/121771 |
| Bug commit | `46f4248d56da` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apimachinery/pkg/runtime/helper.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00019e540 by goroutine 9:
  ase/kubernetes-121771.(*fakeObj).GetGroupVersionKind()
      /work/verified_test.go:68 +0x4a
  ase/kubernetes-121771.WithVersionEncoder.Encode()
      /work/helper.go:37 +0xac
  ase/kubernetes-121771.TestWithVersionEncoderRace.func1()
      /work/verified_test.go:86 +0x1f5

Previous write at 0x00c00019e540 by goroutine 8:
  ase/kubernetes-121771.(*fakeObj).SetGroupVersionKind()
      /work/verified_test.go:69 +0x4f
  ase/kubernetes-121771.WithVersionEncoder.Encode()
      /work/helper.go:40 +0x248
  ase/kubernetes-121771.TestWithVersionEncoderRace.func1()
      /work/verified_test.go:86 +0x1f5

Goroutine 9 (running) created at:
  ase/kubernetes-121771.TestWithVersionEncoderRace()
      /work/verified_test.go:82 +0x153
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/kubernetes-121771.TestWithVersionEncoderRace()
      /work/verified_test.go:82 +0x153
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestWithVersionEncoderRace (0.00s)
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-121771-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-121771-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-121771-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-121771-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-121771-bug .
# (then run as above, no --ssh flag)
```
