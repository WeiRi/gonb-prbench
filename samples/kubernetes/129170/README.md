# kubernetes-129170

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/129170 |
| Bug commit | `bcd65ce24025` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apimachinery/pkg/runtime/serializer/cbor/internal/modes/custom.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00011e000 by goroutine 10:
  ase/kubernetes-129170.(*checkers).getCheckerInternal()
      /work/custom.go:31 +0xdc
  ase/kubernetes-129170.(*checkers).getChecker()
      /work/custom.go:24 +0xa4
  ase/kubernetes-129170.TestRace_PR129170_LazyCheckerInit.func1()
      /work/verified_test.go:27 +0x7f

Previous write at 0x00c00011e000 by goroutine 8:
  ase/kubernetes-129170.(*checkers).getCheckerInternal()
      /work/custom.go:37 +0x184
  ase/kubernetes-129170.(*checkers).getChecker()
      /work/custom.go:24 +0xa4
  ase/kubernetes-129170.TestRace_PR129170_LazyCheckerInit.func1()
      /work/verified_test.go:27 +0x7f

Goroutine 10 (running) created at:
  ase/kubernetes-129170.TestRace_PR129170_LazyCheckerInit()
      /work/verified_test.go:25 +0x164
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/kubernetes-129170.TestRace_PR129170_LazyCheckerInit()
      /work/verified_test.go:25 +0x164
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-129170-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-129170-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-129170-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-129170-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-129170-bug .
# (then run as above, no --ssh flag)
```
