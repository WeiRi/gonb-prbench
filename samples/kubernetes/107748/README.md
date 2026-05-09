# kubernetes-107748

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/107748 |
| Bug commit | `6a1de6b686e4` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/proxy/ipvs/graceful_termination.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0001c2118 by goroutine 9:
  ase/kubernetes-107748.(*T_107748).write()
      /work/verified_test.go:6 +0xb3
  ase/kubernetes-107748.TestRace_107748.func1()
      /work/verified_test.go:13 +0xa3

Previous write at 0x00c0001c2118 by goroutine 8:
  ase/kubernetes-107748.(*T_107748).write()
      /work/verified_test.go:6 +0xb3
  ase/kubernetes-107748.TestRace_107748.func1()
      /work/verified_test.go:13 +0xa3

Goroutine 9 (running) created at:
  ase/kubernetes-107748.TestRace_107748()
      /work/verified_test.go:12 +0x8f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1259 +0x22f
  testing.(*T).Run·dwrap·21()
      /usr/local/go/src/testing/testing.go:1306 +0x47

Goroutine 8 (finished) created at:
  ase/kubernetes-107748.TestRace_107748()
      /work/verified_test.go:12 +0x8f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1259 +0x22f
  testing.(*T).Run·dwrap·21()
      /usr/local/go/src/testing/testing.go:1306 +0x47
==================
--- FAIL: TestRace_107748 (0.01s)
    testing.go:1152: race detected during execution of test
FAIL
FAIL	ase/kubernetes-107748	0.026s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-107748-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-107748-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-107748-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-107748-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-107748-bug .
# (then run as above, no --ssh flag)
```
