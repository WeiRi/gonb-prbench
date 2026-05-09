# kubernetes-109969

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/109969 |
| Bug commit | `564b2049231c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apiserver/pkg/authentication/group/authenticated_group_adder.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000fe020 by goroutine 9:
  ase/kubernetes-109969.(*AuthenticatedGroupAdder).AuthenticateRequest()
      /work/authenticated_group_adder.go:45 +0x2a4
  ase/kubernetes-109969.TestRace_109969_GroupAdderSharedBackingArray.func2()
      /work/verified_test.go:22 +0x1d

Previous write at 0x00c0000fe020 by goroutine 8:
  ase/kubernetes-109969.(*AuthenticatedGroupAdder).AuthenticateRequest()
      /work/authenticated_group_adder.go:45 +0x2a4
  ase/kubernetes-109969.TestRace_109969_GroupAdderSharedBackingArray.func1()
      /work/verified_test.go:21 +0x1d

Goroutine 9 (running) created at:
  ase/kubernetes-109969.TestRace_109969_GroupAdderSharedBackingArray()
      /work/verified_test.go:22 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 8 (finished) created at:
  ase/kubernetes-109969.TestRace_109969_GroupAdderSharedBackingArray()
      /work/verified_test.go:21 +0x44b
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
--- FAIL: TestRace_109969_GroupAdderSharedBackingArray (0.01s)
    testing.go:1465: race detected during execution of test
FAIL
FAIL	ase/kubernetes-109969	0.020s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-109969-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-109969-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-109969-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-109969-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-109969-bug .
# (then run as above, no --ssh flag)
```
